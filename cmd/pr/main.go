package main

import (
	"net/http"

	"github.com/NutsBalls/Backend-trainee-assignment-autumn-2025/pkg/config"
	database "github.com/NutsBalls/Backend-trainee-assignment-autumn-2025/pkg/database"
	"github.com/labstack/echo/v4"
)

func main() {
	cfg := config.Load()
	db := database.NewConn(cfg.DBURL)
	defer db.Close()

	e := echo.New()

	e.GET("/ping", func(c echo.Context) error {
		return c.String(http.StatusOK, "pong")
	})

	e.Logger.Fatal(e.Start(":" + cfg.Port))

}
