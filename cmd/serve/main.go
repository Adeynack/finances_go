package main

import (
	"context"
	"log"
	"os"
	"os/signal"

	"github.com/adeynack/finances/pkg/app"
)

func main() {
	// Bootstrap the server
	shutdownServer, err := app.StartHttpServer()
	if err != nil {
		log.Fatalln(err)
	}

	// Listen for interrupt signal (eg: Ctrl-C) to gracefully shutdown the server
	interruptCtx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	<-interruptCtx.Done()

	log.Println("server shutting down")
	if err := shutdownServer(); err != nil {
		log.Fatal(err)
	} else {
		log.Println("server gracefully terminated")
	}
}
