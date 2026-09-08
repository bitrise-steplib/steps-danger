package main

import (
	"os"
	"path/filepath"
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

func Test_getBundlerVersion(t *testing.T) {
	const lockWithBundler = `GEM
  remote: https://rubygems.org/
  specs:
    danger (8.0.5)

DEPENDENCIES
  danger

BUNDLED WITH
   2.4.12
`
	const lockWithoutBundler = `GEM
  remote: https://rubygems.org/
  specs:
    danger (8.0.5)

DEPENDENCIES
  danger
`

	tests := []struct {
		name string
		// lockFileName is the gem lockfile to write, or empty to write none.
		lockFileName string
		lockContent  string
		wantVersion  string
		wantFound    bool
	}{
		{
			// Not an error: bundler is then installed and invoked without a version selector.
			name: "no gem lockfile",
		},
		{
			name:         "Gemfile.lock naming a bundler version",
			lockFileName: "Gemfile.lock",
			lockContent:  lockWithBundler,
			wantVersion:  "2.4.12",
			wantFound:    true,
		},
		{
			name:         "Gemfile.lock without a BUNDLED WITH section",
			lockFileName: "Gemfile.lock",
			lockContent:  lockWithoutBundler,
		},
		{
			// gems.locked is the other name bundler accepts. The v1 helper this replaced only ever
			// looked for Gemfile.lock.
			name:         "gems.locked is honoured too",
			lockFileName: "gems.locked",
			lockContent:  lockWithBundler,
			wantVersion:  "2.4.12",
			wantFound:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			searchDir := t.TempDir()
			if tt.lockFileName != "" {
				require.NoError(t, os.WriteFile(filepath.Join(searchDir, tt.lockFileName), []byte(tt.lockContent), 0600))
			}

			got, err := getBundlerVersion(searchDir)

			require.NoError(t, err)
			require.Equal(t, tt.wantVersion, got.Version)
			require.Equal(t, tt.wantFound, got.Found)
		})
	}
}
