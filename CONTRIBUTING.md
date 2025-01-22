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

- Clone the repository to `$GOPATH/src/github.com/bpg/terraform-provider-proxmox`:

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

- You also can cross-compile the provider for all supported platforms:

  ```sh
  make build-all
  ```

  The binaries will be placed in the `dist` directory.

## IDE support

If you are using VS Code, feel free to copy `settings.json` from `.vscode/settings.example.json`.

## Testing

The project has a handful of test cases which must pass for a contribution to be
accepted. We also expect that you either create new test cases or modify
existing ones in order to target your changes.

You can run all the test cases by invoking `make test`.

### Acceptance tests

The project has a limited set of acceptance tests which are run against a real Proxmox
instance. These tests are developed alongside the framework-based resource and datasource implementations, and are located in the `fwprovider/tests` directory.

To run the acceptance tests, you need to have a Proxmox instance available. See more details in the [Setup Proxmox for Tests](docs/guides/setup-proxmox-for-tests.md) section.
Create a `testacc.env` file in the project's root directory with the following contents:

```env
TF_ACC=1
PROXMOX_VE_API_TOKEN="root@pam!<token name>=<token value>"
PROXMOX_VE_ENDPOINT="https://<pve instance>:8006/"
PROXMOX_VE_SSH_AGENT="true"
PROXMOX_VE_SSH_USERNAME="root"
# optionally, you can override the default node name and ssh address and other things
#PROXMOX_VE_ACC_NODE_NAME="pve1"
#PROXMOX_VE_ACC_NODE_SSH_ADDRESS="10.0.0.11"
#PROXMOX_VE_ACC_NODE_SSH_PORT="22"
#PROXMOX_VE_ACC_IFACE_NAME="enp1s0"
```

Then use `make testacc` to run all the acceptance tests.
You can also target specific acceptance tests with `make testacc-general`, `make testacc-vm`, `make testacc-ct`.

> [!NOTE]
> The acceptance tests support is still in development. Only handful of resources and data sources are covered by the tests. Some tests may require extra configuration on the Proxmox instance, and fail if the configuration is not present.

## Manual Testing

You can manually test the provider by running it locally. This is useful for
testing changes to the provider before submitting a PR.

- Create a $HOME/.terraformrc (POSIX) or %APPDATA%/terraform.rc (Windows) file with the following contents:

  ```hcl
  provider_installation {

    dev_overrides {
        "bpg/proxmox" = "/home/user/go/bin/" # <- put an absolute path where $GOPATH/bin is pointing to in your system.
    }

    # For all other providers, install them directly from their origin provider
    # registries as normal. If you omit this, Terraform will _only_ use
    # the dev_overrides block, and so no other providers will be available.
    direct {}
  }
  ```

- Build & install the provider by running the following command in the provider directory:

  ```bash
  go install .
  ```

- Run `terraform init` in a directory containing a Terraform configuration
  using the provider. You should see output similar to the following:

  ```bash
  ❯ terraform init -upgrade

  Initializing the backend...

  Initializing provider plugins...

  ...

  ╷
  │ Warning: Provider development overrides are in effect
  │
  │ The following provider development overrides are set in the CLI configuration:
  │  - bpg/proxmox in /home/user/go/bin
  │
  │ Skip terraform init when using provider development overrides. It is not necessary and may error unexpectedly.
  ╵

  Terraform has been successfully initialized!
  ```

- Run `terraform plan` or `terraform apply` to test your changes.

> [!TIP]
> You don't need to run `terraform init` again after making changes to the provider, as long as you have the `dev_overrides` block in your `terraform.rc` file, and the provider is installed in the path specified in the `dev_overrides` block by running `go install .` in the provider directory.

## Coding conventions

We expect that all code contributions have been formatted using `gofmt`.

You can run `make fmt` to format your code.

We also expect that all code contributions have been linted using `golangci-lint`.

You can run `make lint` to lint your code.

## Commit message conventions

We expect that all commit messages follow the
[Conventional Commits](https://www.conventionalcommits.org/) specification.
Please use the `feat`, `fix` or `chore` types for your commits, as they will
be used to automatically generate the changelog. Other types will be ignored
in the changelog.

Please use the `scope` field to indicate the area of the codebase that is being
changed. For example, `vm` for changes in the Virtual Machine resource, or
`lxc` for changes in the Container resource.

Common scopes are:

- `vm` - Virtual Machine resources
- `lxc` - Container resources
- `provider` - Provider configuration and resources
- `core` - Core libraries and utilities
- `docs` - Documentation
- `ci` - Continuous Integration / Actions / GitHub Workflows

Please use lowercase for the description and do not end it with a period.

For example:

```commit
feat(vm): add support for the `clone` operation
```

### Developer Certificate of Origin

In order for a code change to be accepted, you'll also have to accept the
Developer Certificate of Origin (DCO).
It's very lightweight, and you can find it [here](https://developercertificate.org).
Accepting is accomplished by signing off on your commits, you can do this by
adding a `Signed-off-by` line to your commit message, like here:

```commit
feat(vm): add support for the `clone` operation

Signed-off-by: Random Developer <random@developer.example.org>
```

Please use your real name and a valid email address. If you'd like to keep your
email address private, you can use a GitHub-provided `noreply`` email address.
For more information, see "[Setting your commit email address.](https://docs.github.com/en/account-and-profile/setting-up-and-managing-your-personal-account-on-github/managing-email-preferences/setting-your-commit-email-address#setting-your-commit-email-address-on-github)"

If you set your `user.name` and `user.email` in Git, you can sign your commit
automatically with the `-s` flag:

```shell
> git commit -s -m 'feat(vm): add a cool new feature'
```

You can find more details about the DCO checker in the [DCO app repo](https://github.com/dcoapp/app).

## Submitting changes

Please create a new PR against the `main` branch which must be based on the
project's [pull request template](.github/PULL_REQUEST_TEMPLATE.md).

We usually squash all PRs commits on merge, and use the PR title as the commit
message. Therefore, the PR title should follow the
[Conventional Commits](https://www.conventionalcommits.org/) specification as well.

## Releasing

We use automated release management orchestrated
by [release-please](https://github.com/googleapis/release-please) GitHub Action. The action
creates a new release PR with the changelog and bumps the version based on the
commit messages. The release PR is merged by the maintainers.

The release will be published to the GitHub Releases page and the Terraform Registry.

We aim to release a new version every 1-2 weeks.
