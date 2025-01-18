package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	time.Sleep(time.Minute)
	e := echo.New()
	e.Use(middleware.Logger())

	e.GET("/", func(c echo.Context) error {
		header := c.Request().Header
		fmt.Println(header)
		traqID, ok := header["X-Showcase-User"]
		if !ok {
			return c.String(500, "something wrong")
		}

		return c.String(200, fmt.Sprintf("こんにちは、%sさん", traqID[0]))
	})
	e.GET("/hello", func(c echo.Context) error {
		return c.String(200, "Hello, trap!")
	})
	e.GET("/env", func(c echo.Context) error {
		return c.String(200, os.Getenv("EXAMPLE_ENV"))
	})
	e.GET("/sample", func(c echo.Context) error {
		_, err := os.Open("sample.txt")
		if err != nil {
			fmt.Println("sample.txt does not exist")
			fmt.Println(err)
		} else {
			fmt.Println("sample.txt exists!")
		}
		return c.NoContent(200)
	})
	e.GET("/ignore", func(c echo.Context) error {
		_, err := os.Stat(".dockerignore")
		if err != nil {
			fmt.Println(".dockerignore does not exist")
			fmt.Println(err)
		} else {
			fmt.Println(".dockerignore exists!")
		}
		return c.NoContent(200)
	})

	go func() {
		err := e.Start(":8080")
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			e.Logger.Error(err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit
	fmt.Println("Shutting down server...")
	e.Shutdown(context.Background())
}
