# Terraform Provider for Proxmox

[![Go Report Card](https://goreportcard.com/badge/github.com/bpg/terraform-provider-proxmox)](https://goreportcard.com/report/github.com/bpg/terraform-provider-proxmox)
[![GoDoc](https://godoc.org/github.com/bpg/terraform-provider-proxmox?status.svg)](http://godoc.org/github.com/bpg/terraform-provider-proxmox)
[![GitHub release (latest by date)](https://img.shields.io/github/v/release/bpg/terraform-provider-proxmox)](https://github.com/bpg/terraform-provider-proxmox/releases/latest)
[![GitHub Release Date](https://img.shields.io/github/release-date/bpg/terraform-provider-proxmox)](https://github.com/bpg/terraform-provider-proxmox/releases/latest)
[![GitHub stars](https://img.shields.io/github/stars/bpg/terraform-provider-proxmox)](https://github.com/bpg/terraform-provider-proxmox/stargazers)
[![All Contributors](https://img.shields.io/github/all-contributors/bpg/terraform-provider-proxmox)](#contributors)
[![Conventional Commits](https://img.shields.io/badge/conventional%20commits-v1.0.0-ff69b4)](https://www.conventionalcommits.org/en/v1.0.0/)
[![Buy Me A Coffee](https://img.shields.io/badge/-buy%20me%20a%20coffee-5F7FFF?logo=buymeacoffee&labelColor=gray&logoColor=FFDD00)](https://www.buymeacoffee.com/bpgca)

A Terraform Provider which adds support for Proxmox solutions.

This repository is a fork
of <https://github.com/danitso/terraform-provider-proxmox>
which is no longer maintained.

## Compatibility promise

This provider is compatible with the latest version of Proxmox VE (currently
8.0). While it may work with older 7.x versions, it is not guaranteed to do so.

While provider is on version 0.x, it is not guaranteed to be backwards
compatible with all previous minor versions. However, we will try to keep the
backwards compatibility between provider versions as much as possible.

## Requirements

- [Terraform](https://www.terraform.io/downloads.html) 1.2+
- [Go](https://golang.org/doc/install) 1.20+ (to build the provider plugin)

## Building the provider

- Clone the repository
  to `$GOPATH/src/github.com/bpg/terraform-provider-proxmox`:

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

## Using the provider

You can find the latest release and its documentation in
the [Terraform Registry](https://registry.terraform.io/providers/bpg/proxmox/latest).

## Testing the provider

In order to test the provider, you can simply run `make test`.

```sh
make test
```

Tests are limited to regression tests, ensuring backwards compatibility.

## Deploying the example resources

There are number of TF examples in the `examples` directory, which can be used
to deploy a Container, VM, or other Proxmox resources on your test Proxmox
cluster. The following assumptions are made about the test Proxmox cluster:

- It has one node named `pve`
- The node has local storages named `local` and `local-lvm`

Create `examples/terraform.tfvars` with the following variables:

```sh
virtual_environment_username = "root@pam"
virtual_environment_password = "put-your-password-here"
virtual_environment_endpoint = "https://<your-cluster-endpoint>:8006/"
```

Then run `make example` to deploy the example resources.

## Future work

The provider is using
the [Terraform SDKv2](https://developer.hashicorp.com/terraform/plugin/sdkv2),
which is considered legacy and is in maintenance mode.
The work has started to migrate the provider to the
new [Terraform Plugin Framework](https://www.terraform.io/docs/extend/plugin-sdk.html),
with aim to release it as a new major version **1.0**.

## Known issues

### Disk images cannot be imported by non-PAM accounts

Due to limitations in the Proxmox VE API, certain actions need to be performed
using SSH. This requires the use of a PAM account (standard Linux account).

### Disk images from VMware cannot be uploaded or imported

Proxmox VE is not currently supporting VMware disk images directly. However, you
can still use them as disk images by using this workaround:

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

Due to limitations in the Proxmox VE API, certain files need to be uploaded
using SFTP. This requires the use of a PAM account (standard Linux account).

## Contributors

<!-- ALL-CONTRIBUTORS-LIST:START - Do not remove or modify this section -->
<!-- prettier-ignore-start -->
<!-- markdownlint-disable -->
<table>
  <tbody>
    <tr>
      <td align="center" valign="top" width="14.28%"><a href="https://github.com/luhahn"><img src="https://avatars.githubusercontent.com/u/61747797?v=4?s=100" width="100px;" alt="Lucas Hahn"/><br /><sub><b>Lucas Hahn</b></sub></a><br /><a href="https://github.com/bpg/terraform-provider-proxmox/commits?author=luhahn" title="Code">💻</a></td>
    </tr>
  </tbody>
</table>

<!-- markdownlint-restore -->
<!-- prettier-ignore-end -->

<!-- ALL-CONTRIBUTORS-LIST:END -->

## Stargazers

[![Star History Chart](https://api.star-history.com/svg?repos=bpg/terraform-provider-proxmox&type=Date)](https://star-history.com/#bpg/terraform-provider-proxmox&Date)
