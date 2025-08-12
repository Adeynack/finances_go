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
	"github.com/go-chi/chi/v5/middleware"
	slogctx "github.com/veqryn/slog-context"
)

type ServerShutdownFunc func() error

func MustStartHttpServer() ServerShutdownFunc {
	handler := mustCreateHandler()

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
	return shutdown
}

func mustCreateHandler() http.Handler {
	apiImpl := &api.Service{
		DB: mustConnectDatabase(),
	}
	middlewares := []api.StrictMiddlewareFunc{}
	router := chi.NewMux()
	router.Use(
		middleware.RequestID,
		RequestIDStructuredLog,
		middleware.Logger,
		middleware.Timeout(30*time.Second),
	)
	strictHandler := api.NewStrictHandler(apiImpl, middlewares)
	return api.HandlerFromMux(strictHandler, router)
}

func RequestIDStructuredLog(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r = r.WithContext(
			slogctx.With(
				r.Context(),
				"request_id",
				r.Context().Value(middleware.RequestIDKey),
			))
		next.ServeHTTP(w, r)
	})
}
