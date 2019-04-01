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
	if err := os.Setenv("GIT_REPOSITORY_URL", cfg.RepositoryURL); err != nil {
		failf("Failed to set env GIT_REPOSITORY_URL, error: %s", err)
	}

	// Github
	if string(cfg.GithubAPIToken) != "" {
		if err := os.Setenv("DANGER_GITHUB_API_TOKEN", string(cfg.GithubAPIToken)); err != nil {
			failf("Failed to set env DANGER_GITHUB_API_TOKEN, error: %s", err)
		}
	}

	if string(cfg.GithubHost) != "" {
		if err := os.Setenv("DANGER_GITHUB_HOST", string(cfg.GithubHost)); err != nil {
			failf("Failed to set env DANGER_GITHUB_HOST, error: %s", err)
		}
	}

	if string(cfg.GithubAPIBaseURL) != "" {
		if err := os.Setenv("DANGER_GITHUB_API_BASE_URL", string(cfg.GithubAPIBaseURL)); err != nil {
			failf("Failed to set env DANGER_GITHUB_API_BASE_URL, error: %s", err)
		}
	}

	//
	// Gitlab
	if string(cfg.GitlabAPIToken) != "" {
		if err := os.Setenv("DANGER_GITLAB_API_TOKEN", string(cfg.GitlabAPIToken)); err != nil {
			failf("Failed to set env DANGER_GITLAB_API_TOKEN, error: %s", err)
		}
	}

	if string(cfg.GitlabHost) != "" {
		if err := os.Setenv("DANGER_GITLAB_HOST", string(cfg.GitlabHost)); err != nil {
			failf("Failed to set env DANGER_GITLAB_HOST, error: %s", err)
		}
	}

	if string(cfg.GitlabAPIBaseURL) != "" {
		if err := os.Setenv("DANGER_GITLAB_API_BASE_URL", string(cfg.GitlabAPIBaseURL)); err != nil {
			failf("Failed to set env DANGER_GITLAB_API_BASE_URL, error: %s", err)
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
					log.Warnf("%s failed: %s", out)
				} else {
					log.Warnf("%s failed: %s", err)
				}
				failf("Failed to install Bundler, error: %s", err)
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

func installBundler() ([]*command.Model, error) {
	cmds, err := rubycommand.GemInstall("bundler", "")
	if err != nil {
		return nil, fmt.Errorf("failed to create command model, error: %s", err)
	}

	return cmds, nil
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
