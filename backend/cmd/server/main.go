package main

import (
	"fmt"
	"os"

	"github.com/go-playground/validator"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/pkg/errors"
	"github.com/wass88/gameai/lib/server"
)

type CustomValidator struct {
	validator *validator.Validate
}

func (cv *CustomValidator) Validate(i interface{}) error {
	return cv.validator.Struct(i)
}

func main() {

	confFile := os.Getenv("CONF_FILE")
	if confFile == "" {
		panic(errors.Errorf("Missing CONF_FILE enviroment variable"))
	}
	conf, err := server.LoadConfig(confFile)
	if err != nil {
		panic(err)
	}
	addr := conf.APIAddr
	if addr == "" {
		panic(errors.Errorf("Missing api_addr"))
	}
	insideAddr := conf.InsideAPIAddr
	if addr == "" {
		panic(errors.Errorf("Missing inside_addr"))
	}

	db := conf.NewDB()

	e := echo.New()
	eInside := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Validator = &CustomValidator{validator: validator.New()}

	eInside.Use(middleware.Logger())
	eInside.Use(middleware.Recover())
	eInside.Validator = &CustomValidator{validator: validator.New()}

	eInside.POST("/inside_api/container/ready", server.HandlerReadyContainer(db))
	eInside.POST("/inside_api/results/:id/update", server.HandlerResultsUpdate(db))
	eInside.POST("/inside_api/results/:id/complete", server.HandlerResultsComplete(db))

	e.GET("/api/games/:id/matches", server.HandlerViewMatches(db))
	e.GET("/api/matches/:id", server.HandlerViewMatch(db))
	e.GET("/api/games/:id/ai-githubs", server.HandlerViewAIGithubByGame(db))
	e.GET("/api/games/:id/latest-ai", server.HandlerViewLatestByGame(db))

	e.POST("/api/matches", server.HandlerAddMatch(db))
	e.POST("/api/ai-githubs", server.HandlerAddAIGithub(db))

	db.SetSessionHandler(e)
	e.GET("/api/you", server.HandlerYou(db))

	go func() {
		fmt.Printf("Listening API on %s\n", addr)
		err = e.Start(addr)
		if err != nil {
			panic(err)
		}
	}()
	fmt.Printf("Listening Inside API on %s\n", insideAddr)
	err = eInside.Start(insideAddr)
	if err != nil {
		panic(err)
	}
}
