// Package config (testing suite) contains tests for the configuration management functionalities in the GopherSignal application.

package config

import (
	"os"
	"testing"
)

// TestGetEnv verifies the behavior of the GetEnv function, which retrieves environment variables.
// It tests the retrieval of both existing and non-existing environment variables,
// ensuring that the function returns the expected values or fallback values.
func TestGetEnv(t *testing.T) {
	// Setting up a test environment variable.
	os.Setenv("TEST_KEY", "test_value")
	defer os.Unsetenv("TEST_KEY")

	// Case 1: Retrieve an existing environment variable.
	if got := GetEnv("TEST_KEY", "fallback"); got != "test_value" {
		t.Errorf("GetEnv() = %s; want test_value", got)
	}

	// Case 2: Retrieve a non-existing variable, expect fallback value.
	if got := GetEnv("NON_EXISTENT_KEY", "fallback_value"); got != "fallback_value" {
		t.Errorf("GetEnv() = %s; want fallback_value", got)
	}
}

// TestNewConfig ensures that the NewConfig function correctly initializes a Config struct using environment variables.
func TestNewConfig(t *testing.T) {
	// Preparing test environment variables.
	os.Setenv("MYSQL_USER", "test_user")
	os.Setenv("MYSQL_PASSWORD", "test_password")
	os.Setenv("MYSQL_HOST", "test_host")
	os.Setenv("MYSQL_PORT", "3306")
	os.Setenv("MYSQL_DATABASE", "test_database")
	os.Setenv("GO_ENV", "test_go_env")
	os.Setenv("SERVER_ADDRESS", "test_server_addr")
	defer func() {
		os.Unsetenv("MYSQL_USER")
		os.Unsetenv("MYSQL_PASSWORD")
		os.Unsetenv("MYSQL_HOST")
		os.Unsetenv("MYSQL_PORT")
		os.Unsetenv("MYSQL_DATABASE")
		os.Unsetenv("GO_ENV")
		os.Unsetenv("SERVER_ADDRESS")
	}()

	// Creating a new configuration to test environment variable integration.
	appConfig := NewConfig()

	// Validating the config values.
	expectedDSN := "test_user:test_password@tcp(test_host:3306)/test_database?parseTime=true"
	if appConfig.DataSourceName != expectedDSN || appConfig.Environment != "test_go_env" || appConfig.ServerAddress != "test_server_addr" {
		t.Errorf("NewConfig() = %+v; want DataSourceName=%s, Environment=%s, ServerAddress=%s",
			appConfig, expectedDSN, "test_go_env", "test_server_addr")
	}

	// Validate the default CacheMaxAge value (should be 5400 if CACHE_MAX_AGE is not set).
	if appConfig.CacheMaxAge != 5400 {
		t.Errorf("NewConfig() CacheMaxAge = %d; want 5400", appConfig.CacheMaxAge)
	}
}

// TestGetDefaultSwaggerHost verifies the behavior of the GetDefaultSwaggerHost function.
// It tests the function's ability to return the default Swagger host based on the environment.
func TestGetDefaultSwaggerHost(t *testing.T) {
	// Case 1: Get Swagger host for "development" environment.
	if got := GetDefaultSwaggerHost("development"); got != "localhost:8080" {
		t.Errorf("GetDefaultSwaggerHost('development') = %s; want localhost:8080", got)
	}

	// Case 2: Get Swagger host for any other environment.
	if got := GetDefaultSwaggerHost("production"); got != "gophersignal.com" {
		t.Errorf("GetDefaultSwaggerHost('production') = %s; want gophersignal.com", got)
	}
}

// TestCacheMaxAge verifies that the CACHE_MAX_AGE environment variable is correctly integrated.
func TestCacheMaxAge(t *testing.T) {
	t.Run("valid cache value", func(t *testing.T) {
		os.Setenv("CACHE_MAX_AGE", "6000")
		defer os.Unsetenv("CACHE_MAX_AGE")
		cfg := NewConfig()
		if cfg.CacheMaxAge != 6000 {
			t.Errorf("Expected CacheMaxAge to be 6000, got %d", cfg.CacheMaxAge)
		}
	})

	t.Run("invalid cache value falls back", func(t *testing.T) {
		os.Setenv("CACHE_MAX_AGE", "invalid")
		defer os.Unsetenv("CACHE_MAX_AGE")
		cfg := NewConfig()
		if cfg.CacheMaxAge != 5400 {
			t.Errorf("Expected default CacheMaxAge of 5400 on invalid input, got %d", cfg.CacheMaxAge)
		}
	})

	t.Run("default cache value", func(t *testing.T) {
		os.Unsetenv("CACHE_MAX_AGE")
		cfg := NewConfig()
		if cfg.CacheMaxAge != 5400 {
			t.Errorf("Expected default CacheMaxAge of 5400 when not set, got %d", cfg.CacheMaxAge)
		}
	})
}