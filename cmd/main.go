package main

import (
	"context"
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
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"http://localhost:4200"},
		AllowMethods: []string{echo.GET, echo.PUT, echo.POST, echo.DELETE},
	}))

	// Initialize database with a context that can be canceled
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	dbinst, err := DB.NewPG(ctx)
	if err != nil {
		e.Logger.Fatal("Unable to initialize the database:", err)
	}
	defer dbinst.Close()

	CreateRoutes(e)

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
