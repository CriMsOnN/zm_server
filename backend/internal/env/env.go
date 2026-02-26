package env

import (
	"os"
	"strconv"
)

func GetEnvString(key string, fallback string) string {
	if found, ok := os.LookupEnv(key); ok {
		return found
	}

	return fallback
}

func GetEnvNumber(key string, fallback int) int {
	if found, ok := os.LookupEnv(key); ok {
		intValue, err := strconv.Atoi(found)
		if err != nil {
			return fallback
		}
		return intValue
	}

	return fallback
}

func GetEnvBool(key string, fallback bool) bool {
	if found, ok := os.LookupEnv(key); ok {
		boolValue, err := strconv.ParseBool(found)
		if err != nil {
			return fallback
		}
		return boolValue
	}

	return fallback
}
