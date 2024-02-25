package app

import (
	"os"

	"github.com/joho/godotenv"
)

func InitAppEnvironment() {
	appEnv := os.Getenv("APP_ENV")
	if appEnv == "" {
		// defaulting to development
		appEnv = "dev"
		os.Setenv("APP_ENV", appEnv)
	}
	godotenv.Load( //             Example with APP_ENV 'dev':
		".env."+appEnv+".local", // .env.dev.local
		".env.local",            // .env.local
		".env."+appEnv,          // .env.dev
		".env",                  // .env
	)
}
