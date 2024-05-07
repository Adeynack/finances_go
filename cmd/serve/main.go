package main

import (
	"log"
	"net/http"

	"github.com/adeynack/finances/pkg/api"
	"github.com/go-chi/chi/v5"
)

func main() {
	apiImpl := &api.Implementation{}
	middlewares := []api.StrictMiddlewareFunc{}
	router := chi.NewMux()
	strictHandler := api.NewStrictHandler(apiImpl, middlewares)
	handler := api.HandlerFromMux(strictHandler, router)
	server := &http.Server{Handler: handler, Addr: "0.0.0.0:8080"}
	log.Fatal(server.ListenAndServe())
}
