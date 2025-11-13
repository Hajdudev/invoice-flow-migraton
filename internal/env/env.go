package env

import (
	"fmt"
	"os"
	"strconv"
	"sync"
)

var (
	cache = make(map[string]string)
	mu    sync.RWMutex
)

// getFromCacheOrEnv retrieves value from cache or environment
func getFromCacheOrEnv(key string) (string, bool) {
	mu.RLock()
	val, exists := cache[key]
	mu.RUnlock()

	if exists {
		return val, true
	}

	val, ok := os.LookupEnv(key)
	if !ok {
		return "", false
	}

	mu.Lock()
	cache[key] = val
	mu.Unlock()
	return val, true
}

// ClearCache clears the entire cache (useful for testing)
func ClearCache() {
	mu.Lock()
	defer mu.Unlock()
	cache = make(map[string]string)
}

// ============================================================================
// STRING FUNCTIONS
// ============================================================================

// GetString retrieves a string value with a fallback
func GetString(key, fallback string) string {
	val, ok := getFromCacheOrEnv(key)
	if !ok {
		return fallback
	}
	return val
}

// GetStringOrDefault retrieves a string value with a fallback (alias for GetString)
func GetStringOrDefault(key, fallback string) string {
	return GetString(key, fallback)
}

// GetStringOrThrow retrieves a string value or returns an error if not found
func GetStringOrThrow(key string) (string, error) {
	val, ok := getFromCacheOrEnv(key)
	if !ok {
		return "", fmt.Errorf("environment variable %s is not set", key)
	}
	if val == "" {
		return "", fmt.Errorf("environment variable %s is empty", key)
	}
	return val, nil
}

// MustGetString retrieves a string value or panics if not found
func MustGetString(key string) string {
	val, err := GetStringOrThrow(key)
	if err != nil {
		panic(err)
	}
	return val
}

// ============================================================================
// INT FUNCTIONS
// ============================================================================

// GetInt retrieves an int value with a fallback
func GetInt(key string, fallback int) int {
	val, ok := getFromCacheOrEnv(key)
	if !ok {
		return fallback
	}

	valAsInt, err := strconv.Atoi(val)
	if err != nil {
		return fallback
	}

	return valAsInt
}

// GetIntOrDefault retrieves an int value with a fallback (alias for GetInt)
func GetIntOrDefault(key string, fallback int) int {
	return GetInt(key, fallback)
}

// GetIntOrThrow retrieves an int value or returns an error if not found or invalid
func GetIntOrThrow(key string) (int, error) {
	val, ok := getFromCacheOrEnv(key)
	if !ok {
		return 0, fmt.Errorf("environment variable %s is not set", key)
	}
	if val == "" {
		return 0, fmt.Errorf("environment variable %s is empty", key)
	}

	valAsInt, err := strconv.Atoi(val)
	if err != nil {
		return 0, fmt.Errorf("environment variable %s is not a valid integer: %w", key, err)
	}

	return valAsInt, nil
}

// MustGetInt retrieves an int value or panics if not found or invalid
func MustGetInt(key string) int {
	val, err := GetIntOrThrow(key)
	if err != nil {
		panic(err)
	}
	return val
}

// ============================================================================
// BOOL FUNCTIONS
// ============================================================================

// GetBool retrieves a bool value with a fallback
func GetBool(key string, fallback bool) bool {
	val, ok := getFromCacheOrEnv(key)
	if !ok {
		return fallback
	}

	boolVal, err := strconv.ParseBool(val)
	if err != nil {
		return fallback
	}

	return boolVal
}

// GetBoolOrDefault retrieves a bool value with a fallback (alias for GetBool)
func GetBoolOrDefault(key string, fallback bool) bool {
	return GetBool(key, fallback)
}

// GetBoolOrThrow retrieves a bool value or returns an error if not found or invalid
func GetBoolOrThrow(key string) (bool, error) {
	val, ok := getFromCacheOrEnv(key)
	if !ok {
		return false, fmt.Errorf("environment variable %s is not set", key)
	}
	if val == "" {
		return false, fmt.Errorf("environment variable %s is empty", key)
	}

	boolVal, err := strconv.ParseBool(val)
	if err != nil {
		return false, fmt.Errorf("environment variable %s is not a valid boolean: %w", key, err)
	}

	return boolVal, nil
}

// MustGetBool retrieves a bool value or panics if not found or invalid
func MustGetBool(key string) bool {
	val, err := GetBoolOrThrow(key)
	if err != nil {
		panic(err)
	}
	return val
}
