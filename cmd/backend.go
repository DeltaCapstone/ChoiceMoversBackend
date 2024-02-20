package main

import (
	"context"
	"net/http"

	db "github.com/DeltaCapstone/ChoiceMoversBackend/database"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	// Echo instance
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	dbinst, _ := db.NewPG(context.Background())
	defer dbinst.Close()

	dbinst.Ping(context.Background())

	// Routes
	e.GET("/", hello)
	e.GET("/dakota", getDakota)

	// Start server
	e.Logger.Fatal(e.Start(":8080"))
}

// Handler
func hello(c echo.Context) error {
	return c.String(http.StatusOK, "Hello, World!")
}

func getDakota(c echo.Context) error {
	value := dbinst.getName(c, "dakota")
	return c.String(http.StatusOK, value)
}
