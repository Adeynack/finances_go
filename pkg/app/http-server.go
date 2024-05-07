package app

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/adeynack/finances/pkg/api"
	"github.com/go-chi/chi/v5"
)

type ServerShutdownFunc func() error

func StartHttpServer() (ServerShutdownFunc, error) {
	apiImpl := &api.Implementation{}
	middlewares := []api.StrictMiddlewareFunc{}
	router := chi.NewMux()
	strictHandler := api.NewStrictHandler(apiImpl, middlewares)
	handler := api.HandlerFromMux(strictHandler, router)

	address := fmt.Sprintf("localhost:%s", os.Getenv("PORT"))
	server := &http.Server{Handler: handler, Addr: address}

	// Start the server in the background
	go func() {
		log.Printf("Starting HTTP server at http://%s", address)
		if err := server.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("error with server: %v", err)
		}
	}()

	// Return the function to shutdown the server.
	shutdown := func() error {
		ctx, cancelTimeout := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancelTimeout()

		return server.Shutdown(ctx)
	}
	return shutdown, nil
}
