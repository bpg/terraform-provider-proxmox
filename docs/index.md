---
layout: home
title: Introduction
nav_order: 1
---

# Proxmox Provider

This provider for [Terraform](https://www.terraform.io/) is used for interacting with resources supported by [Proxmox](https://www.proxmox.com/en/). The provider needs to be configured with the proper endpoints and credentials before it can be used.

Use the navigation to the left to read about the available resources.

## Authentication

The Proxmox provider offers a flexible means of providing credentials for authentication. Static credentials can be provided to the `proxmox` block through either a `api_token` or a combination of `username` and `password` arguments.

> Warning: Hard-coding credentials into any Terraform configuration is not recommended. The practice risks secret leakage should the configuration ever be committed to a public version control system. See the [API Token authentication](#api-token-authentication) for the best approach.

```terraform
provider "proxmox" {
  username = "username@realm"
  password = "a-strong-password"
}
```

Instead of using static arguments, credentials can be handled through the use of environment variables. See the [Argument Reference](#argument-reference) section for the names and use cases.

```terraform
provider "proxmox" {}
```

Usage:

```sh
export PROXMOX_VE_USERNAME="username@realm"
export PROXMOX_VE_PASSWORD="a-strong-password"
terraform plan
```

## SSH connection

The Proxmox provider can connect to a Proxmox node via SSH. This is used in the `proxmox_virtual_environment_vm` or `proxmox_virtual_environment_file` resource to execute commands on the node to perform actions that are not supported by Proxmox API. For example, to import VM disks, or to uploading certain type of resources, such as snippets.

The SSH connection configuration is provided via the optional `ssh` block in the `provider` block:

```terraform
provider "proxmox" {

  // ...
  ssh {
    agent = true
  }
}
```

If no `ssh` block is provided, the provider will attempt to connect to the target node using the credentials provided in the `username` and `password` fields. Note that the target node is identified by the `node` argument in the resource, and may be different from the Proxmox API endpoint. Please refer to the [Argument Reference](#argument-reference) section to view the available arguments of the `ssh` block.

## Node IP address used for SSH connection

In order to make the SSH connection, the provider needs to be able to resolve the target node to an IP. When the target node is represented as a name, the provider will enumerate the node's interfaces via the Proxmox API; using the first one that is not a loopback device. In some cases this may not be the desired behavior, for example when the node has multiple network interfaces, and the one that should be used for SSH is not the first one.

To override the node IP address used for SSH connection, you can use the optional `node` blocks in the `ssh` block. For example:

```terraform
provider "proxmox" {
  // ...
  ssh {
    // ...
    node {
      name    = "pve1"
      address = "192.168.10.1"
    }
    node {
      name    = "pve2"
      address = "192.168.10.2"
    }
  }
}

```

## API Token Authentication

API Token authentication can be used to authenticate against the Proxmox API without the need to provide a password. In combination with the `ssh` block, this allows for a fully password-less authentication.

To set this up, SSH into the Proxmox cluster or host in order to accomplish the following (or use the GUI):

Create a user:

```sh
sudo pveum user add terraform-prov@pve
```

Create a role for the user:

```sh
sudo pveum role add TerraformProv -privs "Datastore.Allocate Datastore.AllocateSpace Datastore.AllocateTemplate Datastore.Audit Pool.Allocate Sys.Audit Sys.Console Sys.Modify SDN.Use VM.Allocate VM.Audit VM.Clone VM.Config.CDROM VM.Config.Cloudinit VM.Config.CPU VM.Config.Disk VM.Config.HWType VM.Config.Memory VM.Config.Network VM.Config.Options VM.Migrate VM.Monitor VM.PowerMgmt"
```

Assign the role to the previously created user:

```sh
sudo pveum aclmod / -user terraform-prov@pve -role TerraformProv
```

Create an API token for the user:

```
sudo pveum user token add terraform-prov@pve terraform --privsep=0
```

Generating the token will output a table containing the token's ID and secret which are meant to be concatenated into a single string for use with either the `api_token` field of the `provider` block (fine for testing but should be avoided) or sourced from the `PROXMOX_VE_API_TOKEN` environment variable.

Refer to the upstream docs as needed for additional details concerning [PVE User Management](https://pve.proxmox.com/wiki/User_Management).

```terraform
provider "proxmox" {
  api_token = "terraform-prov@pve!terraform=xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx"
  // ...
}
```

## Temporary directory

Using `proxmox_virtual_environment_file` with `.iso` files or disk images can require large amount of space in the temporary directory of the computer running terraform.

Consider pointing `tmp_dir` to a directory with enough space, especially if the default temporary directory is limited by the system memory (e.g. `tmpfs` mounted
on `/tmp`).

## Argument Reference

In addition to [generic provider arguments](https://www.terraform.io/docs/configuration/providers.html) ( e.g. `alias` and `version`), the following arguments are supported in the Proxmox `provider` block:

- `api_token` - (Optional) The API Token for the Proxmox Virtual Environment API (better to avoid this and source from the `PROXMOX_VE_API_TOKEN` environment variable instead). Either way, the string will be formatted as `terraform-prov@pve!terraform=xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx`.
- `endpoint` - (Required) The endpoint for the Proxmox Virtual Environment API (can also be sourced from `PROXMOX_VE_ENDPOINT`). Usually this is `https://<your-cluster-endpoint>:8006/`.
- `insecure` - (Optional) Whether to skip the TLS verification step (can also be sourced from `PROXMOX_VE_INSECURE`). If omitted, defaults to `false`.
- `otp` - (Optional) The one-time password for the Proxmox Virtual Environment API (can also be sourced from `PROXMOX_VE_OTP`).
- `password` - (Required) The password for the Proxmox Virtual Environment API (can also be sourced from `PROXMOX_VE_PASSWORD`).
- `ssh` - (Optional) The SSH connection configuration to a Proxmox node. This is a block, whose fields are documented below.
  - `username` - (Optional) The username to use for the SSH connection. Defaults to the username used for the Proxmox API connection. Can also be sourced from `PROXMOX_VE_SSH_USERNAME`. At this time, this must be set to 'root' if provider.username is set otherwise.
  - `password` - (Optional) The password to use for the SSH connection. Defaults to the password used for the Proxmox API connection. Can also be sourced from `PROXMOX_VE_SSH_PASSWORD`.
  - `agent` - (Optional) Whether to use the SSH agent for the SSH authentication. Defaults to `false`. Can also be sourced from `PROXMOX_VE_SSH_AGENT`.
  - `agent_socket` - (Optional) The path to the SSH agent socket. Defaults to the value of the `SSH_AUTH_SOCK` environment variable. Can also be sourced from `PROXMOX_VE_SSH_AUTH_SOCK`.
  - `node` - (Optional) The node configuration for the SSH connection. Can be specified multiple times to provide configuration for multiple nodes.
    - `name` - (Required) The name of the node.
    - `address` - (Required) The IP address of the node.
    - `port` - (Optional) SSH port of the node. Defaults to 22.
- `tmp_dir` - (Optional) Use custom temporary directory. (can also be sourced from `PROXMOX_VE_TMPDIR`)
- `username` - (Optional) The username and realm for the Proxmox Virtual Environment API (can also be sourced from `PROXMOX_VE_USERNAME`).
