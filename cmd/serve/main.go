package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strconv"

	"github.com/adeynack/finances/server"
)

func main() {
	config, err := parseConfiguration()
	if err != nil {
		log.Fatalln(err)
	}
	shutdownServer := server.StartHttpServer(config)

	// Listen for interrupt signal (eg: Ctrl-C) to gracefully shutdown the server
	interruptCtx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()
	<-interruptCtx.Done()
	log.Println("server shutting down")
	if err := shutdownServer(); err != nil {
		log.Fatalf("error shutting server down: %s", err)
	} else {
		log.Println("server gracefully terminated")
	}
}

func parseConfiguration() (server.ServerConfiguration, error) {
	var err error
	config := server.ServerConfiguration{}

	if portEnv := os.Getenv("PORT"); portEnv == "" {
		config.Port = 40001
	} else {
		config.Port, err = strconv.Atoi(portEnv)
		if err != nil {
			return config, fmt.Errorf("could not convert %q to a port number: %s", portEnv, err)
		}
	}

	return config, nil
}
