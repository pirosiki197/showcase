package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-sql-driver/mysql"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	conf := mysql.Config{
		User:                 os.Getenv("DB_USER"),
		Passwd:               os.Getenv("DB_PASSWORD"),
		Net:                  "tcp",
		Addr:                 net.JoinHostPort(os.Getenv("DB_HOST"), os.Getenv("DB_PORT")),
		DBName:               os.Getenv("DB_NAME"),
		AllowNativePasswords: true,
	}
	db, err := sql.Open("mysql", conf.FormatDSN())
	if err != nil {
		panic(err)
	}

	e := echo.New()
	e.Use(middleware.Logger())

	e.GET("/", func(c echo.Context) error {
		header := c.Request().Header
		fmt.Println(header)
		traqID, ok := header["X-Showcase-User"]
		if !ok {
			return c.String(400, "X-Showcase-User header is required")
		}

		return c.String(200, fmt.Sprintf("こんにちは、%sさん", traqID[0]))
	})
	e.GET("/hello", func(c echo.Context) error {
		return c.String(200, "Hello, trap!")
	})
	e.GET("/env", func(c echo.Context) error {
		return c.String(200, os.Getenv("EXAMPLE_ENV"))
	})
	e.GET("/json", func(c echo.Context) error {
		data := map[string]any{
			"message": "Hello, JSON!",
			"nested": map[string]any{
				"key":   "value",
				"array": []any{"item1", "item2"},
				"int":   123,
				"bool":  true,
			},
		}
		encoded, err := json.Marshal(data)
		if err != nil {
			log.Println("Error encoding JSON:", err)
			return c.String(http.StatusInternalServerError, "Internal Server Error")
		}
		fmt.Println(string(encoded))
		return c.JSON(http.StatusOK, data)
	})
	e.GET("/db", func(c echo.Context) error {
		rows, err := db.Query("SELECT name FROM users")
		if err != nil {
			log.Println("Error from DB:", err)
		}
		defer rows.Close()

		var names []string
		for rows.Next() {
			var name string
			rows.Scan(&name)
			names = append(names, name)
		}
		if err := rows.Err(); err != nil {
			log.Println("Error from DB:", err)
		}

		return c.JSON(http.StatusOK, names)
	})

	go func() {
		for range time.Tick(10 * time.Second) {
			fmt.Println("10 seconds passed")
		}
	}()

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
	time.Sleep(10 * time.Second)
	e.Shutdown(context.Background())
}
