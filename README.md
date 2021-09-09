[![Go Report Card](https://goreportcard.com/badge/github.com/danitso/terraform-provider-proxmox)](https://goreportcard.com/report/github.com/danitso/terraform-provider-proxmox)
[![GoDoc](https://godoc.org/github.com/danitso/terraform-provider-proxmox?status.svg)](http://godoc.org/github.com/danitso/terraform-provider-proxmox)

# Terraform Provider for Proxmox
A Terraform Provider which adds support for Proxmox solutions.

## Requirements
- [Terraform](https://www.terraform.io/downloads.html) 0.13+
- [Go](https://golang.org/doc/install) 1.16+ (to build the provider plugin)
- [GoReleaser](https://goreleaser.com/install/) 0.155+ (to build the provider plugin)

## Table of Contents
- [Building the provider](#building-the-provider)
- [Using the provider](#using-the-provider)
- [Testing the provider](#testing-the-provider)
- [Known issues](#known-issues)

## Building the provider
- Clone the repository to `$GOPATH/src/github.com/danitso/terraform-provider-proxmox`:

    ```sh
    $ mkdir -p "${GOPATH}/src/github.com/danitso"
    $ cd "${GOPATH}/src/github.com/danitso"
    $ git clone git@github.com:danitso/terraform-provider-proxmox
    ```

- Enter the provider directory and build it:

    ```sh
    $ cd "${GOPATH}/src/github.com/danitso/terraform-provider-proxmox"
    $ make build
    ```

## Using the provider
You can find the latest release and its documentation in the [Terraform Registry](https://registry.terraform.io/providers/danitso/proxmox/latest).

## Testing the provider
In order to test the provider, you can simply run `make test`.

```sh
$ make test
```

Tests are limited to regression tests, ensuring backwards compability.

## Known issues

### Disk images cannot be imported by non-PAM accounts
Due to limitations in the Proxmox VE API, certain actions need to be performed using SSH. This requires the use of a PAM account (standard Linux account).

### Disk images from VMware cannot be uploaded or imported
Proxmox VE is not currently supporting VMware disk images directly. However, you can still use them as disk images by using this workaround:

```hcl
resource "proxmox_virtual_environment_file" "vmdk_disk_image" {
  content_type = "iso"
  datastore_id = "datastore-id"
  node_name    = "node-name"

  source_file {
    # We must override the file extension to bypass the validation code in the Proxmox VE API.
    file_name = "vmdk-file-name.img"
    path      = "path-to-vmdk-file"
  }
}

resource "proxmox_virtual_environment_vm" "example" {
  ...

  disk {
    datastore_id = "datastore-id"
    # We must tell the provider that the file format is vmdk instead of qcow2.
    file_format  = "vmdk"
    file_id      = "${proxmox_virtual_environment_file.vmdk_disk_image.id}"
  }

  ...
}
```

### Snippets cannot be uploaded by non-PAM accounts
Due to limitations in the Proxmox VE API, certain files need to be uploaded using SFTP. This requires the use of a PAM account (standard Linux account).
