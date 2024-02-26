package app

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/adeynack/finances/app/routes"
	"github.com/adeynack/finances/db"
	"github.com/adeynack/finances/model"
	"github.com/labstack/echo/v4"
)

func StartHttpServer() (ServerShutdownFunc, error) {
	db, err := db.Connect()
	if err != nil {
		return nil, err
	}
	// TODO: Temporary. Just doing something with the DB to assert it works & logs properly.
	_ = db
	var user model.User
	db.First(&user)
	fmt.Printf("\nUser:\n%s\n\n", user.MustJSON())
	// end TODO

	// Create & initialize http server
	e := echo.New()
	e.HideBanner = true
	routes.Draw(e)

	// Start the server in the background
	go func() {
		address := fmt.Sprintf("localhost:%s", os.Getenv("PORT"))
		if err := e.Start(address); err != nil && err != http.ErrServerClosed {
			e.Logger.Fatalf("error shutting down server: %s", err)
		}
	}()

	// Return the function to shutdown the server.
	shutdown := func() error {
		ctx, cancelTimeout := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancelTimeout()
		return e.Shutdown(ctx)
	}
	return shutdown, nil
}

type ServerShutdownFunc func() error
