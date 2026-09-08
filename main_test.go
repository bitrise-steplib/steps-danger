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

func Test_trimURLScheme(t *testing.T) {
	tests := []struct {
		name string
		url  string
		want string
	}{
		{
			name: "https scheme is removed",
			url:  "https://github.com/bitrise-io/sample-apps-ios-simple-objc.git",
			want: "github.com/bitrise-io/sample-apps-ios-simple-objc.git",
		},
		{
			name: "ssh scheme is removed",
			url:  "ssh://git@github.com/bitrise-io/sample-apps-ios-simple-objc.git",
			want: "git@github.com/bitrise-io/sample-apps-ios-simple-objc.git",
		},
		{
			name: "ssh scheme with a port is removed",
			url:  "ssh://git@github.com:22/bitrise-io/sample-apps-ios-simple-objc.git",
			want: "git@github.com:22/bitrise-io/sample-apps-ios-simple-objc.git",
		},
		{
			name: "an scp style URL has no scheme to remove",
			url:  "git@github.com:bitrise-io/sample-apps-ios-simple-objc.git",
			want: "git@github.com:bitrise-io/sample-apps-ios-simple-objc.git",
		},
		{
			// The regression: these hosts begin with characters that were in the old cutset, so
			// part of the host name was cut off with the scheme.
			name: "a host starting with a scheme character keeps its name",
			url:  "https://tools.corp.com/bitrise-io/repo.git",
			want: "tools.corp.com/bitrise-io/repo.git",
		},
		{
			name: "a host starting with s keeps its name",
			url:  "https://ssh.example.com/bitrise-io/repo.git",
			want: "ssh.example.com/bitrise-io/repo.git",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.want, trimURLScheme(tt.url))
		})
	}
}
