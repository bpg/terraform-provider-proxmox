# Terraform / OpenTofu Provider for Proxmox VE

[![GitHub release (latest by date)](https://img.shields.io/github/v/release/bpg/terraform-provider-proxmox)](https://github.com/bpg/terraform-provider-proxmox/releases/latest)
[![GitHub Release Date](https://img.shields.io/github/release-date/bpg/terraform-provider-proxmox?display_date=published_at)](https://github.com/bpg/terraform-provider-proxmox/releases/latest)
[![Terraform Provider Downloads](https://img.shields.io/terraform/provider/dw/2185)](https://registry.terraform.io/providers/bpg/proxmox/latest)
[![All Contributors](https://img.shields.io/github/all-contributors/bpg/terraform-provider-proxmox)](#contributors)
[![GitHub stars](https://img.shields.io/github/stars/bpg/terraform-provider-proxmox?style=flat)](https://github.com/bpg/terraform-provider-proxmox)
![GitHub Sponsors](https://img.shields.io/github/sponsors/bpg-dev)

A Terraform / OpenTofu Provider that adds support for Proxmox Virtual Environment.

Maintained by [BPG Labs](https://bpg.sh).

## Disclaimer

This project is a personal open-source initiative by [@bpg-dev](https://github.com/bpg-dev) and is not affiliated with, endorsed by, or associated with any of their current or former employers. All opinions, code, and documentation are solely those of the maintainer and the individual contributors.

The project is not affiliated with [Proxmox Server Solutions GmbH](https://www.proxmox.com/en/about/about-us/company) or any of its subsidiaries. The use of the Proxmox name and/or logo is for informational purposes only and does not imply any endorsement or affiliation with the Proxmox project.

## Quick Start

```hcl
terraform {
  required_providers {
    proxmox = {
      source = "bpg/proxmox"
    }
  }
}

provider "proxmox" {
  endpoint  = "https://pve.example.com:8006/"
  api_token = "terraform@pve!provider=xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx"
}
```

Most resources work entirely through the Proxmox VE API. SSH access to the nodes is only needed for a few operations, such as uploading snippets.
See the provider documentation at [bpg.sh/docs](https://bpg.sh/docs) for authentication options, SSH configuration, the full list of resources and data sources, and the [known issues](https://bpg.sh/docs/guides/known-issues).

Releases are published to the [Terraform Registry](https://registry.terraform.io/providers/bpg/proxmox/latest) and the [OpenTofu Registry](https://search.opentofu.org/provider/bpg/proxmox/latest).
Binaries for manual installation are on the [Releases](https://github.com/bpg/terraform-provider-proxmox/releases) page, and their provenance can be verified with `gh`, see [attestations](https://github.com/bpg/terraform-provider-proxmox/attestations/).

## Compatibility

This provider targets Proxmox VE 9.x and is acceptance-tested against the current 9.x release.

> [!IMPORTANT]
> Proxmox VE 8.x is supported, but some functionality might be limited or not work as expected. Testing against 8.x is not a priority, and issues specific to 8.x will not be addressed.
>
> Proxmox VE 7.x is NOT supported. While some features might work with 7.x, we do not test against it, and issues specific to 7.x will not be addressed.

While the provider is on version 0.x, it is not guaranteed to be backward compatible with all previous minor versions.
However, we will try to maintain backward compatibility between provider versions as much as possible.

Terraform and OpenTofu version requirements are listed in the [provider documentation](https://bpg.sh/docs#requirements).

## Security

- Report vulnerabilities privately as described in [SECURITY.md](.github/SECURITY.md). Do not open public issues for them.
- Published advisories are listed on the [security advisories](https://github.com/bpg/terraform-provider-proxmox/security/advisories) page.
- Nearly all operations use the Proxmox VE API. Use an API token with only the privileges your resources need.
- The few operations that require SSH run commands on the node through `sudo`. Use the sudoers configuration from the [SSH User](https://bpg.sh/docs#ssh-user) section of the documentation as written. Granting the SSH user unrestricted `sudo` access to binaries such as `qm` or `pvesm` is equivalent to giving it root on the host.

## Building and Testing the Provider

Building requires [Go](https://golang.org/doc/install) 1.26. [Docker](https://www.docker.com/products/docker-desktop/) is optional, for running the dev tools.

```sh
make build
```

Unit tests cover the API client, model conversion and schema logic, and run with:

```sh
make test
```

Most of the test coverage is acceptance tests, which exercise the provider end-to-end against a real Proxmox VE instance.
They live next to the resource or data source they cover (and in `fwprovider/test` for cross-resource scenarios), are grouped into light/medium/heavy tiers, and are run with `./testacc` (requires `testacc.env` in the project root, see [CONTRIBUTING.md](CONTRIBUTING.md#acceptance-tests)).
They are not run in CI by default, as they require a Proxmox VE environment; maintainers run them against a lab cluster before merging changes, and new or changed functionality must come with acceptance tests.

## Example Resources

The `example` directory contains sample configurations, mostly for the legacy SDKv2 resources; Framework resources are covered by their acceptance tests instead.
`make example` builds the provider, applies the samples against a test Proxmox VE environment, and destroys them again, so it needs a dedicated environment and an `example/terraform.tfvars`.
See [Setting up Proxmox in a VM for development](docs/guides/dev-proxmox-setup.md) for the prerequisites and configuration.

## Future Work

The provider was originally built on the [Terraform SDKv2](https://developer.hashicorp.com/terraform/plugin/sdkv2), which is considered legacy and is in maintenance mode.
New resources and data sources are implemented with the [Terraform Plugin Framework](https://developer.hashicorp.com/terraform/plugin/framework), and both are served from a single binary.
The core resources, including `proxmox_virtual_environment_vm`, `proxmox_virtual_environment_container`, `proxmox_virtual_environment_file` and the firewall and access-control resources, are still on SDKv2 and continue to receive fixes and enhancements there.
Migrating them to the Plugin Framework, starting with the VM resource ([#1231](https://github.com/bpg/terraform-provider-proxmox/issues/1231)), is a prerequisite for the **1.0** release.

## Contributors

See [CONTRIBUTORS.md](CONTRIBUTORS.md) for a list of contributors to this project.

## History

The provider was forked in 2021 from [danitso/terraform-provider-proxmox](https://github.com/danitso/terraform-provider-proxmox) after upstream maintenance stopped, and has been developed independently since.

## Repository Metrics

![Alt](https://repobeats.axiom.co/api/embed/bd0eca87c8a61f50b5fb6ff49a0d6c34de918963.svg "Repobeats analytics image")

## Sponsorship

❤️ This project is sponsored by:

- [Elias Alvord](https://github.com/elias314)
- [laktosterror](https://github.com/laktosterror)
- [Greg Brant](https://github.com/gregbrant2)
- [Serge](https://github.com/sergelogvinov)
- [Daniel Brennand](https://github.com/dbrennand)
- [Brian King](https://github.com/inflatador)
- [Marshall Ford](https://github.com/marshallford)
- [Simon Caron](https://github.com/simoncaron)

Thanks again for your continuous support, it is much appreciated! 🙏

## Acknowledgements

This project has been developed with **GoLand** IDE under the [JetBrains Open Source license](https://www.jetbrains.com/community/opensource/#support), generously provided by JetBrains s.r.o.

<img src="https://resources.jetbrains.com/storage/products/company/brand/logos/GoLand_icon.png" alt="GoLand logo" width="80">
