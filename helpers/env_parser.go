package helpers

import (
	"os"
	"strconv"
	"time"
)

// Get Value(int) from env
// And give default value if err parsing
func EnvGetInt(key string, defaultValue int) int {
	if val, err := strconv.ParseInt(os.Getenv(key), 10, 64); err == nil {
		return int(val)
	}
	return defaultValue
}

// Get Value(string) from env
// And give default value if err or empty
func EnvGetString(key string, defaultValue string) string {
	if val := os.Getenv((key)); val != "" {
		return val
	}
	return defaultValue
}

// Get Value(bool) from env
// And give default value if err parsing
func EnvGetBool(key string, defaultValue bool) bool {
	if val, err := strconv.ParseBool(os.Getenv((key))); err == nil {
		return val
	}
	return defaultValue
}

// Get Value(bool) from env
// And give default value if err parsing
func EnvGetTimeDuration(key string, defaultValue time.Duration) time.Duration {
	if val, err := time.ParseDuration(os.Getenv((key))); err != nil {
		return val
	}
	return defaultValue
}
