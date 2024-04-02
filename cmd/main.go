package main

import (
	"context"
	"os"
	"os/signal"
	"time"

	DB "github.com/DeltaCapstone/ChoiceMoversBackend/database"
	"github.com/DeltaCapstone/ChoiceMoversBackend/utils"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"go.uber.org/zap"
)

func main() {
	err := utils.LoadConfig()
	if err != nil {
		zap.L().Sugar().Fatalf(err.Error())
	}
	// Echo instance
	e := echo.New()

	// Middleware
	logger, _ := zap.NewDevelopment()
	e.Use(middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		LogURI:    true,
		LogStatus: true,
		LogMethod: true,
		LogValuesFunc: func(c echo.Context, v middleware.RequestLoggerValues) error {
			logger.Info("request",
				zap.String("Method", v.Method),
				zap.String("URI", v.URI),
				zap.Int("status", v.Status),
			)

			return nil
		},
	}))
	zap.ReplaceGlobals(logger)

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
		logger.Sugar().Fatalf("Unable to initialize the database:", err)
	}
	defer dbinst.Close()

	CreateRoutes(e)

	// Start server with graceful shutdown
	go func() {
		if err := e.Start(":8080"); err != nil {
			logger.Info("Shutting down the server")
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	cancel()                    // Cancel the context to stop database operations
	time.Sleep(2 * time.Second) // Allow time for existing connections to finish
}
