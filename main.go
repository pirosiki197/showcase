package main

import (
	"fmt"
	"os"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
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
	e.Get("/sample", func(c echo.Context) error {
		_, err := os.Open("sample.txt")
		if err != nil {
			fmt.Println("sample.txt does not exist")
			fmt.Println(err)
		}
		return c.NoContent(204)
	})

	e.Logger.Fatal(e.Start(":8080"))
}
