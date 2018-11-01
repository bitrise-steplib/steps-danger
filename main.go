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
	RepositoryURL  string          `env:"repository_url,required"`
	GithubAPIToken stepconf.Secret `env:"github_api_token"`
	GitlabAPIToken stepconf.Secret `env:"gitlab_api_token"`
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

	os.Setenv("GIT_REPOSITORY_URL", cfg.RepositoryURL)

	if string(cfg.GithubAPIToken) != "" {
		os.Setenv("DANGER_GITHUB_API_TOKEN", string(cfg.GithubAPIToken))
	}

	if string(cfg.GitlabAPIToken) != "" {
		os.Setenv("DANGER_GITLAB_API_TOKEN", string(cfg.GitlabAPIToken))
	}

	//
	// Check dependencies
	fmt.Println()
	log.Infof("Checking dependencies")
	{
		log.Printf("Bundler...")
		if ok, err := bundlerInstalled(); err != nil {
			failf("Failed to check bundler, error: %s", err)
		} else if !ok {
			log.Warnf(`Bundler is not installed`)
			fmt.Println()
			log.Printf("Install Bundler")

			if cmds, err := installBundler(); err != nil {
				failf("Failed to create Bundler install command, error: %s", err)
			} else {
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
		}
	}

	// Danger
	{
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
	}

	fmt.Println()
	log.Donef("Done")

}

func bundlerInstalled() (bool, error) {
	return rubycommand.IsGemInstalled("bundler", "")
}

// Install ...
func installBundler() ([]*command.Model, error) {
	cmds, err := rubycommand.GemInstall("bundler", "")
	if err != nil {
		return nil, fmt.Errorf("failed to create command model, error: %s", err)
	}

	return cmds, nil
}

func failf(format string, v ...interface{}) {
	log.Errorf(format, v...)
	os.Exit(1)
}
