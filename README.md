[![Build Status](https://api.travis-ci.com/danitso/terraform-provider-proxmox.svg?branch=master)](https://travis-ci.com/danitso/terraform-provider-proxmox)

# Terraform Provider for Proxmox
A Terraform Provider which adds support for Proxmox solutions.

## Requirements
- [Terraform](https://www.terraform.io/downloads.html) 0.11+
- [Go](https://golang.org/doc/install) 1.13 (to build the provider plugin)

## Table of contents
- [Building the Provider](#building-the-provider)
- [Using the Provider](#using-the-provider)
    - [Configuration](#configuration)
        - [Arguments](#arguments)
        - [Environment variables](#environment-variables)
    - [Data Sources](#data-sources)
        - [Virtual Environment](#virtual-environment)
            - [Datastores](#datastores-proxmox_virtual_environment_datastores)
            - [Group](#group-proxmox_virtual_environment_group)
            - [Groups](#groups-proxmox_virtual_environment_groups)
            - [Nodes](#nodes-proxmox_virtual_environment_nodes)
            - [Pool](#pool-proxmox_virtual_environment_pool)
            - [Pools](#pools-proxmox_virtual_environment_pools)
            - [Role](#role-proxmox_virtual_environment_role)
            - [Roles](#roles-proxmox_virtual_environment_roles)
            - [User](#user-proxmox_virtual_environment_user)
            - [Users](#users-proxmox_virtual_environment_users)
            - [Version](#version-proxmox_virtual_environment_version)
    - [Resources](#resources)
        - [Virtual Environment](#virtual-environment-1)
            - [File](#file-proxmox_virtual_environment_file)
            - [Group](#group-proxmox_virtual_environment_group-1)
            - [Pool](#pool-proxmox_virtual_environment_pool-1)
            - [Role](#role-proxmox_virtual_environment_role-1)
            - [User](#user-proxmox_virtual_environment_user-1)
            - [VM](#vm-proxmox_virtual_environment_vm)
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

### Configuration

#### Arguments
* `virtual_environment` - (Optional) This is the configuration block for the Proxmox Virtual Environment
    * `endpoint` - (Required) The endpoint for the Proxmox Virtual Environment API
    * `insecure` - (Optional) Whether to skip the TLS verification step (defaults to `false`)
    * `password` - (Required) The password for the Proxmox Virtual Environment API
    * `username` - (Required) The username for the Proxmox Virtual Environment API

#### Environment variables
You can set up the provider by passing environment variables instead of specifying arguments.

* `PROXMOX_VE_ENDPOINT` or `PM_VE_ENDPOINT` - The endpoint for the Proxmox Virtual Environment API
* `PROXMOX_VE_INSECURE` or `PM_VE_INSECURE` - Whether to skip the TLS verification step
* `PROXMOX_VE_PASSWORD` or `PM_VE_PASSWORD` - The password for the Proxmox Virtual Environment API
* `PROXMOX_VE_USERNAME` or `PM_VE_USERNAME` - The username for the Proxmox Virtual Environment API

```hcl
provider "proxmox" {
  virtual_environment {}
}
```

##### Usage

```sh
$ export PROXMOX_VE_ENDPOINT="https://hostname:8006"
$ export PROXMOX_VE_INSECURE="true"
$ export PROXMOX_VE_PASSWORD="a-strong-password"
$ export PROXMOX_VE_USERNAME="username@realm"
$ terraform plan
```

You can omit `PROXMOX_VE_INSECURE`, if the Proxmox Virtual Environment API is exposing a certificate trusted by your operating system.

### Data Sources

#### Virtual Environment

##### Datastores (proxmox_virtual_environment_datastores)

###### Arguments
* `node_name` - (Required) A node name

###### Attributes
* `active` - Whether the datastore is active
* `content_types` - The allowed content types
* `datastore_ids` - The datastore ids
* `enabled` - Whether the datastore is enabled
* `shared` - Whether the datastore is shared
* `space_available` - The available space in bytes
* `space_total` - The total space in bytes
* `space_used` - The used space in bytes
* `types` - The storage types

##### Group (proxmox_virtual_environment_group)

###### Arguments
* `group_id` - (Required) The group id

###### Attributes
* `acl` - The access control list
    * `path` - The path
    * `propagate` - Whether to propagate to child paths
    * `role_id` - The role id
* `comment` - The group comment
* `members` - The group members as a list with `username@realm` entries

##### Groups (proxmox_virtual_environment_groups)

###### Arguments
This data source doesn't accept arguments.

###### Attributes
* `comments` - The group comments
* `group_ids` - The group ids

##### Nodes (proxmox_virtual_environment_nodes)

###### Arguments
This data source doesn't accept arguments.

###### Attributes
* `cpu_count` - The CPU count for each node
* `cpu_utilization` - The CPU utilization on each node
* `memory_available` - The memory available on each node
* `memory_used` - The memory used on each node
* `names` - The node names
* `online` - Whether a node is online
* `ssl_fingerprints` - The SSL fingerprint for each node
* `support_levels` - The support level for each node
* `uptime` - The uptime in seconds for each node

##### Pool (proxmox_virtual_environment_pool)

###### Arguments
* `pool_id` - (Required) The pool id

###### Attributes
* `comment` - The pool comment
* `members` - The pool members
    * `datastore_id` - The datastore id
    * `id` - The member id
    * `node_name` - The node name
    * `type` - The member type
    * `vm_id` - The virtual machine id

##### Pools (proxmox_virtual_environment_pools)

###### Arguments
This data source doesn't accept arguments.

###### Attributes
* `pool_ids` - The pool ids

##### Role (proxmox_virtual_environment_role)

###### Arguments
* `role_id` - (Required) The role id

###### Attributes
* `privileges` - The role privileges

##### Roles (proxmox_virtual_environment_roles)

###### Arguments
This data source doesn't accept arguments.

###### Attributes
* `privileges` - The role privileges
* `role_ids` - The role ids
* `special` - Whether the role is special (built-in)

##### User (proxmox_virtual_environment_user)

###### Arguments
* `user_id` - (Required) The user id.

###### Attributes
* `acl` - The access control list
    * `path` - The path
    * `propagate` - Whether to propagate to child paths
    * `role_id` - The role id
* `comment` - The user comment
* `email` - The user's email address
* `enabled` - Whether the user account is enabled
* `expiration_date` - The user account's expiration date (RFC 3339)
* `first_name` - The user's first name
* `groups` - The user's groups
* `keys` - The user's keys
* `last_name` - The user's last name

##### Users (proxmox_virtual_environment_users)

###### Arguments
This data source doesn't accept arguments.

###### Attributes
* `comments` - The user comments
* `emails` - The users' email addresses
* `enabled` - Whether a user account is enabled
* `expiration_dates` - The user accounts' expiration dates (RFC 3339)
* `first_names` - The users' first names
* `groups` - The users' groups
* `keys` - The users' keys
* `last_names` - The users' last names
* `user_ids` - The user ids

##### Version (proxmox_virtual_environment_version)

###### Arguments
This data source doesn't accept arguments.

###### Attributes
* `keyboard_layout` - The keyboard layout
* `release` - The release number
* `repository_id` - The repository id
* `version` - The version string

### Resources

#### Virtual Environment

##### File (proxmox_virtual_environment_file)

###### Arguments
* `content_type` - (Optional) The content type
    * `backup`
    * `iso`
    * `snippets`
    * `vztmpl`
* `datastore_id` - (Required) The datastore id
* `node_name` - (Required) The node name
* `source_file` - (Optional) The source file (conflicts with `source_raw`)
    * `checksum` - (Optional) The SHA256 checksum of the source file
    * `file_name` - (Optional) The file name to use instead of the source file name
    * `insecure` - (Optional) Whether to skip the TLS verification step for HTTPS sources (defaults to `false`)
    * `path` - (Required) A path to a local file or a URL
* `source_raw` - (Optional) The raw source (conflicts with `source_file`)
    * `data` - (Required) The raw data
    * `file_name` - (Required) The file name
    * `resize` - (Optional) The number of bytes to resize the file to

###### Attributes
* `file_modification_date` - The file modification date (RFC 3339)
* `file_name` - The file name
* `file_size` - The file size in bytes
* `file_tag` - The file tag

###### Notes

The Proxmox VE API endpoint for file uploads does not support chunked transfer encoding, which means that we must first store the source file as a temporary file locally before uploading it.

You must ensure that you have at least `Size-in-MB * 2 + 1` MB of storage space available (twice the size plus overhead because a multipart payload needs to be created as another temporary file).

##### Group (proxmox_virtual_environment_group)

###### Arguments
* `acl` - (Optional) The access control list (multiple blocks supported)
    * `path` - The path
    * `propagate` - Whether to propagate to child paths
    * `role_id` - The role id
* `comment` - (Optional) The group comment
* `group_id` - (Required) The group id

###### Attributes
* `members` - The group members as a list with `username@realm` entries

##### Pool (proxmox_virtual_environment_pool)

###### Arguments
* `comment` - (Optional) The pool comment
* `pool_id` - (Required) The pool id

###### Attributes
* `members` - The pool members
    * `datastore_id` - The datastore id
    * `id` - The member id
    * `node_name` - The node name
    * `type` - The member type
    * `vm_id` - The virtual machine id

##### Role (proxmox_virtual_environment_role)

###### Arguments
* `privileges` - (Required) The role privileges
* `role_id` - (Required) The role id

###### Attributes
This resource doesn't expose any additional attributes.

##### User (proxmox_virtual_environment_user)

###### Arguments
* `acl` - (Optional) The access control list (multiple blocks supported)
    * `path` - The path
    * `propagate` - Whether to propagate to child paths
    * `role_id` - The role id
* `comment` - (Optional) The user comment
* `email` - (Optional) The user's email address
* `enabled` - (Optional) Whether the user account is enabled
* `expiration_date` - (Optional) The user account's expiration date (RFC 3339)
* `first_name` - (Optional) The user's first name
* `groups` - (Optional) The user's groups
* `keys` - (Optional) The user's keys
* `last_name` - (Optional) The user's last name
* `password` - (Required) The user's password
* `user_id` - (Required) The user id

###### Attributes
This resource doesn't expose any additional attributes.

##### VM (proxmox_virtual_environment_vm)

###### Arguments
* `agent` - (Optional) The QEMU agent configuration
    * `enabled` - (Optional) Whether to enable the QEMU agent (defaults to `false`)
    * `trim` - (Optional) Whether to enable the FSTRIM feature in the QEMU agent (defaults to `false`)
    * `type` - (Optional) The QEMU agent interface type (defaults to `virtio`)
        * `isa` - ISA Serial Port
        * `virtio` - VirtIO (paravirtualized)
* `cdrom` - (Optional) The CDROM configuration
    * `enabled` - (Optional) Whether to enable the CDROM drive (defaults to `false`)
    * `file_id` - (Optional) A file ID for an ISO file (defaults to `cdrom` as in the physical drive)
* `cloud_init` - (Optional) The cloud-init configuration (conflicts with `cdrom`)
    * `dns` - (Optional) The DNS configuration
        * `domain` - (Optional) The DNS search domain
        * `server` - (Optional) The DNS server
    * `ip_config` - (Optional) The IP configuration (one block per network device)
        * `ipv4` - (Optional) The IPv4 configuration
            * `address` - (Optional) The IPv4 address (use `dhcp` for autodiscovery)
            * `gateway` - (Optional) The IPv4 gateway (must be omitted when `dhcp` is used as the address)
        * `ipv6` - (Optional) The IPv4 configuration
            * `address` - (Optional) The IPv6 address (use `dhcp` for autodiscovery)
            * `gateway` - (Optional) The IPv6 gateway (must be omitted when `dhcp` is used as the address)
    * `user_account` - (Required) The user account configuration (conflicts with `user_data_file_id`)
        * `keys` - (Required) The SSH keys
        * `password` - (Optional) The SSH password
        * `username` - (Required) The SSH username
    * `user_data_file_id` - (Optional) The ID of a file containing custom user data (conflicts with `user_account`)
* `cpu` - (Optional) The CPU configuration
    * `cores` - (Optional) The number of CPU cores (defaults to `1`)
    * `flags` - (Optional) The CPU flags
        * `+aes`/`-aes` - Activate AES instruction set for HW acceleration
        * `+amd-no-ssb`/`-amd-no-ssb` - Notifies guest OS that host is not vulnerable for Spectre on AMD CPUs
        * `+amd-ssbd`/`-amd-ssbd` - Improves Spectre mitigation performance with AMD CPUs, best used with "virt-ssbd"
        * `+hv-evmcs`/`-hv-evmcs` - Improve performance for nested virtualization (only supported on Intel CPUs)
        * `+hv-tlbflush`/`-hv-tlbflush` - Improve performance in overcommitted Windows guests (may lead to guest BSOD on old CPUs)
        * `+ibpb`/`-ibpb` - Allows improved Spectre mitigation on AMD CPUs
        * `+md-clear`/`-md-clear` - Required to let the guest OS know if MDS is mitigated correctly
        * `+pcid`/`-pcid` - Meltdown fix cost reduction on Westmere, Sandy- and Ivy Bridge Intel CPUs
        * `+pdpe1gb`/`-pdpe1gb` - Allows guest OS to use 1 GB size pages, if host HW supports it
        * `+spec-ctrl`/`-spec-ctrl` - Allows improved Spectre mitigation with Intel CPUs
        * `+ssbd`/`-ssbd` - Protection for "Speculative Store Bypass" for Intel models
        * `+virt-ssbd`/`-virt-ssbd` - Basis for "Speculative Store Bypass" protection for AMD models
    * `hotplugged` - (Optional) The number of hotplugged vCPUs (defaults to `0`)
    * `sockets` - (Optional) The number of CPU sockets (defaults to `1`)
    * `type` - (Optional) The emulated CPU type (defaults to `qemu64`)
        * `486` - Intel 486
        * `Broadwell`/`Broadwell-IBRS`/`Broadwell-noTSX`/`Broadwell-noTSX-IBRS` - Intel Core Processor (Broadwell, 2014)
        * `Cascadelake-Server` - Intel Xeon 32xx/42xx/52xx/62xx/82xx/92xx (2019)
        * `Conroe` - Intel Celeron_4x0 (Conroe/Merom Class Core 2, 2006)
        * `EPYC`/`EPYC-IBPB` - AMD EPYC Processor (2017)
        * `Haswell`/`Haswell-IBRS`/`Haswell-noTSX`/`Haswell-noTSX-IBRS` - Intel Core Processor (Haswell, 2013)
        * `IvyBridge`/`IvyBridge-IBRS` - Intel Xeon E3-12xx v2 (Ivy Bridge, 2012)
        * `KnightsMill` - Intel Xeon Phi 72xx (2017)
        * `Nehalem`/`Nehalem-IBRS` - Intel Core i7 9xx (Nehalem Class Core i7, 2008)
        * `Opteron_G1` - AMD Opteron 240 (Gen 1 Class Opteron, 2004)
        * `Opteron_G2` - AMD Opteron 22xx (Gen 2 Class Opteron, 2006)
        * `Opteron_G3` - AMD Opteron 23xx (Gen 3 Class Opteron, 2009)
        * `Opteron_G4` - AMD Opteron 62xx class CPU (2011)
        * `Opteron_G5` - AMD Opteron 63xx class CPU (2012)
        * `Penryn` - Intel Core 2 Duo P9xxx (Penryn Class Core 2, 2007)
        * `SandyBridge`/`SandyBridge-IBRS` - Intel Xeon E312xx (Sandy Bridge, 2011)
        * `Skylake-Client`/`Skylake-Client-IBRS` - Intel Core Processor (Skylake, 2015)
        * `Skylake-Server`/`Skylake-Server-IBRS` - Intel Xeon Processor (Skylake, 2016)
        * `Westmere`/`Westmere-IBRS` - Intel Westmere E56xx/L56xx/X56xx (Nehalem-C, 2010)
        * `athlon` - AMD Athlon
        * `core2duo` - Intel Core 2 Duo
        * `coreduo` - Intel Core Duo
        * `host` - Host passthrough
        * `kvm32`/`kvm64` - Common KVM processor (32 & 64 bit variants)
        * `max` - Maximum amount of features from host CPU
        * `pentium` - Intel Pentium (1993)
        * `pentium2` - Intel Pentium 2 (1997-1999)
        * `pentium3` - Intel Pentium 3 (1999-2001)
        * `phenom` - AMD Phenom (2010)
        * `qemu32`/`qemu64` - QEMU Virtual CPU version 2.5+ (32 & 64 bit variants)
    * `units` - (Optional) The CPU units (defaults to `1024`)
* `description` - (Optional) The description
* `disk` - (Optional) The disk configuration (multiple blocks supported)
    * `datastore_id` - (Optional) The ID of the datastore to create the disk in (defaults to `local-lvm`)
    * `file_format` - (Optional) The file format (defaults to `qcow2`)
        * `qcow2` - QEMU Disk Image v2
        * `raw` - Raw Disk Image
        * `vmdk` - VMware Disk Image
    * `file_id` - (Optional) The file ID for a disk image (experimental - might cause high CPU utilization during import, especially with large disk images)
    * `size` - (Optional) The disk size in gigabytes (defaults to `8`)
    * `speed` - (Optional) The speed limits
        * `read` - (Optional) The maximum read speed in megabytes per second
        * `read_burstable` - (Optional) The maximum burstable read speed in megabytes per second
        * `write` - (Optional) The maximum write speed in megabytes per second
        * `write_burstable` - (Optional) The maximum burstable write speed in megabytes per second
* `keyboard_layout` - (Optional) The keyboard layout (defaults to `en-us`)
    * `da` - Danish
    * `de` - German
    * `de-ch` - Swiss German
    * `en-gb` - British English
    * `en-us` - American English
    * `es` - Spanish
    * `fi` - Finnish
    * `fr` - French
    * `fr-be` - Belgian French
    * `fr-ca` - French Canadian
    * `fr-ch` - Swish French
    * `hu` - Hungarian
    * `is` - Icelandic
    * `it` - Italian
    * `ja` - Japanese
    * `lt` - Lithuanian
    * `mk` - Macedonian
    * `nl` - Dutch
    * `no` - Norwegian
    * `pl` - Polish
    * `pt` - Portuguese
    * `pt-br` - Brazilian Portuguese
    * `sl` - Slovenian
    * `sv` - Swedish
    * `tr` - Turkish
* `memory` - (Optional) The memory configuration
    * `dedicated` - (Optional) The dedicated memory in megabytes (defaults to `512`)
    * `floating` - (Optional) The floating memory in megabytes (defaults to `0`)
    * `shared` - (Optional) The shared memory in megabytes (defaults to `0`)
* `name` - (Optional) The name
* `network_device` - (Optional) The network device configuration (multiple blocks supported)
    * `bridge` - (Optional) The name of the network bridge (defaults to `vmbr0`)
    * `enabled` - (Optional) Whether to enable the network device (defaults to `true`)
    * `mac_address` - (Optional) The MAC address
    * `model` - (Optional) The network device model (defaults to `virtio`)
        * `e1000` - Intel E1000
        * `rtl8139` - Realtek RTL8139
        * `virtio` - VirtIO (paravirtualized)
        * `vmxnet3` - VMware vmxnet3
    * `rate_limit` - (Optional) The rate limit in megabytes per second
    * `vlan_id` - (Optional) The VLAN identifier
* `node_name` - (Required) The name of the node to assign the virtual machine to
* `os_type` - (Optional) The OS type (defaults to `other`)
    * `l24` - Linux Kernel 2.4
    * `l26` - Linux Kernel 2.6 - 5.X
    * `other` - Unspecified OS
    * `solaris` - OpenIndiania, OpenSolaris og Solaris Kernel
    * `w2k` - Windows 2000
    * `w2k3` - Windows 2003
    * `w2k8` - Windows 2008
    * `win7` - Windows 7
    * `win8` - Windows 8, 2012 or 2012 R2
    * `win10` - Windows 10 or 2016
    * `wvista` - Windows Vista
    * `wxp` - Windows XP
* `pool_id` - (Optional) The ID of a pool to assign the virtual machine to
* `started` - (Optional) Whether to start the virtual machine (defaults to `true`)
* `tablet_device` - (Optional) Whether to enable the USB tablet device (defaults to `true`)
* `vga` - (Optional) The VGA configuration
    * `enabled` - (Optional) Whether to enable the VGA device (defaults to `true`)
    * `memory` - (Optional) The VGA memory in megabytes (defaults to `16`)
    * `type` - (Optional) The VGA type (defaults to `std`)
        * `cirrus` - Cirrus (deprecated since QEMU 2.2)
        * `qxl` - SPICE
        * `qxl2` - SPICE Dual Monitor
        * `qxl3` - SPICE Triple Monitor
        * `qxl4` - SPICE Quad Monitor
        * `serial0` - Serial Terminal 0
        * `serial1` - Serial Terminal 1
        * `serial2` - Serial Terminal 2
        * `serial3` - Serial Terminal 3
        * `std` - Standard VGA
        * `virtio` - VirtIO-GPU
        * `vmware` - VMware Compatible
* `vm_id` - (Optional) The ID

###### Attributes
* `ipv4_addresses` - The IPv4 addresses per network interface published by the QEMU agent (empty list when `agent.enabled` is `false`)
* `ipv6_addresses` - The IPv6 addresses per network interface published by the QEMU agent (empty list when `agent.enabled` is `false`)
* `mac_addresses` - The MAC addresses published by the QEMU agent with fallback to the network device configuration, if the agent is disabled
* `network_interface_names` - The network interface names published by the QEMU agent (empty list when `agent.enabled` is `false`)

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
