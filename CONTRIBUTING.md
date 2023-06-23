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

## Testing

The project has a handful of test cases which must pass for a contribution to be
accepted. We also expect that you either create new test cases or modify
existing ones in order to target your changes.

You can run all the test cases by invoking `make test`.

## Coding conventions

We expect that all code contributions have been formatted using `gofmt`. You can
run `make fmt` to format your code.

We also expect that all code contributions have been linted
using `golangci-lint`.
You can run `make lint` to lint your code.

## Commit message conventions

We expect that all commit messages follow the
[Conventional Commits](https://www.conventionalcommits.org/) specification.

Please use the `scope` field to indicate the area of the codebase that is being
changed. For example, `vm` for changes in the Virtual Machine resource, or
`lcx` for changes in the Container resource.

Common scopes are:

- `vm` - Virtual Machine resources
- `lcx` - Container resources
- `provider` - Provider configuration and resources
- `core` - Core libraries and utilities
- `docs` - Documentation
- `ci` - Continuous Integration / Actions / GitHub Workflows

Please use lowercase for the description and do not end it with a period.

For example:

```
feat(vm): add support for the `clone` operation
```

## Submitting changes

Please create a new PR against the `main` branch which must be based on the
project's [pull request template](.github/PULL_REQUEST_TEMPLATE.md).

We usually squash all PRs commits on merge, and use the PR title as the commit
message. Therefore, the PR title should follow the
[Conventional Commits](https://www.conventionalcommits.org/) specification as
well.

## Releasing

We use automated release management orchestrated
by https://github.com/googleapis/release-please GitHub Action. The action
creates a new release PR with the changelog and bumps the version based on the
commit messages. The release PR is merged by the maintainers.

The release will be published to the GitHub Releases page and the Terraform
Registry.

We aim to release a new version every 1-2 weeks.
