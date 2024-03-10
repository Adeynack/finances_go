//go:build !test
// +build !test

package appenv

import (
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/adeynack/finances/app/utils"
	"github.com/joho/godotenv"
)

const (
	EnvAppEnv       = "APP_ENV"
	EnvDotEnvLoaded = "DOT_ENV_LOADED"
	EnvServerSecret = "SERVER_SECRET"
)

var Env string

// Init makes sure the `dotenv` environment is loaded.
// This should be called in the `init()` of every package that
// needs environment variables to be properly set (eg: the database).
func Init() {
	if Env != "" {
		return // means this function was already executed once.
	}
	if utils.ReadEnvBoolean(EnvDotEnvLoaded, false) {
		return // means the app is already loaded through dotenv-cli
	}

	Env = os.Getenv(EnvAppEnv)
	if Env == "" {
		Env = "dev" // defaulting to development
		err := os.Setenv(EnvAppEnv, Env)
		if err != nil {
			panic(fmt.Errorf("unable to set ENV %q: %v", EnvAppEnv, err))
		}
	}

	// This load order must match the `dotenv-cli` cascading logic.
	// Example (from executing `dotenv --debug -c dev`):
	// [ '.env.dev.local', '.env.local', '.env.dev', '.env' ]
	dotEnvFiles := make([]string, 0)
	dotEnvFiles = appendFileIfExist(dotEnvFiles, ".env."+Env+".local")
	dotEnvFiles = appendFileIfExist(dotEnvFiles, ".env.local")
	dotEnvFiles = append(dotEnvFiles, ".env."+Env, ".env")
	err := godotenv.Load(dotEnvFiles...)
	if err != nil {
		log.Fatalf("error loading godotenv: %v", err)
	}
}

func appendFileIfExist(files []string, path string) []string {
	if _, err := os.Stat(path); err == nil {
		return append(files, path)
	} else if errors.Is(err, os.ErrNotExist) {
		return files
	} else {
		panic(fmt.Errorf("error getting file stat for %q: %v", path, err))
	}
}
