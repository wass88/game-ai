package main

import (
	"fmt"
	"os"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/wass88/gameai/lib/server"
)

func main() {
	dbname := os.Getenv("MYSQL_DATABASE")
	db := server.NewDB(dbname)

	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	//e.GET("/hello", handler.MainPage())
	e.POST("/api/results/:id/update", server.HandlerResultsUpdate(db))
	e.POST("/api/results/:id/complete", server.HandlerResultsComplete(db))

	e.GET("/api/games/:id/matches", server.HandlerViewMatches(db))

	addr := os.Getenv("LISTEN_ADDR")
	if addr == "" {
		addr = ":3000"
	}

	fmt.Printf("Listening on %s\n", addr)
	err := e.Start(addr)
	if err != nil {
		panic(err)
	}
}
