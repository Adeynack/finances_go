package app

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/adeynack/finances/app/appenv"
	"github.com/adeynack/finances/controller"
	"github.com/adeynack/finances/model/database"
	"github.com/labstack/echo/v4"
)

func init() {
	appenv.Init()
}

func StartHttpServer() (ServerShutdownFunc, error) {
	db, err := database.Connect()
	if err != nil {
		return nil, err
	}

	// Create & initialize HTTP server
	e := echo.New()
	e.HideBanner = true

	appContext := controller.AppContext{
		ServerSecret: extractServerSecretOrFail(),
		DB:           db,
	}
	controller.DrawRoutes(e, &appContext)

	// Start the server in the background
	go func() {
		address := fmt.Sprintf("localhost:%s", os.Getenv("PORT"))
		if err := e.Start(address); err != nil && !errors.Is(err, http.ErrServerClosed) {
			e.Logger.Fatalf("error shutting down server: %w", err)
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

func extractServerSecretOrFail() string {
	secret := strings.Trim(os.Getenv(appenv.EnvServerSecret), " ")
	if secret == "" {
		log.Fatalln("Server secret must be set. TIP: If you see this, either set ENV %q on your server, or specify it in `.env.local` if developping.", appenv.EnvServerSecret)
	}
	return secret
}
