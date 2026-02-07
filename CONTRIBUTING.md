# How to contribute

**First:** if you're unsure or afraid of _anything_, ask for help! You can
submit a work in progress (WIP) pull request, or file an [issue](https://github.com/bpg/terraform-provider-proxmox/issues) with the parts
you know. For questions and general discussion, use [GitHub Discussions](https://github.com/bpg/terraform-provider-proxmox/discussions).
We'll do our best to guide you in the right direction, and let you
know if there are guidelines we will need to follow. We want people to be able
to participate without fear of doing the wrong thing.

Below are our expectations for contributors. Following these guidelines gives us
the best opportunity to work with you, by making sure we have the things we need
in order to make it happen. Doing your best to follow it will speed up our
ability to merge PRs and respond to issues.

## Quick start

1. Fork and clone the repository.
2. Create a feature branch.
3. Implement your changes.
4. Run `make test` and `make lint`.
5. Commit with `git commit -s` (signs off on DCO).
6. Open a PR against `main`.

For common issues during development, see [Development Troubleshooting](docs/guides/dev-troubleshooting.md).

## Table of contents

- [Quick start](#quick-start)
- [Build the provider](#build-the-provider)
- [IDE support](#ide-support)
- [Devcontainer support](#devcontainer-support)
- [Testing](#testing)
- [Provider implementation guidance](#provider-implementation-guidance)
- [Coding conventions](#coding-conventions)
- [Commit message conventions](#commit-message-conventions)
- [Developer Certificate of Origin (DCO)](#developer-certificate-of-origin-dco)
- [Submitting changes](#submitting-changes)
- [Using AI assistants and LLM agents](#using-ai-assistants-and-llm-agents)
- [Releasing](#releasing)

## Build the provider

> [!TIP]
> `$GOPATH` is the path to your Go workspace. If undefined, it defaults to `$HOME/go` on Linux and macOS, and `%USERPROFILE%\go` on Windows.

> [!NOTE]
> The provider requires Go 1.25 or later to build.

- Clone the repository to: `$GOPATH/src/github.com/bpg/terraform-provider-proxmox`:

  ```sh
  mkdir -p "${GOPATH}/src/github.com/bpg"
  cd "${GOPATH}/src/github.com/bpg"
  git clone git@github.com:bpg/terraform-provider-proxmox
  ```

- Enter the provider directory and build it:

  ```sh
  cd "${GOPATH}/src/github.com/bpg/terraform-provider-proxmox"
  make build
  ```

- To cross-compile the provider for all supported platforms:

  ```sh
  make release-build
  ```

  The compiled binaries will be placed in the `dist` directory.

- A portion of the documentation is generated from the source code. To update the documentation, run:

  ```sh
  make docs
  ```

## IDE support

If you are using VS Code, feel free to copy `settings.json` from `.vscode/settings.example.json`.

## Devcontainer support

Prerequisites:

- Docker (or Docker Desktop) installed on your machine
- [VS Code Remote - Containers extension](https://marketplace.visualstudio.com/items?itemName=ms-vscode-remote.remote-containers)

To launch the devcontainer:

1. Open the project in VS Code.
2. Run **Remote-Containers: Open Folder in Container** from the Command Palette.

See [Developing inside a Container](https://code.visualstudio.com/docs/devcontainers/containers) for more details.

## Testing

### Unit tests

The project has a test suite that must pass for contributions to be accepted. When making changes:

1. Run all tests with:

   ```sh
   make test
   ```

2. Add or modify test cases to cover your changes.
3. Ensure all tests pass before submitting your PR.

### Acceptance tests

Acceptance tests run against a real Proxmox instance and verify the provider's functionality end-to-end.

**Where to put tests:** place acceptance tests alongside the resource or data source implementation (same package/folder) whenever possible. The shared `fwprovider/test/` directory is reserved for cross-resource integration tests or cases that require expensive Proxmox setup shared across multiple suites.

#### Prerequisites

1. A running Proxmox instance (see [Development Proxmox Setup](docs/guides/dev-proxmox-setup.md))
2. Create a `testacc.env` file in the project root with:

   ```env
   TF_ACC=1
   PROXMOX_VE_API_TOKEN="root@pam!<token name>=<token value>"
   PROXMOX_VE_ENDPOINT="https://<pve instance>:8006/"
   PROXMOX_VE_SSH_AGENT="true"
   PROXMOX_VE_SSH_USERNAME="root"
   ```

   Optional configuration:

   ```env
   # Override default node name and SSH settings
   PROXMOX_VE_ACC_NODE_NAME="pve1"
   PROXMOX_VE_ACC_NODE_SSH_ADDRESS="10.0.0.11"
   PROXMOX_VE_ACC_NODE_SSH_PORT="22"
   PROXMOX_VE_ACC_IFACE_NAME="enp1s0"
   ```

#### Running acceptance tests

Run the full acceptance test suite with:

```sh
make testacc
```

To run a single test or a group of tests, use the helper script:

```sh
./testacc <test_name>
```

For example, to run all VM-related tests: `./testacc TestAccResourceVM.*`

> [!NOTE]
>
> - Acceptance test coverage is still in development.
> - Only some resources and data sources are currently tested.
> - Some tests may require specific Proxmox configuration.

### Manual testing

You can test the provider locally before submitting changes:

1. Create a provider override configuration in one of these locations:

   **Linux/macOS** (`$HOME/.terraformrc`):

   ```hcl
   provider_installation {
     dev_overrides {
       "bpg/proxmox" = "/home/user/go/bin/"  # Replace with your $GOPATH/bin
     }
     direct {}
   }
   ```

   **Windows** (`%APPDATA%/terraform.rc`):

   ```hcl
   provider_installation {
     dev_overrides {
       "bpg/proxmox" = "C:\\Users\\user\\go\\bin"  # Replace with your %GOPATH%/bin
     }
     direct {}
   }
   ```

2. Build and install the provider:

   ```sh
   go install .
   ```

3. Test your changes:

   ```sh
   terraform plan   # Preview changes
   terraform apply  # Apply changes
   ```

> [!TIP]
> After the initial setup, you only need to run `go install .` when rebuilding the provider.

## Provider implementation guidance

New resources and data sources **must** be implemented using the Terraform Plugin Framework. The framework provider lives under the `fwprovider/` directory. The legacy SDK implementation in `proxmoxtf/` is feature-frozen; PRs that add new SDK-based resources or data sources will not be accepted.

### Reference implementations

See [docs/adr/reference-examples.md](docs/adr/reference-examples.md) for annotated walkthroughs and a new-resource checklist. Quick summary:

| Complexity | Reference | When to use |
| ---------- | --------- | ----------- |
| Basic CRUD | SDN VNet (`fwprovider/cluster/sdn/vnet/`) | Start here for any new resource |
| Many optional fields | Metrics Server (`fwprovider/cluster/metrics/`) | Sensitive attributes, bool-to-int conversion |
| Cross-field validation | ACL (`fwprovider/access/`) | `ConfigValidators`, custom import ID parsing |

Architecture decisions are documented in [docs/adr/](docs/adr/README.md).

### Documentation workflow

Documentation for Framework resources is **auto-generated** from schema definitions:

1. **Write descriptive schema attributes** — Add `Description` and/or `MarkdownDescription` fields to your schema attributes. These become the docs.
2. **Optional: Create a template** — For custom formatting, create a template in `/templates/resources/` or `/templates/data-sources/`.
3. **Add a `go:generate` directive** — For new resources/data sources, add a `cp` command in `main.go` to copy the generated doc from `build/docs-gen/` to `docs/`.
4. **Run `make docs`** — This generates documentation in `/docs/` from schemas and templates.

> [!IMPORTANT]
> Do **not** manually edit files in `/docs/` for Framework resources. Your changes will be overwritten by `make docs`. Edit schema descriptions or templates instead.

**Description vs MarkdownDescription:**

| Field                   | Format     | Used for                                      |
| ----------------------- | ---------- | --------------------------------------------- |
| `Description`           | Plain text | CLI help, simple tooltips                     |
| `MarkdownDescription`   | Markdown   | Registry docs, rich formatting                |

- If only one is set, it's used for both purposes.
- Use `MarkdownDescription` when you need inline code (backticks), links, or HTML (`<br>`).
- Keep descriptions concise: explain what the attribute does, valid values, and defaults.

Example:

```go
schema.StringAttribute{
    Description:         "The name of the VM.",
    MarkdownDescription: "The name of the VM. Must be a valid DNS name (`[a-zA-Z0-9-]+`).",
    Optional:            true,
}
```

**When to use templates:**

- Adding usage examples beyond auto-generated ones
- Custom warnings, notes, or formatting
- Import instructions with specific syntax

See existing templates in `/templates/` for examples.

**Admonitions in docs:**

Use Terraform registry admonition syntax in `/docs/` files and templates:

| Symbol | Type    | Usage                           |
| ------ | ------- | ------------------------------- |
| `->`   | Note    | General information, tips       |
| `~>`   | Warning | Cautions, important caveats     |
| `!>`   | Danger  | Critical warnings, "do not use" |

Example:

```markdown
-> Consider using `proxmox_virtual_environment_download_file` resource instead.

~> Never commit proxy configurations or credentials to the repository.

!> **DO NOT USE** — This resource is experimental and will change.
```

For more details, see the [Terraform Plugin Framework documentation on descriptions](https://developer.hashicorp.com/terraform/plugin/framework/handling-data/attributes#description).

### Best practices

- Keep validation logic consistent across framework components.
- Place acceptance tests alongside the implementation (same package/folder).
- Reuse shared helpers from `fwprovider/attribute/` and `fwprovider/validators/`.

## Coding conventions

We expect all code contributions to follow these guidelines:

1. Code must pass linting with `golangci-lint`
   - Run `make lint` to format and lint your code (includes formatting via `golangci-lint fmt`).
   - The project uses `.golangci.yml` for linting configuration.

## Commit message conventions

We follow the [Conventional Commits](https://www.conventionalcommits.org/) specification. Please use the following types for your commits:

- `feat`: New features
- `fix`: Bug fixes
- `chore`: Maintenance tasks

These types are used to automatically generate the changelog. Other types will be ignored.

Use the `scope` field to indicate the area of the codebase being changed:

- `vm` – Virtual Machine resources
- `lxc` – Container resources
- `provider` – Provider configuration and resources
- `core` – Core libraries and utilities
- `docs` – Documentation
- `ci` – Continuous Integration / Actions / GitHub Workflows

Guidelines:

- Use lowercase for descriptions.
- Do not end descriptions with a period.
- Keep the first line under 72 characters.

Example:

```text
feat(vm): add support for the `clone` operation
```

## Developer Certificate of Origin (DCO)

All contributions must be signed off according to the Developer Certificate of Origin (DCO). The DCO is a lightweight way of certifying that you wrote or have the right to submit the code you are contributing. It provides legal protection for the project by ensuring contributors have the necessary rights to their contributions and agree to license them under the project's terms. You can find the full text [here](https://developercertificate.org).

To sign off your commits, add a `Signed-off-by` line to your commit message:

```text
feat(vm): add support for the `clone` operation

Signed-off-by: Random Developer <random@developer.example.org>
```

> [!NOTE]
>
> - **Name**: Use your real name (preferred) or GitHub username if you prefer privacy.
> - **Email**: Use a valid email address (GitHub's 'noreply' email is acceptable for privacy, see [GitHub docs](https://docs.github.com/en/account-and-profile/setting-up-and-managing-your-personal-account-on-github/managing-email-preferences/setting-your-commit-email-address#setting-your-commit-email-address-on-github)).
> - **Auto-sign**: If your Git config has `user.name` and `user.email` set, use `git commit -s` to automatically add the sign-off.

For more details about the DCO checker, see the [DCO app repo](https://github.com/dcoapp/app).

## Submitting changes

### Pull request scope

Please keep PRs small and focused. Small PRs are easier to review, easier to test, and get merged faster.

Guidelines:

- **One change per PR**: a single fix/feature, or a refactor with no behavior change.
- **Avoid multi-resource PRs**: if you need to change multiple resources/data sources, split the work into separate PRs (you can stack them and link follow-ups).
- **Do not mix concerns**: avoid combining formatting-only changes, refactors, and behavior changes in the same PR.
- **Iterate quickly**: open a draft/WIP PR early if you want feedback on approach before polishing edge cases.

### Proof of work

**Proof of work is mandatory for all code changes.** Every PR must include evidence that the change works as expected:

- Test output (unit tests, acceptance tests)
- Logs, screenshots, or terminal output demonstrating the fix/feature
- Any other relevant information that demonstrates the change works as expected

> [!WARNING]
> PRs without proof of work may be rejected. Trivial changes (typo fixes, documentation-only updates that don't affect code behavior) are exempt from this requirement.
If you use AI assistants, they are expected to generate a proof of work document as a `.dev/*_REPORT.md` file. Review this file and use its contents when completing the PR template.

### How to submit

1. Create a new PR against the `main` branch using the project's [pull request template](.github/PULL_REQUEST_TEMPLATE.md).
2. Ensure your PR title follows the Conventional Commits specification (we use this as the squash commit message).
3. **Include proof of work** in the PR description (test results, logs, screenshots).
4. All commits in a PR are typically squashed on merge.

## Using AI assistants and LLM agents

We welcome contributions that use AI assistants, LLM agents, or AI-powered coding tools. These tools can help with code generation, testing, documentation, and other development tasks.

### Guidelines

**Allowed and encouraged:**

- Using AI assistants (GitHub Copilot, Cursor, Claude, ChatGPT, etc.) to help write code
- Using LLM agents to automate repetitive tasks
- Leveraging AI for test generation, documentation, or debugging
- Any tool that helps you complete the task effectively

**Contributor responsibility:**

While AI tools can assist with contributions, **the person submitting the change is fully responsible** for:

1. **Code quality and correctness** — Review all AI-generated code carefully. You are accountable for what you submit.
2. **DCO sign-off** — By signing off (`git commit -s`), you personally certify that you have the right to submit the code under the project's license, regardless of how it was generated.
3. **Reproducible proof of work** — See [Proof of work](#proof-of-work) requirements above.
4. **Understanding the change** — Be prepared to explain and defend your contribution during code review.

> [!IMPORTANT]
> The DCO sign-off is a legal certification. When you sign off on a commit, you are affirming that you wrote or have the right to submit the code, and that you agree to license it under the project's terms. This applies equally to human-written and AI-assisted code.

### Agent instructions

For AI agents working on this repository:

- **[CLAUDE.md](CLAUDE.md)** — Development guidelines and critical rules
- **[GEMINI.md](GEMINI.md)** — PR review instructions
- **[.dev/README.md](.dev/README.md#working-with-llm-agents)** — Detailed workflow with skills (`/start-issue`, `/ready`, `/debug-api`, `/prepare-pr`, `/resume`)

The skills automate common workflows like setting up branches, running checklists, and preparing PR submissions.

## Releasing

We use [release-please](https://github.com/googleapis/release-please) GitHub Action for automated release management. The process works as follows:

1. The action creates a release PR based on commit messages.
2. The PR includes an auto-generated changelog and version bump.
3. Maintainers review and merge the release PR.
4. The release is automatically published to:
   - GitHub Releases
   - Terraform Registry

We aim to release new versions every 1–2 weeks.
