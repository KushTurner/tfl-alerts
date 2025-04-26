package config

import (
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestGetEnv(t *testing.T) {

	t.Run("Will use env if env exists", func(t *testing.T) {
		os.Setenv("ENVIRONMENT", "not_fallback")
		envVar := getEnv("ENVIRONMENT", "fallback")

		assert.Equal(t, "not_fallback", envVar)
	})

	t.Run("Will use fallback if env doesn't exist", func(t *testing.T) {
		os.Clearenv()
		envVar := getEnv("ENVIRONMENT", "fallback")

		assert.Equal(t, "fallback", envVar)
	})
}
