package appenv

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

func Init() {
	appEnv := os.Getenv("APP_ENV")
	if appEnv == "" {
		// defaulting to development
		appEnv = "dev"
		err := os.Setenv("APP_ENV", appEnv)
		if err != nil {
			panic(fmt.Errorf("unable to set ENV \"APP_ENV\": %v", err))
		}
	}
	godotenv.Load( //             Example with APP_ENV 'dev':
		".env."+appEnv+".local", // .env.dev.local
		".env.local",            // .env.local
		".env."+appEnv,          // .env.dev
		".env",                  // .env
	)
}
