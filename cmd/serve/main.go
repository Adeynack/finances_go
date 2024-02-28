package main

import (
	"context"
	"log"
	"os"
	"os/signal"

	"github.com/adeynack/finances/app"
	"github.com/adeynack/finances/model/database"
)

func main() {
	shutdownServer, err := app.StartHttpServer()
	if err != nil && !database.FatalLogIfTrappedMigrationError(err) {
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
