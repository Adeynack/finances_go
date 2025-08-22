package app

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/adeynack/finances/pkg/api/apiserver"
	"github.com/adeynack/finances/pkg/platform/ctxval"
	"github.com/adeynack/finances/pkg/repository"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	slogctx "github.com/veqryn/slog-context"
)

type ServerShutdownFunc func() error

func MustStartHttpServer() ServerShutdownFunc {
	handler := mustCreateHandler(false)

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

func mustCreateTestHandler() http.Handler {
	return mustCreateHandler(true)
}

func mustCreateHandler(testing bool) http.Handler {
	db := repository.NewDB(mustConnectDatabase())

	repo, err := repository.New()
	if err != nil {
		panic(err)
	}

	apiImpl := &apiserver.Service{Repo: repo}
	middlewares := []apiserver.StrictMiddlewareFunc{}
	router := chi.NewMux()
	router.Use(
		middleware.RequestID,
		requestIDStructuredLog,
		middleware.Logger,
		middleware.Timeout(30*time.Second),
		injectDBConnection(db, testing),
	)
	strictHandler := apiserver.NewStrictHandler(apiImpl, middlewares)
	return apiserver.HandlerFromMux(strictHandler, router)
}

func requestIDStructuredLog(next http.Handler) http.Handler {
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

func injectDBConnection(db repository.DB, testing bool) func(next http.Handler) http.Handler {
	if testing {
		return func(next http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				ctx := r.Context()
				db, err := db.BeginTx(ctx)
				if err != nil {
					w.WriteHeader(http.StatusInternalServerError)
					return
				}

				defer db.Rollback()
				next.ServeHTTP(w, r.WithContext(ctxval.Register(ctx, db)))
			})
		}
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			next.ServeHTTP(w, r.WithContext(ctxval.Register(r.Context(), db)))
		})
	}
}
