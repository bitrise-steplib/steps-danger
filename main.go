package main

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/Masterminds/semver"
	"github.com/bitrise-io/go-steputils/stepconf"
	"github.com/bitrise-io/go-steputils/v2/ruby"
	"github.com/bitrise-io/go-utils/log"
	"github.com/bitrise-io/go-utils/v2/command"
	"github.com/bitrise-io/go-utils/v2/env"
	v2log "github.com/bitrise-io/go-utils/v2/log"
	"github.com/kballard/go-shellquote"
)

// Config ...
type Config struct {
	RepositoryURL     string `env:"repository_url,required"`
	AdditionalOptions string `env:"additional_options"`

	GithubAPIToken   stepconf.Secret `env:"github_api_token"`
	GithubHost       string          `env:"github_host"`
	GithubAPIBaseURL string          `env:"github_api_base_url"`

	GitlabAPIToken   stepconf.Secret `env:"gitlab_api_token"`
	GitlabHost       string          `env:"gitlab_host"`
	GitlabAPIBaseURL string          `env:"gitlab_api_base_url"`
}

func validateInputs(cfg Config) {
	if cfg.GithubAPIToken == "" && cfg.GitlabAPIToken == "" {
		failf("None of the API tokens have been set.  If you want to use GitHub you need to set github_api_token. If you want to use GitLab you need to set gitlab_api_token")
	}

	// GitHub enterprise
	if (cfg.GithubHost != "" || cfg.GithubAPIBaseURL != "") && (cfg.GithubHost == "" || cfg.GithubAPIBaseURL == "") {
		failf("If you want to use GitHub Enterprise you need to set both of the github_host and the github_api_base_url")
	}

	// GitLab enterprise
	if (cfg.GitlabHost != "" || cfg.GitlabAPIBaseURL != "") && (cfg.GitlabHost == "" || cfg.GitlabAPIBaseURL == "") {
		failf("If you want to use GitLab Enterprise you need to set both of the gitlab_host and the gitlab_api_base_url")
	}

}

func failf(format string, v ...interface{}) {
	log.Errorf(format, v...)
	os.Exit(1)
}

// stdOpts returns the command options the Step's commands share.
func stdOpts() *command.Opts {
	return &command.Opts{
		Stdout: os.Stdout,
		Stderr: os.Stderr,
	}
}

// getBundlerVersion returns the bundler version the gem lockfile in searchDir was created with.
// A missing or unreadable lockfile is not an error: bundler is then installed and invoked without a
// version selector, which is what this Step did before too.
func getBundlerVersion(searchDir string) (ruby.Version, error) {
	lockFileContent, err := ruby.GemFileLockContent(searchDir)
	if err != nil {
		if errors.Is(err, ruby.ErrGemLockNotFound) {
			log.Warnf("No gem lockfile found in %s", searchDir)
		} else {
			log.Warnf("Could not read the gem lockfile, error: %s", err)
		}
		log.Infof("Using unspecified bundler version")

		return ruby.Version{}, nil
	}

	return ruby.ParseBundlerVersion(lockFileContent)
}

