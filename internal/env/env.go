package env

import (
	"fmt"
	"os"
	"strconv"
)

func GetEnvString(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	// Add this line temporarily
	fmt.Printf("Key %s not found, using default: %s\n", key, defaultValue)
	return defaultValue
}

func GetEnvInt(key string, defaultValue int) int {
	if value, exists := os.LookupEnv(key); exists {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}

	return defaultValue
}
