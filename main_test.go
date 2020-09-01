package main

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_ShouldTrimScheme(t *testing.T) {
	scenarios := []struct {
		input    string
		expected bool
	}{
		{"8", true},
		{"8.0", true},
		{"8.0.0", true},
		{"8.0.4", true},
		{"8.0.5", false},
		{"8.0.6", false},
		{"8.1.0", false},
		{"9", false},
	}

	for _, scenario := range scenarios {
		acutalResult := shouldTrimScheme(scenario.input)
		require.Equal(t, scenario.expected, acutalResult)
	}
}