func main() {
	var cfg Config
	if err := stepconf.Parse(&cfg); err != nil {
		failf("Issue with input: %s", err)
	}

	logger := v2log.NewLogger()
	envRepository := env.NewRepository()
	cmdFactory := command.NewFactory(envRepository)
	cmdLocator := env.NewCommandLocator()

	// This Step runs danger through bundler, so unlike Steps that can fall back to an already
	// installed executable, it cannot do anything without Ruby.
	rubyFactory, err := ruby.NewCommandFactory(cmdFactory, cmdLocator, logger)
	if err != nil {
		failf("Failed to check the Ruby installation: %s", err)
	}
	rubyEnvironment := ruby.NewEnvironment(rubyFactory, cmdLocator, logger)

	cfg.RepositoryURL = trimScheme(cmdFactory, cfg.RepositoryURL)

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

	bundlerVersion, err := getBundlerVersion(".")
	if err != nil {
		failf("Could not determine required bundler version, error: %s", err)
	}

	if ok, err := rubyEnvironment.IsGemInstalled("bundler", bundlerVersion.Version); err != nil {
		failf("Failed to check bundler, error: %s", err)
	} else if !ok {
		log.Warnf(`Bundler is not installed`)
		fmt.Println()
		log.Printf("Installing Bundler")

		// force = true: in some configurations `bundler _1.2.3_` reports "Command not found" until
		// bundler is reinstalled.
		installBundlerCommands := rubyFactory.CreateGemInstall("bundler", bundlerVersion.Version, false, true, stdOpts())

		for _, installBundlerCommand := range installBundlerCommands {
			log.Donef("$ %s", installBundlerCommand.PrintableCommandArgs())
			fmt.Println()

			if err := installBundlerCommand.Run(); err != nil {
				failf("command failed, error: %s", err)
			}
		}
	}
	log.Printf("Bundler installed")

	//
	// Danger
	fmt.Println()
	log.Infof("Installing dependencies from your gem file")

	cmd := rubyFactory.CreateBundleInstall(bundlerVersion.Version, stdOpts())
	log.Printf("$ %s", cmd.PrintableCommandArgs())

	if err := cmd.Run(); err != nil {
		failf("Failed to run bundle install, error: %s", err)
	}

	fmt.Println()
	log.Infof("Running danger")

	additionalOptions, err := shellquote.Split(cfg.AdditionalOptions)
	if err != nil {
		failf("Failed to shell-quote additional options (%s): %s", cfg.AdditionalOptions, err)
	}

	cmd = rubyFactory.CreateBundleExec("danger", additionalOptions, bundlerVersion.Version, stdOpts())
	log.Printf("$ %s", cmd.PrintableCommandArgs())

	if err := cmd.Run(); err != nil {
		failf("Failed to run bundle exec danger, error: %s", err)
	}

	fmt.Println()
	log.Donef("Done")
}

// trimScheme trims the URL if danger version is <8.0.5
func trimScheme(cmdFactory command.Factory, url string) string {
	cmd := cmdFactory.Create("danger", []string{"--version"}, nil)
	log.Printf("$ %s", cmd.PrintableCommandArgs())

	dangerVersion, err := cmd.RunAndReturnTrimmedCombinedOutput()
	if err != nil {
		log.Errorf("Could not determine danger vesion: %s", err)
		return url
	}

	log.Printf("Found danger version: %s", dangerVersion)

	if shouldTrimScheme(dangerVersion) {
		return trimURLScheme(url)
	}

	return url
}

// trimURLScheme removes the scheme from a repository URL, if it has one.
//
// It used to be strings.TrimLeft(url, "https://"), whose second argument is a *cutset*, not a
// prefix: every leading character in "htps:/" was removed, so a host starting with one of those
// letters lost part of its name -- https://tools.corp.com became ools.corp.com. That went
// unnoticed because github.com starts with a letter outside the cutset. It also means the old code
// stripped ssh:// as well, which the tests below keep, so scm URLs without a scheme
// (git@github.com:owner/repo.git) are still returned untouched.
func trimURLScheme(url string) string {
	const schemeSeparator = "://"

	if i := strings.Index(url, schemeSeparator); i >= 0 {
		return url[i+len(schemeSeparator):]
	}

	return url
}

func shouldTrimScheme(rawDangerVersion string) bool {
	dangerVersion, err := semver.NewVersion(rawDangerVersion)
	if err != nil {
		log.Errorf("Could not parse danger vesion: %s", err)
		return false
	}

	versionConstraint, err := semver.NewConstraint("<8.0.5")
	if err != nil {
		log.Errorf("Could not parse version constraint: %s", err)
		return false
	}

	return versionConstraint.Check(dangerVersion)
}
