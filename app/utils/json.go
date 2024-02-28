package utils

import (
	"encoding/json"
	"fmt"
)

// MustJSON marhals an object to a JSON string. As it panics if an error occurs,
// it is meant as a quick debug convenience. Use proper `json.Marshal` with error
// handling for production-path JSON serialization.
func MustJSON(obj any) []byte {
	marshaled, err := json.Marshal(obj)
	if err != nil {
		panic(fmt.Sprintf("error mashaling model %T to JSON: %v", obj, err))
	}
	return marshaled
}

func MustJSONPretty(obj any) []byte {
	marshaled, err := json.MarshalIndent(obj, "", "  ")
	if err != nil {
		panic(fmt.Sprintf("error mashaling model %T to JSON: %v", obj, err))
	}
	return marshaled
}
