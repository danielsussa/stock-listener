package main

import (
	"encoding/json"
	"net/http"

	"github.com/labstack/echo"
	"golang.org/x/net/websocket"
)

func serveWeb() {
	e := echo.New()
	e.GET("/stats", func(c echo.Context) error {
		_, err := json.Marshal(stockMap)
		if err != nil {
			panic(err)
		}
		return c.JSON(http.StatusOK, map[string]interface{}{
			"stockMap": stockMap,
			"snapshot": snapshotMap,
		})
	})
	e.GET("/ws", hello)
	e.GET("/options", optionsHandler)
	e.Static("/", "api-core/assets/index.html")
	e.Logger.Fatal(e.Start(":8099"))
}

func optionsHandler(c echo.Context) error {
	options := make([]*option, 0)
	for _, v := range stockMap {
		if v.Kind() == (option{}).Kind() {
			options = append(options, v.(*option))
		}
	}
	c.JSON(http.StatusOK, options)
	return nil
}

var currentOptUpdate *option

func hello(c echo.Context) error {
	websocket.Handler(func(ws *websocket.Conn) {
		defer ws.Close()
		for {
			if currentOptUpdate != nil {
				err := websocket.JSON.Send(ws, currentOptUpdate)
				if err != nil {
					c.Logger().Error(err)
				}
				currentOptUpdate = nil
			}
		}
	}).ServeHTTP(c.Response(), c.Request())
	return nil
}
