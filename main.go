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
		User:   os.Getenv("NS_MARIADB_USER"),
		Passwd: os.Getenv("NS_MARIADB_PASSWORD"),
		Net:    "tcp",
		Addr:   os.Getenv("NS_MARIADB_HOSTNAME") + os.Getenv("NS_MARIADB_PORT"),
		DBName: os.Getenv("NS_MARIADB_DATABASE"),
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
		return c.String(200, "NeoShowcaseを使いたい")
	})
	e.GET("/hello", func(c echo.Context) error {
		return c.String(200, "Hello, trap!")
	})

	e.Logger.Fatal(e.Start(":8080"))
}
