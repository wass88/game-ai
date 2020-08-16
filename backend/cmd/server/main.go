package main

import (
	"fmt"
	"os"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/wass88/gameai/lib"
)

func main() {
	dbname := os.Getenv("MYSQL_DATABASE")
	db := lib.NewDB(dbname)

	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	//e.GET("/hello", handler.MainPage())
	e.POST("/api/results/:id/update", lib.HandlerResultsUpdate(db))
	e.POST("/api/results/:id/complete", lib.HandlerResultsComplete(db))

	addr := os.Getenv("LISTEN_ADDR")
	if addr == "" {
		addr = ":3000"
	}

	fmt.Printf("Listening on %s\n", addr)
	e.Start(addr)
}
