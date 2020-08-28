package main

import (
	"fmt"
	"os"

	"github.com/go-playground/validator"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/wass88/gameai/lib/server"
)

type CustomValidator struct {
	validator *validator.Validate
}

func (cv *CustomValidator) Validate(i interface{}) error {
	return cv.validator.Struct(i)
}

func main() {
	addr := os.Getenv("LISTEN_ADDR")
	if addr == "" {
		addr = ":3000"
	}

	confFile := os.Getenv("CONF_FILE")
	conf, err := server.LoadConfig(confFile)
	if err != nil {
		panic(err)
	}

	db := conf.NewDB()

	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Validator = &CustomValidator{validator: validator.New()}

	// Change Listen port
	e.POST("/api/container/ready", server.HandlerReadyContainer(db))
	e.POST("/api/results/:id/update", server.HandlerResultsUpdate(db))
	e.POST("/api/results/:id/complete", server.HandlerResultsComplete(db))

	e.GET("/api/games/:id/matches", server.HandlerViewMatches(db))
	e.GET("/api/matches/:id", server.HandlerViewMatch(db))
	e.GET("/api/games/:id/ai-githubs", server.HandlerViewAIGithubByGame(db))
	e.GET("/api/games/:id/latest-ai", server.HandlerViewLatestByGame(db))

	e.POST("/api/matches", server.HandlerAddMatch(db))
	e.POST("/api/ai-githubs", server.HandlerAddAIGithub(db))

	db.SetSessionHandler(e)
	e.GET("/api/you", server.HandlerYou(db))

	fmt.Printf("Listening on %s\n", addr)
	err = e.Start(addr)
	if err != nil {
		panic(err)
	}
}
