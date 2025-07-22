# Terraform / OpenTofu Provider for Proxmox VE

[![GitHub release (latest by date)](https://img.shields.io/github/v/release/bpg/terraform-provider-proxmox)](https://github.com/bpg/terraform-provider-proxmox/releases/latest)
[![GitHub Release Date](https://img.shields.io/github/release-date/bpg/terraform-provider-proxmox)](https://github.com/bpg/terraform-provider-proxmox/releases/latest)
[![GitHub stars](https://img.shields.io/github/stars/bpg/terraform-provider-proxmox?style=flat)](https://github.com/bpg/terraform-provider-proxmox/stargazers)
[![All Contributors](https://img.shields.io/github/all-contributors/bpg/terraform-provider-proxmox)](#contributors)
[![Go Report Card](https://goreportcard.com/badge/github.com/bpg/terraform-provider-proxmox)](https://goreportcard.com/report/github.com/bpg/terraform-provider-proxmox)
[![Conventional Commits](https://img.shields.io/badge/conventional%20commits-v1.0.0-ff69b4)](https://www.conventionalcommits.org/en/v1.0.0/)

A Terraform / OpenTofu Provider that adds support for Proxmox Virtual Environment.

This repository is a fork of <https://github.com/danitso/terraform-provider-proxmox> which is no longer maintained.

## Disclaimer

This project is a personal open-source initiative and is not affiliated with, endorsed by, or associated with any of my current or former employers. All opinions, code, and documentation are solely those of myself and the individual contributors.

The project is not affiliated with [Proxmox Server Solutions GmbH](https://www.proxmox.com/en/about/about-us/company) or any of its subsidiaries. The use of the Proxmox name and/or logo is for informational purposes only and does not imply any endorsement or affiliation with the Proxmox project.

## Compatibility Promise

This provider is compatible with Proxmox VE 8.x (currently **8.4**).

> [!IMPORTANT]
> Proxmox VE 7.x is NOT supported. While some features might work with 7.x, we do not test against it, and issues specific to 7.x will not be addressed.

While the provider is on version 0.x, it is not guaranteed to be backward compatible with all previous minor versions.
However, we will try to maintain backward compatibility between provider versions as much as possible.

## Requirements

### Production Requirements

- [Proxmox Virtual Environment](https://www.proxmox.com/en/proxmox-virtual-environment/) 8.x
- TLS 1.3 for the Proxmox API endpoint (legacy TLS 1.2 is optionally supported)
- [Terraform](https://www.terraform.io/downloads.html) 1.5+ or [OpenTofu](https://opentofu.org) 1.6+

### Development Requirements

- [Go](https://golang.org/doc/install) 1.24 (to build the provider plugin)
- [Docker](https://www.docker.com/products/docker-desktop/) (optional, for running dev tools)

## Using the Provider

You can find the latest release and its documentation in the [Terraform Registry](https://registry.terraform.io/providers/bpg/proxmox/latest) or [OpenTofu Registry](https://search.opentofu.org/provider/bpg/proxmox/latest).

For manual provider installation, you can download the binaries from the [Releases](https://github.com/bpg/terraform-provider-proxmox/releases) page.
You also can use `gh` tool to verify the binaries provenance, see more details [here](https://github.com/bpg/terraform-provider-proxmox/attestations/).

## Testing the Provider

To test the provider, simply run `make test`.

```sh
make test
```

Tests are limited to regression tests, ensuring backward compatibility.

A limited number of acceptance tests are available in the `proxmoxtf/test` directory, mostly for "new" functionality implemented using the Terraform Provider Framework.
These tests are not run by default, as they require a Proxmox VE environment to be available.
They can be run using `make testacc`. The Proxmox connection can be configured using environment variables; see the provider documentation for details.

## Deploying the Example Resources

There are a number of TF examples in the `example` directory, which can be used to deploy a Container, VM, or other Proxmox resources in your test Proxmox environment.

### Prerequisites

The following assumptions are made about the test environment:

- It has one node named `pve`
- The node has local storages named `local` and `local-lvm`
- The "Snippets" and "Import" content types are enabled in the `local` storage
- Default Linux Bridge "vmbr0" is VLAN aware (datacenter -> pve -> network -> edit & apply)

Create `example/terraform.tfvars` with the following variables:

```sh
virtual_environment_endpoint                 = "https://pve.example.doc:8006/"
virtual_environment_ssh_username             = "terraform"
virtual_environment_api_token                = "root@pam!terraform=00000000-0000-0000-0000-000000000000"
```

Then run `make example` to deploy the example resources.

If you don't have a free Proxmox cluster to play with, there is a dedicated [how-to tutorial](docs/guides/setup-proxmox-for-tests.md) on how to set up Proxmox inside a VM and run `make example` on it.

## Future Work

The provider is using the [Terraform SDKv2](https://developer.hashicorp.com/terraform/plugin/sdkv2), which is considered legacy and is in maintenance mode.
Work has started to migrate the provider to the new [Terraform Plugin Framework](https://developer.hashicorp.com/terraform/plugin/framework), with the aim of releasing it as a new major version **1.0**.

## Known Issues

### HA VMs / containers

If a VM or container resource is created with the provider but managed by an HA cluster, it might be migrated to a different node without the provider being aware of the change.
This causes a "configuration drift" and the provider will report an error when managing the resource.
You would need to manually reconcile the resource state stored in the backend to match the actual state of the resource, or remove the resource from the provider management.

### Serial Device Required for Debian 12 / Ubuntu VMs

Debian 12 and Ubuntu VMs throw kernel panic when resizing a cloud image boot disk, as they require a serial device configured.
Add the following block to your VM config:

```hcl
  serial_device {
    device = "socket"
  }
```

For more context, see #1639 and #1770.

### Disk Images from VMware Cannot Be Uploaded or Imported

Proxmox VE does not currently support VMware disk images directly.
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

### Snippets Cannot Be Uploaded by Non-PAM Accounts

Due to limitations in the Proxmox VE API, certain files (snippets, backups) need to be uploaded using SFTP.
This requires the use of a PAM account (standard Linux account).

### Cluster Hardware Mappings Cannot Be Created by Non-PAM Accounts

Due to limitations in the Proxmox VE API, cluster hardware mappings must be created using the `root` PAM account (standard Linux account) due to [IOMMU](https://en.wikipedia.org/wiki/Input%E2%80%93output_memory_management_unit#Virtualization) interactions.
Hardware mappings allow the use of [PCI "passthrough"](https://pve.proxmox.com/wiki/PCI_Passthrough) and [map physical USB ports](https://pve.proxmox.com/wiki/USB_Physical_Port_Mapping).

## Contributors

See [CONTRIBUTORS.md](CONTRIBUTORS.md) for a list of contributors to this project.

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
- [Radosław Szamszur](https://github.com/rszamszur)
- [Marshall Ford](https://github.com/marshallford)

Thanks again for your continuous support, it is much appreciated! 🙏

## Acknowledgements

This project has been developed with **GoLand** IDE under the [JetBrains Open Source license](https://www.jetbrains.com/community/opensource/#support), generously provided by JetBrains s.r.o.

<img src="https://resources.jetbrains.com/storage/products/company/brand/logos/GoLand_icon.png" alt="GoLand logo" width="80">
