package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"time"

	DB "github.com/DeltaCapstone/ChoiceMoversBackend/database"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	// Echo instance
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Initialize database with a context that can be canceled
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	dbinst, err := DB.NewPG(ctx)
	if err != nil {
		e.Logger.Fatal("Unable to initialize the database:", err)
	}
	defer dbinst.Close()

	dbinst.Ping(ctx)

	// Routes
	e.GET("/", hello)
	//test route to run a very basic querry
	e.GET("/dakota", getDakota)

	// Start server with graceful shutdown
	go func() {
		if err := e.Start(":8080"); err != nil {
			e.Logger.Info("Shutting down the server")
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	cancel()                    // Cancel the context to stop database operations
	time.Sleep(2 * time.Second) // Allow time for existing connections to finish
}

// Handler
func hello(c echo.Context) error {
	return c.String(http.StatusOK, "Hello, World!")
}

// Test with basic query
func getDakota(c echo.Context) error {
	value, err := DB.PgInstance.GetName(c.Request().Context(), "dakota")
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error()+" Error retrieving data")
	}
	return c.String(http.StatusOK, "dakota's id is: "+value)
}
