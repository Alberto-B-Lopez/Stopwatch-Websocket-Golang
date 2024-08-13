package main

import (
	"context"
	"fmt"
	"math/rand"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/lopez/websockets/components"
)

func homeHandler(c echo.Context) error {
	return render(c, components.Base())
}

func wsHandler(c echo.Context, ctx context.Context) error {
	var watcher *Watcher
	id := c.Param("id")

	conn, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		return err
	}

	watcher, ok := watchers[id]

	if !ok {
		fmt.Println(id)
		watcher = NewWatcher(id, conn)
		watchers[id] = watcher
	}

	if !watcher.Paused {
		go watcher.Start(ctx)
	}

	go watcher.Read(ctx)
	go watcher.Status(ctx)

	return nil
}

func addTimerHandler(c echo.Context) error {
	randInt := rand.Intn(1000000000000)
	id := strconv.Itoa(randInt)
	return render(c, components.Websocket(id))
}
