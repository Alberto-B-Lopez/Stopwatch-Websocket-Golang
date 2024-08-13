package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"

	"github.com/labstack/echo/v4"
)

var (
	upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
	watchers = make(map[string]*Watcher)
)

func main() {
	e := echo.New()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	e.GET("/", homeHandler)

	e.GET("/ws/:id", func(c echo.Context) error {
		fmt.Println("we going to upgrade now")
		return wsHandler(c, ctx)
	})

	e.GET("/addTimer", addTimerHandler)

	log.Fatal(e.Start(":8080"))
}
