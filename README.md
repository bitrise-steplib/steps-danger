# Run Danger

[![Step changelog](https://shields.io/github/v/release/bitrise-steplib/steps-danger?include_prereleases&label=changelog&color=blueviolet)](https://github.com/bitrise-steplib/steps-danger/releases)

[Danger](https://danger.systems/python/) automates common code review chores during your CI process.

<details>
<summary>Description</summary>

If inserted in your Workflow on Bitrise, it runs the Danger CI to check any linting issues in your project's code.

### Configuring the Step

1. The **Repository URL of your project** input is automatically filled out.
2. If you add any additional options in the **Additional options for the command call**, they will be added to your `bundle exec danger` command call.
3. Select a git provider's input section: GitHub or GitLab.
4. If you're using GitHub:
- Add your access token in the **Access token for your project** input. Click the input's description for more information on how to set up the access token.
- Add the host GitHub is running on, for example, `git.corp.evilcorp.com`. Read more about [how to set it up](https://danger.systems/guides/getting_started.html).
- Add the GitHub API Enterprise API URL in the **GitHub API base URL** input.
5. If you are using GitLab:
- Add your access token in the **Access token for your project** input. Click the input's description for more information on how to set up the access token.
- Add the host GitLab is running on in the **GitLab host** input. You must add this if you are using Self-Managed GitLab.
- Add the **GitLab API base URL**. You must add this if you are using Self-Managed GitLab.

### Useful links
- [No activity summaries found for test with Danger](https://devcenter.bitrise.io/troubleshooting/no-activity-summaries-found-for-test-with-danger/#the-issue)
</details>

## 🧩 Get started

Add this step directly to your workflow in the [Bitrise Workflow Editor](https://docs.bitrise.io/en/bitrise-ci/workflows-and-pipelines/steps/adding-steps-to-a-workflow.html).

You can also run this step directly with [Bitrise CLI](https://github.com/bitrise-io/bitrise).

## ⚙️ Configuration

<details>
<summary>Inputs</summary>

| Key | Description | Flags | Default |
| --- | --- | --- | --- |
| `repository_url` | Repository URL of your project | required | `$GIT_REPOSITORY_URL` |
| `github_api_token` | **SETTING UP AN ACCESS TOKEN**  [Here’s the link](https://github.com/settings/tokens/new), you should open this in the private session where you just created the new GitHub account. Again, the rights that you give to the token depend on the openness of your projects. You’ll want to save for later, when you add a `github_api_token` to the step.  **TOKENS FOR OSS PROJECTS**  We recommend giving the token the smallest scope possible. This means just public\_repo, this scopes limits Danger’s abilities to just writing comments on OSS projects. Because the token can be quite easily be extracted from the CI environment, this minimizes the chance for bad actors to cause chaos with it.  **TOKENS FOR CLOSED SOURCE PROJECTS**  We recommend giving access to the whole repo scope, and its children.  **You can read more about it here:** [https://danger.systems/guides/getting_started.html](https://danger.systems/guides/getting_started.html) | sensitive |  |
| `github_host` | The host that GitHub is running on. You need to set this if you are using **Enterprise GitHub**. You can work with GitHub Enterprise by setting the `github_host` and the `github_api_base_url` inputs.  **For example:** `git.corp.evilcorp.com`  **You can read more about it here:** [https://danger.systems/guides/getting_started.html](https://danger.systems/guides/getting_started.html) |  |  |
| `github_api_base_url` | The host that the GitHub Enterprise API is reachable on. You need to set this if you are using **Enterprise GitHub**. You can work with GitHub Enterprise by setting the `github_host` and the `github_api_base_url` inputs.  **For example:** `https://git.corp.evilcorp.com/api/v3`  **You can read more about it here:** [https://danger.systems/guides/getting_started.html](https://danger.systems/guides/getting_started.html) |  |  |
| `gitlab_api_token` | **SETTING UP AN ACCESS TOKEN**  Here’s the link, you should open this in the private session where you have just created the new GitLab account. You’ll want to copy the token for later, when you add a `gitlab_api_token` to the step.  If you are self hosting GitLab, you’ll need to generate an access token for the bot user. You can find the section under “Access Tokens” in the bot user’s profile.  Find more information about Danger in their guides: [https://danger.systems/guides/getting_started.html](https://danger.systems/guides/getting_started.html) | sensitive |  |
| `gitlab_host` | The host that GitLab is running on. You need to set this if you are using **Self-Managed GitLab**. You can work with Self-Managed GitLab by setting the `gitlab_host` and the `gitlab_api_base_url` inputs.  **For example:** `git.corp.evilcorp.com`  **You can read more about it here:** [https://danger.systems/guides/getting_started.html](https://danger.systems/guides/getting_started.html) |  |  |
| `gitlab_api_base_url` | The host that the Self-Managed GitLab API is reachable on. You need to set this if you are using **Self-Managed GitLab**. You can work with Self-Managed GitLab by setting the `gitlab_host` and the `gitlab_api_base_url` inputs.  **For example:** `https://git.corp.evilcorp.com/api/v4`  **You can read more about it here:** [https://danger.systems/guides/getting_started.html](https://danger.systems/guides/getting_started.html) |  |  |
| `additional_options` | Additional commands and options to append to the danger command call. The provided value will be appended to the `bundle exec danger` command call, as is. |  | `--fail-on-errors=true` |
</details>

<details>
<summary>Outputs</summary>
There are no outputs defined in this step
</details>

## 🙋 Contributing

We welcome [pull requests](https://github.com/bitrise-steplib/steps-danger/pulls) and [issues](https://github.com/bitrise-steplib/steps-danger/issues) against this repository.

For pull requests, work on your changes in a forked repository and use the Bitrise CLI to [run step tests locally](https://docs.bitrise.io/en/bitrise-ci/bitrise-cli/running-your-first-local-build-with-the-cli.html).

Learn more about developing steps:

- [Create your own step](https://docs.bitrise.io/en/bitrise-ci/workflows-and-pipelines/developing-your-own-bitrise-step/developing-a-new-step.html)
