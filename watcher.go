package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/gorilla/websocket"
	"github.com/lopez/websockets/components"
)

type Command struct {
	Command string `json:"command"`
}

type Watcher struct {
	Id        string
	Conn      *websocket.Conn
	Done      chan bool
	Paused    bool
	Listen    chan string
	Ticker    *time.Ticker
	StartTime time.Time
	CurrTime  time.Duration
}

func NewWatcher(id string, conn *websocket.Conn) *Watcher {
	return &Watcher{
		Id:        id,
		Conn:      conn,
		Done:      make(chan bool),
		Paused:    false,
		Listen:    make(chan string),
		Ticker:    time.NewTicker(1 * time.Second),
		StartTime: time.Now(),
	}
}

func (w *Watcher) Start(ctx context.Context) {
	defer func() {
		fmt.Println("we closeing connection")
	}()

	for {
		select {
		case <-w.Ticker.C:
			if w.Paused {
				fmt.Println("we are paused")
				continue
			}
			w.CurrTime = time.Since(w.StartTime)
			timeValue := time.Since(w.StartTime).String()
			component := components.Stopwatch(timeValue, w.Id)
			fmt.Println("Time sent to client:", timeValue)
			fmt.Println("Watcher:", w.Id, "Paused:", w.Paused)
			buffer := &bytes.Buffer{}
			component.Render(ctx, buffer)
			err := w.Conn.WriteMessage(websocket.TextMessage, buffer.Bytes())
			if err != nil {
				fmt.Println("Error writing to websocket:", err)
				return
			}

		case <-ctx.Done():
			fmt.Println("Client disconnected")
			return
		}
	}
}

func (w *Watcher) Status(ctx context.Context) {
	defer func() {
		fmt.Println("we closeing connection")
		w.Conn.Close()
	}()

	for {
		select {
		case status := <-w.Listen:

			switch status {
			case "stop":
				w.Paused = true
				w.Ticker.Stop()
				fmt.Println("We are stopping the watcher")
			case "resume":
				w.Paused = false
				w.Ticker.Reset(1 * time.Second)
				w.StartTime = time.Now().Add(-w.CurrTime)
				fmt.Println("We are resuming the watcher")
			case "reset":
				w.Paused = false
				w.Ticker.Reset(1 * time.Second)
				w.StartTime = time.Now()
				fmt.Println("Watcher:", w.Id, "Reset")
			}
		case <-ctx.Done():
			return
		}
	}
}

func (w *Watcher) Read(ctx context.Context) {
	defer func() {
		w.Conn.Close()
	}()
	var command Command

	for {
		_, msg, err := w.Conn.ReadMessage()
		if err != nil {
			fmt.Println("read error:", err)
			return
		}
		fmt.Printf("%s\n", msg)
		jerr := json.Unmarshal(msg, &command)
		if jerr != nil {
			fmt.Println("error:", err)
		}
		fmt.Println("Command:", command.Command)
		w.Listen <- command.Command

	}
}
