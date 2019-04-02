package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/bitrise-io/go-utils/command"
	"github.com/bitrise-io/go-utils/command/rubycommand"
	"github.com/bitrise-io/go-utils/errorutil"
	"github.com/bitrise-io/go-utils/log"
	"github.com/bitrise-tools/go-steputils/stepconf"
)

// Config ...
type Config struct {
	RepositoryURL string `env:"repository_url,required"`

	GithubAPIToken   stepconf.Secret `env:"github_api_token"`
	GithubHost       string          `env:"github_host"`
	GithubAPIBaseURL string          `env:"github_api_base_url"`

	GitlabAPIToken   stepconf.Secret `env:"gitlab_api_token"`
	GitlabHost       string          `env:"gitlab_host"`
	GitlabAPIBaseURL string          `env:"gitlab_api_base_url"`
}

func main() {
	var cfg Config
	if err := stepconf.Parse(&cfg); err != nil {
		failf("Issue with input: %s", err)
	}

	// Fix repository URL
	cfg.RepositoryURL = strings.TrimLeft(cfg.RepositoryURL, "https://")

	stepconf.Print(cfg)
	fmt.Println()

	validateInputs(cfg)

	//
	// Set local envs for the step
	for key, value := range map[string]string{
		"GIT_REPOSITORY_URL":         cfg.RepositoryURL,
		"DANGER_GITHUB_API_TOKEN":    string(cfg.GithubAPIToken),
		"DANGER_GITHUB_HOST":         cfg.GithubHost,
		"DANGER_GITHUB_API_BASE_URL": cfg.GithubAPIBaseURL,
		"DANGER_GITLAB_API_TOKEN":    string(cfg.GitlabAPIToken),
		"DANGER_GITLAB_HOST":         cfg.GitlabHost,
		"DANGER_GITLAB_API_BASE_URL": cfg.GitlabAPIBaseURL,
	} {
		if value != "" {
			if err := os.Setenv(key, value); err != nil {
				failf("Failed to set env %s, error: %s", key, err)
			}
		}
	}

	//
	// Check dependencies
	log.Infof("Checking dependencies")
	log.Printf("Bundler...")
	if ok, err := rubycommand.IsGemInstalled("bundler", ""); err != nil {
		failf("Failed to check bundler, error: %s", err)
	} else if !ok {
		log.Warnf(`Bundler is not installed`)
		fmt.Println()
		log.Printf("Install Bundler")

		cmds, err := rubycommand.GemInstall("bundler", "")
		if err != nil {
			failf("failed to create command model, error: %s", err)
		}

		for _, cmd := range cmds {
			if out, err := cmd.RunAndReturnTrimmedCombinedOutput(); err != nil {
				if errorutil.IsExitStatusError(err) {
					failf("%s failed: %s", cmd.PrintableCommandArgs(), out)
				}
				failf("%s failed: %s", cmd.PrintableCommandArgs(), err)
			}
		}
	}
	log.Printf("Bundler installed")

	//
	// Danger
	fmt.Println()
	log.Infof("Installing dependencies from your gem file")

	cmd := command.New("bundle", "install")
	cmd.SetStdout(os.Stdout)
	cmd.SetStderr(os.Stderr)
	log.Printf("$ %s", cmd.PrintableCommandArgs())

	if err := cmd.Run(); err != nil {
		failf("Failed to run bundle install, error: %s", err)
	}

	fmt.Println()
	log.Infof("Running danger")

	cmd = command.New("bundle", "exec", "danger")
	cmd.SetStdout(os.Stdout)
	cmd.SetStderr(os.Stderr)
	log.Printf("$ %s", cmd.PrintableCommandArgs())

	if err := cmd.Run(); err != nil {
		failf("Failed to run bundle exec danger, error: %s", err)
	}

	fmt.Println()
	log.Donef("Done")
}

func validateInputs(cfg Config) {
	if cfg.GithubAPIToken == "" && cfg.GitlabAPIToken == "" {
		failf("None of the API token have been set.  If you want to use Github you need to set github_api_token. If you want to use Gitlab you need to set gitlab_api_token")
	}

	// Github enterprise
	if (cfg.GithubHost != "" || cfg.GithubAPIBaseURL != "") && (cfg.GithubHost == "" || cfg.GithubAPIBaseURL == "") {
		failf("If you want to use Github Enterprise you need to set both of the github_host and the github_api_base_url")
	}

	// Gitlab enterprise
	if (cfg.GitlabHost != "" || cfg.GitlabAPIBaseURL != "") && (cfg.GitlabHost == "" || cfg.GitlabAPIBaseURL == "") {
		failf("If you want to use Github Enterprise you need to set both of the gitlab_host and the gitlab_api_base_url")
	}

}

func failf(format string, v ...interface{}) {
	log.Errorf(format, v...)
	os.Exit(1)
}
