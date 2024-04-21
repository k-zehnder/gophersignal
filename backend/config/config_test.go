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
// It sets up test environment variables, calls NewConfig to create a configuration, and validates the correctness of the config fields.
func TestNewConfig(t *testing.T) {
	// Preparing test environment variables.
	os.Setenv("MYSQL_DSN", "test_mysql_dsn")
	os.Setenv("GO_ENV", "test_go_env")
	os.Setenv("SERVER_ADDRESS", "test_server_addr")
	defer func() {
		os.Unsetenv("MYSQL_DSN")
		os.Unsetenv("GO_ENV")
		os.Unsetenv("SERVER_ADDRESS")
	}()

	// Creating a new configuration to test environment variable integration.
	appConfig := NewConfig()

	// Validating the config values.
	if appConfig.DataSourceName != "test_mysql_dsn" || appConfig.Environment != "test_go_env" || appConfig.ServerAddress != "test_server_addr" {
		t.Error("NewConfig() did not set config fields correctly")
	}
}

// TestGetDefaultSwaggerHost verifies the behavior of the GetDefaultSwaggerHost function.
// It tests the function's ability to return the default Swagger host based on the environment.
func TestGetDefaultSwaggerHost(t *testing.T) {
	// Case 1: Get Swagger host for "dev" environment.
	if got := GetDefaultSwaggerHost("dev"); got != "localhost:8080" {
		t.Errorf("GetDefaultSwaggerHost('dev') = %s; want localhost:8080", got)
	}

	// Case 2: Get Swagger host for any other environment.
	if got := GetDefaultSwaggerHost("prod"); got != "gophersignal.com" {
		t.Errorf("GetDefaultSwaggerHost('prod') = %s; want gophersignal.com", got)
	}
}
