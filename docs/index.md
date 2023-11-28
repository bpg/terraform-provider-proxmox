---
layout: home
title: Introduction
nav_order: 1
---

# Proxmox Provider

This provider for [Terraform](https://www.terraform.io/) is used for interacting
with resources supported by [Proxmox](https://www.proxmox.com/en/). The provider
needs to be configured with the proper endpoints and credentials before it can
be used.

Use the navigation to the left to read about the available resources.

## Example Usage

```terraform
provider "proxmox" {
  endpoint = "https://10.0.0.2:8006/"
  username = "root@pam"
  password = "the-password-set-during-installation-of-proxmox-ve"
  insecure = true
  tmp_dir  = "/var/tmp"
}
```

## Authentication

The Proxmox provider offers a flexible means of providing credentials for
authentication. The following methods are supported, in this order, and
explained below:

- Static credentials
- Environment variables

### Static credentials

> Warning: Hard-coding credentials into any Terraform configuration is not
> recommended, and risks secret leakage should this file ever be committed to a
> public version control system.

Static credentials can be provided by adding a `username` and `password` in-line
in the Proxmox provider block:

```terraform
provider "proxmox" {
  username = "username@realm"
  password = "a-strong-password"
}
```

### Environment variables

You can provide your credentials via the `PROXMOX_VE_USERNAME`
and `PROXMOX_VE_PASSWORD`, environment variables, representing your Proxmox
username, realm and password, respectively:

```terraform
provider "proxmox" {}
```

Usage:

```sh
export PROXMOX_VE_USERNAME="username@realm"
export PROXMOX_VE_PASSWORD="a-strong-password"
terraform plan
```

### SSH connection

The Proxmox provider can connect to a Proxmox node via SSH. This is used in
the `proxmox_virtual_environment_vm` or `proxmox_virtual_environment_file`
resource to execute commands on the node to perform actions that are not
supported by Proxmox API. For example, to import VM disks, or to uploading
certain type of resources, such as snippets.

The SSH connection configuration is provided via the optional `ssh` block in
the `provider` block:

```terraform
provider "proxmox" {
  endpoint = "https://10.0.0.2:8006/"
  username = "username@realm"
  password = "a-strong-password"
  insecure = true
  ssh {
    agent = true
  }
}
```

If no `ssh` block is provided, the provider will attempt to connect to the
target node using the credentials provided in the `username` and `password`
fields.
Note that the target node is identified by the `node` argument in the resource,
and may be different from the Proxmox API endpoint. Please refer to the
"Argument Reference" section below for all the available arguments in the `ssh`
block.

#### Node IP address used for SSH connection

In order to make the SSH connection, the provider needs to know the IP address
of the target node. The provider will attempt to resolve the
node name to an IP address using Proxmox API to enumerate the node network
interfaces, and use the first one that is not a loopback interface. In some
cases this may not be the desired behavior, for example when the node has
multiple network interfaces, and the one that should be used for SSH is not the
first one.

To override the node IP address used for SSH connection, you can use the
optional `node` blocks in the `ssh` block. For example:

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

### API Token authentication

API Token authentication can be used to authenticate with the Proxmox API
without the need to provide a password. In combination with the `ssh` block,
this allows for a fully password-less authentication.

To create an API Token, log in to the Proxmox web interface, and navigate to
`Datacenter` > `Permissions` > `API Tokens`. Click on `Add` to create a new
token. You can then use the `api_token` field in the `provider` block to provide
the token. `api_token` can also be sourced from `PROXMOX_VE_API_TOKEN`
environment variable. The token authentication is taking precedence over the
password authentication.

```terraform
provider "proxmox" {
  endpoint  = var.virtual_environment_endpoint
  api_token = "root@pam!for-terraform-provider=xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx"
  insecure  = true
  ssh {
    agent    = true
    username = "root"
  }
}
```

-> **Note:** Note1: The `username` field in the `ssh` block (or alternatively a username in `PROXMOX_VE_USERNAME` or `PROXMOX_VE_SSH_USERNAME` environment variable) is required when using API Token authentication. This is because the provider needs to know which user to use for the SSH connection.

-> **Note:** Note2: Not all Proxmox API operations are supported via API Token. You may see errors like `error creating container: received an HTTP 403 response - Reason: Permission check failed (changing feature flags for privileged container is only allowed for root@pam)` or `error creating VM: received an HTTP 500 response - Reason: only root can set 'arch' config` when using API Token authentication, even when `Administrator` role or the `root@pam` user is used with the token. The workaround is to use password authentication for those operations.

### Temporary directory

Using `proxmox_virtual_environment_file` with `.iso` files or disk images can require
large amount of space in the temporary directory of the computer running terraform.

Consider pointing `tmp_dir` to a directory with enough space, especially if the default
temporary directory is limited by the system memory (e.g. `tmpfs` mounted
on `/tmp`).

## Argument Reference

In addition
to [generic provider arguments](https://www.terraform.io/docs/configuration/providers.html) (
e.g. `alias` and `version`), the following arguments are supported in the
Proxmox `provider` block:

- `endpoint` - (Required) The endpoint for the Proxmox Virtual Environment
  API (can also be sourced from `PROXMOX_VE_ENDPOINT`). Usually this is
  `https://<your-cluster-endpoint>:8006/`. **Do not** include `/api2/json` at
  the end.
- `insecure` - (Optional) Whether to skip the TLS verification step (can
  also be sourced from `PROXMOX_VE_INSECURE`). If omitted, defaults
  to `false`.
- `otp` - (Optional) The one-time password for the Proxmox Virtual
  Environment API (can also be sourced from `PROXMOX_VE_OTP`).
- `password` - (Required) The password for the Proxmox Virtual Environment
  API (can also be sourced from `PROXMOX_VE_PASSWORD`).
- `username` - (Required) The username and realm for the Proxmox Virtual
  Environment API (can also be sourced from `PROXMOX_VE_USERNAME`). For
  example, `root@pam`.
- `api_token` - (Optional) The API Token for the Proxmox Virtual
  Environment API (can also be sourced from `PROXMOX_VE_API_TOKEN`). For
  example, `root@pam!for-terraform-provider=xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx`.
- `ssh` - (Optional) The SSH connection configuration to a Proxmox node. This is
  a block, whose fields are documented below.
  - `username` - (Optional) The username to use for the SSH connection.
      Defaults to the username used for the Proxmox API connection. Can also be
      sourced from `PROXMOX_VE_SSH_USERNAME`. Required when using API Token.
  - `password` - (Optional) The password to use for the SSH connection.
      Defaults to the password used for the Proxmox API connection. Can also be
      sourced from `PROXMOX_VE_SSH_PASSWORD`.
  - `agent` - (Optional) Whether to use the SSH agent for the SSH
      authentication. Defaults to `false`. Can also be sourced
      from `PROXMOX_VE_SSH_AGENT`.
  - `agent_socket` - (Optional) The path to the SSH agent socket.
      Defaults to the value of the `SSH_AUTH_SOCK` environment variable. Can
      also be sourced from `PROXMOX_VE_SSH_AUTH_SOCK`.
  - `node` - (Optional) The node configuration for the SSH connection. Can be
      specified multiple times to provide configuration fo multiple nodes.
    - `name` - (Required) The name of the node.
    - `address` - (Required) The IP address of the node.
    - `port` - (Optional) SSH port of the node. Defaults to 22.
- `tmp_dir` - (Optional) Use custom temporary directory. (can also be sourced from `PROXMOX_VE_TMPDIR`)
