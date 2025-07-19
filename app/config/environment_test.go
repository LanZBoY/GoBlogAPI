package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetEnvOrDefault(t *testing.T) {
	os.Setenv("TEST_ENV", "value")
	defer os.Unsetenv("TEST_ENV")
	assert.Equal(t, "value", getEnvOrDefault("TEST_ENV", "default"))
	assert.Equal(t, "default", getEnvOrDefault("NOT_SET", "default"))
}
