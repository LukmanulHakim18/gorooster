package helpers

import (
	"os"
	"strconv"
)

// Get Value(Int) from env
// And give default value if err parsing
func EnvGetInt(key string, defaultValue int) int {
	if val, err := strconv.ParseInt(os.Getenv(key), 10, 64); err == nil {
		return int(val)
	}
	return defaultValue
}

// Get Value(String) from env
// And give default value if err or empty
func EnvGetString(key string, defaultValue string) string {
	if val := os.Getenv((key)); val != "" {
		return val
	}
	return defaultValue
}

// Get Value(Bool) from env
// And give default value if err parsing
func EnvGetBool(key string, defaultValue bool) bool {
	if val, err := strconv.ParseBool(os.Getenv((key))); err == nil {
		return val
	}
	return defaultValue
}
