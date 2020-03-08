[![Build Status](https://api.travis-ci.com/danitso/terraform-provider-proxmox.svg?branch=master)](https://travis-ci.com/danitso/terraform-provider-proxmox)

# Terraform Provider for Proxmox
A Terraform Provider which adds support for Proxmox solutions.

## Requirements
- [Terraform](https://www.terraform.io/downloads.html) 0.11+
- [Go](https://golang.org/doc/install) 1.13 (to build the provider plugin)

## Table of contents
- [Building the Provider](#building-the-provider)
- [Using the Provider](#using-the-provider)
- [Developing the Provider](#developing-the-provider)
- [Testing the Provider](#testing-the-provider)
- [Known issues](#known-issues)

## Building the Provider
Clone repository to: `$GOPATH/src/github.com/danitso/terraform-provider-proxmox`

```sh
$ mkdir -p $GOPATH/src/github.com/danitso; cd $GOPATH/src/github.com/danitso
$ git clone git@github.com:danitso/terraform-provider-proxmox
```

Enter the provider directory, initialize and build the provider

```sh
$ cd $GOPATH/src/github.com/danitso/terraform-provider-proxmox
$ make init
$ make build
```

## Using the Provider
If you're building the provider, follow the instructions to [install it as a plugin.](https://www.terraform.io/docs/plugins/basics.html#installing-plugins) After placing it into your plugins directory, run `terraform init` to initialize it.

For more information, please refer to the [official documentation](http://danitso.com/terraform-provider-proxmox/).

## Developing the Provider
If you wish to work on the provider, you'll first need [Go](http://www.golang.org) installed on your machine (version 1.13+ is *required*). You'll also need to correctly setup a [GOPATH](http://golang.org/doc/code.html#GOPATH), as well as adding `$GOPATH/bin` to your `$PATH`.

To compile the provider, run `make build`. This will build the provider and put the provider binary in the `$GOPATH/bin` directory.

```sh
$ make build
...
$ $GOPATH/bin/terraform-provider-proxmox
...
```

If you wish to contribute to the provider, please see [CONTRIBUTING.md](CONTRIBUTING.md).

## Testing the Provider
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
