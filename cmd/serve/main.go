package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"strings"

	"github.com/adeynack/finances/app"
	"github.com/adeynack/finances/db"
)

func main() {
	app.InitAppEnvironment()

	shutdownServer, err := app.StartHttpServer()
	if trappedMigrationsError, ok := err.(*db.TrappedMigrationsError); ok {
		// TODO: Move to future `cmd/tool` or `cmd/dev` dev-ops binary
		log.Fatalf("\n\nPending migrations:\n%s\n\n", strings.Join(trappedMigrationsError.PendingMigrations, "\n"))
	} else if err != nil {
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
