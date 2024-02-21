package server

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/adeynack/finances/server/routes"
	"github.com/labstack/echo/v4"
)

func StartHttpServer(config ServerConfiguration) ServerShutdownFunc {
	// Create & initialize http server
	e := echo.New()
	routes.Draw(e)

	// Start the server in the background
	go func() {
		address := fmt.Sprintf(":%d", config.Port)
		if err := e.Start(address); err != nil && err != http.ErrServerClosed {
			e.Logger.Fatalf("error shutting down server: %s", err)
		}
	}()

	// Return the function to shutdown the server.
	return func() error {
		ctx, cancelTimeout := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancelTimeout()
		return e.Shutdown(ctx)
	}
}

type ServerShutdownFunc func() error

type ServerConfiguration struct {
	Port int
}
