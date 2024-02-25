package utils

import (
	"os"
	"strings"
)

func ReadEnvBoolean(key string, defaultValue bool) bool {
	switch strings.ToUpper(os.Getenv(key)) {
	case "1", "true", "t", "yes", "y":
		return true
	case "0", "false", "f", "no", "n":
		return false
	default:
		return defaultValue
	}
}
