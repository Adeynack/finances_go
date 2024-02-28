//go:build !test
// +build !test

package app

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

var Env string

func init() {
	Env = os.Getenv("APP_ENV")
	if Env == "" {
		// checking if running tests
		Env = "dev" // defaulting to development
		err := os.Setenv("APP_ENV", Env)
		if err != nil {
			panic(fmt.Errorf("unable to set ENV \"APP_ENV\": %v", err))
		}
	}
	godotenv.Load( //           Example with APP_ENV 'dev':
		".env."+Env+".local", // .env.dev.local
		".env.local",         // .env.local
		".env."+Env,          // .env.dev
		".env",               // .env
	)
}
