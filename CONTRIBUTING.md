# How to contribute

**First:** if you're unsure or afraid of _anything_, ask for help! You can
submit a work in progress (WIP) pull request, or file an issue with the parts
you know. We'll do our best to guide you in the right direction, and let you
know if there are guidelines we will need to follow. We want people to be able
to participate without fear of doing the wrong thing.

Below are our expectations for contributors. Following these guidelines gives us
the best opportunity to work with you, by making sure we have the things we need
in order to make it happen. Doing your best to follow it will speed up our
ability to merge PRs and respond to issues.

## Build the provider

> [!TIP]
> `$GOPATH` is the path to your Go workspace. If undefined, it defaults to `$HOME/go` on Linux and macOS, and `%USERPROFILE%\go` on Windows.

> [!NOTE]
> The provider requires Go 1.24 or later to build.

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
  make build-all
  ```

  The compiled binaries will be placed in the `dist` directory.

A portion of the documentation is generated from the source code. To update the documentation, run:

```sh
make docs
```

## IDE support

If you are using VS Code, feel free to copy `settings.json` from `.vscode/settings.example.json`.

## Devcontainer support

The project uses a devcontainer to provide a consistent development environment.
If you are using VS Code, you can use the devcontainer by opening the project in a container.
See [Developing inside a Container](https://code.visualstudio.com/docs/devcontainers/containers) for more details.

## Testing

### Unit Tests

The project has a test suite that must pass for contributions to be accepted. When making changes:

1. Run all tests with:

   ```sh
   make test
   ```

2. Add or modify test cases to cover your changes
3. Ensure all tests pass before submitting your PR

### Acceptance Tests

Acceptance tests run against a real Proxmox instance and verify the provider's functionality end-to-end. These tests are located in the `fwprovider/tests` directory.

#### Prerequisites

1. A running Proxmox instance (see [Setup Proxmox for Tests](docs/guides/setup-proxmox-for-tests.md))
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

#### Running Acceptance Tests

Run the acceptance test suite with:

```sh
make testacc
```

If you want to run a single test or a group of tests, use the helper script:

```sh
./testacc <test_name>
```

For example, to run all VM-related tests: `./testacc.sh TestAccResourceVM.*`

> [!NOTE]
>
> - Acceptance test coverage is still in development
> - Only some resources and data sources are currently tested
> - Some tests may require specific Proxmox configuration

### Manual Testing

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

## Coding conventions

We expect all code contributions to follow these guidelines:

1. Code must be formatted using `gofmt`
   - Run `make fmt` to format your code

2. Code must be linted using `golangci-lint`
   - Run `make lint` to lint your code

## Commit message conventions

We follow the [Conventional Commits](https://www.conventionalcommits.org/) specification. Please use the following types for your commits:

- `feat`: New features
- `fix`: Bug fixes
- `chore`: Maintenance tasks

These types are used to automatically generate the changelog. Other types will be ignored.

Use the `scope` field to indicate the area of the codebase being changed:

- `vm` - Virtual Machine resources
- `lxc` - Container resources
- `provider` - Provider configuration and resources
- `core` - Core libraries and utilities
- `docs` - Documentation
- `ci` - Continuous Integration / Actions / GitHub Workflows

Guidelines:

- Use lowercase for descriptions
- Do not end descriptions with a period
- Keep the first line under 72 characters

Example:

```commit
feat(vm): add support for the `clone` operation
```

### Developer Certificate of Origin (DCO)

All contributions must be signed off according to the Developer Certificate of Origin (DCO). The DCO is a lightweight way of certifying that you wrote or have the right to submit the code you are contributing. You can find the full text [here](https://developercertificate.org).

To sign off your commits, add a `Signed-off-by` line to your commit message:

```commit
feat(vm): add support for the `clone` operation

Signed-off-by: Random Developer <random@developer.example.org>
```

> [!NOTE]
>
> - Use your real name and a valid email address
> - You can use GitHub's `noreply` email address for privacy (see [GitHub docs](https://docs.github.com/en/account-and-profile/setting-up-and-managing-your-personal-account-on-github/managing-email-preferences/setting-your-commit-email-address#setting-your-commit-email-address-on-github))
> - If your Git config has `user.name` and `user.email` set, use `git commit -s` to automatically add the sign-off

For more details about the DCO checker, see the [DCO app repo](https://github.com/dcoapp/app).

## Submitting changes

1. Create a new PR against the `main` branch using the project's [pull request template](.github/PULL_REQUEST_TEMPLATE.md)
2. Ensure your PR title follows the Conventional Commits specification (we use this as the squash commit message)
3. All commits in a PR are typically squashed on merge

## Releasing

We use [release-please](https://github.com/googleapis/release-please) GitHub Action for automated release management. The process works as follows:

1. The action creates a release PR based on commit messages
2. The PR includes an auto-generated changelog and version bump
3. Maintainers review and merge the release PR
4. The release is automatically published to:
   - GitHub Releases
   - Terraform Registry

We aim to release new versions every 1-2 weeks.
