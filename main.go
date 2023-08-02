package main

import (
	"log"
	"os"

	"github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	conf := mysql.Config{
		User:                 getEnvOrDefault("NS_MARIADB_USER", "root"),
		Passwd:               getEnvOrDefault("NS_MARIADB_PASSWORD", "password"),
		Net:                  "tcp",
		Addr:                 getEnvOrDefault("NS_MARIADB_HOSTNAME", "localhost") + ":" + getEnvOrDefault("NS_MARIADB_PORT", "3306"),
		DBName:               getEnvOrDefault("NS_MARIADB_DATABASE", "showcase"),
		AllowNativePasswords: true,
	}
	db, err := sqlx.Open("mysql", conf.FormatDSN())
	if err != nil {
		log.Fatal(err)
	}
	if err := db.Ping(); err != nil {
		log.Fatal(err)
	}

	e := echo.New()
	e.Use(middleware.Logger())

	e.GET("/", func(c echo.Context) error {
		header := c.Request().Header
		log.Println(header)
		return c.String(200, "NeoShowcaseを使いたい")
	})
	e.GET("/hello", func(c echo.Context) error {
		return c.String(200, "Hello, trap!")
	})
	e.GET("/env", func(c echo.Context) error {
		return c.String(200, os.Getenv("EXAMPLE_ENV"))
	})

	e.Logger.Fatal(e.Start(":8080"))
}

func getEnvOrDefault(key, d string) string {
	value := os.Getenv(key)
	if value == "" {
		return d
	}
	return value
}
