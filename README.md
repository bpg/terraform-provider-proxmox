# Terraform Provider for Proxmox

[![Go Report Card](https://goreportcard.com/badge/github.com/bpg/terraform-provider-proxmox)](https://goreportcard.com/report/github.com/bpg/terraform-provider-proxmox)
[![GoDoc](https://godoc.org/github.com/bpg/terraform-provider-proxmox?status.svg)](http://godoc.org/github.com/bpg/terraform-provider-proxmox)
[![GitHub release (latest by date)](https://img.shields.io/github/v/release/bpg/terraform-provider-proxmox)](https://github.com/bpg/terraform-provider-proxmox/releases/latest)
[![GitHub Release Date](https://img.shields.io/github/release-date/bpg/terraform-provider-proxmox)](https://github.com/bpg/terraform-provider-proxmox/releases/latest)
[![GitHub stars](https://img.shields.io/github/stars/bpg/terraform-provider-proxmox?style=flat)](https://github.com/bpg/terraform-provider-proxmox/stargazers)
[![All Contributors](https://img.shields.io/github/all-contributors/bpg/terraform-provider-proxmox)](#contributors)
[![Conventional Commits](https://img.shields.io/badge/conventional%20commits-v1.0.0-ff69b4)](https://www.conventionalcommits.org/en/v1.0.0/)
[![Buy Me A Coffee](https://img.shields.io/badge/-buy%20me%20a%20coffee-5F7FFF?logo=buymeacoffee&labelColor=gray&logoColor=FFDD00)](https://www.buymeacoffee.com/bpgca)

A Terraform / OpenTofu Provider which adds support for Proxmox solutions.

This repository is a fork of <https://github.com/danitso/terraform-provider-proxmox> which is no longer maintained.

## Compatibility promise

This provider is compatible with the latest version of Proxmox VE (currently 8.0).
While it may work with older 7.x versions, it is not guaranteed to do so.

While provider is on version 0.x, it is not guaranteed to be backwards compatible with all previous minor versions.
However, we will try to keep the backwards compatibility between provider versions as much as possible.

## Requirements

- [Proxmox Virtual Environment](https://www.proxmox.com/en/proxmox-virtual-environment/) 8.x
- TLS 1.3 for the Proxmox API endpoint (legacy TLS 1.2 is optionally supported)
- [Terraform](https://www.terraform.io/downloads.html) 1.5.x+ or [OpenTofu](https://opentofu.org) 1.6.x
- [Go](https://golang.org/doc/install) 1.21 (to build the provider plugin)

## Using the provider

You can find the latest release and its documentation in the [Terraform Registry](https://registry.terraform.io/providers/bpg/proxmox/latest).

## Testing the provider

In order to test the provider, you can simply run `make test`.

```sh
make test
```

Tests are limited to regression tests, ensuring backwards compatibility.

A limited number of acceptance tests are available in the `proxmoxtf/test` directory, mostly for "new" functionality implemented using the Terraform Provider Framework.
These tests are not run by default, as they require a Proxmox VE environment to be available.
They can be run using `make testacc`, the Proxmox connection can be configured using environment variables, see provider documentation for details.

## Deploying the example resources

There are number of TF examples in the `example` directory, which can be used to deploy a Container, VM, or other Proxmox resources on your test Proxmox environment.
The following assumptions are made about the test environment:

- It has one node named `pve`
- The node has local storages named `local` and `local-lvm`
- The "Snippets" content type is enabled in `local` storage

Create `example/terraform.tfvars` with the following variables:

```sh
virtual_environment_username = "root@pam"
virtual_environment_password = "put-your-password-here"
virtual_environment_endpoint = "https://<your-cluster-endpoint>:8006/"
```

Then run `make example` to deploy the example resources.

If you don't have free proxmox cluster to play with, there is dedicated [how-to tutorial](docs/guides/setup-proxmox-for-tests.md) how to setup Proxmox inside VM and run `make example` on it.

## Future work

The provider is using the [Terraform SDKv2](https://developer.hashicorp.com/terraform/plugin/sdkv2), which is considered legacy and is in maintenance mode.
The work has started to migrate the provider to the new [Terraform Plugin Framework](https://www.terraform.io/docs/extend/plugin-sdk.html), with aim to release it as a new major version **1.0**.

## Known issues

### Disk images cannot be imported by non-PAM accounts

Due to limitations in the Proxmox VE API, certain actions need to be performed using SSH. This requires the use of a PAM account (standard Linux account).

### Disk images from VMware cannot be uploaded or imported

Proxmox VE is not currently supporting VMware disk images directly.
However, you can still use them as disk images by using this workaround:

```hcl
resource "proxmox_virtual_environment_file" "vmdk_disk_image" {
  content_type = "iso"
  datastore_id = "datastore-id"
  node_name    = "node-name"

  source_file {
    # We must override the file extension to bypass the validation code
    # in the Proxmox VE API.
    file_name = "vmdk-file-name.img"
    path      = "path-to-vmdk-file"
  }
}

resource "proxmox_virtual_environment_vm" "example" {
  //...

  disk {
    datastore_id = "datastore-id"
    # We must tell the provider that the file format is vmdk instead of qcow2.
    file_format  = "vmdk"
    file_id      = "${proxmox_virtual_environment_file.vmdk_disk_image.id}"
  }

  //...
}
```

### Snippets cannot be uploaded by non-PAM accounts

Due to limitations in the Proxmox VE API, certain files (snippets, backups) need to be uploaded using SFTP.
This requires the use of a PAM account (standard Linux account).

## Contributors

See [CONTRIBUTORS.md](CONTRIBUTORS.md) for a list of contributors to this project.

## Repository Metrics

<picture>
  <img src="https://gist.githubusercontent.com/bpg/2cc44ead81225542ed1ef0303d8f9eb9/raw/metrics.svg?p" alt="Metrics">
</picture>

## Sponsorship

❤️ This project is sponsored by:

- [TJ Zimmerman](https://github.com/zimmertr)
- [Elias Alvord](https://github.com/elias314)
- [laktosterror](https://github.com/laktosterror)

Thanks again for your support, it is much appreciated! 🙏
