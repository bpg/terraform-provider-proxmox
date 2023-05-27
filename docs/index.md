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
section below for all the available arguments in the `ssh` block.

### API Token authentication

API Token authentication can be used to authenticate with the Proxmox API
without the need to provide a password. In combination with the `ssh` block,
this allows for a fully password-less authentication.

To create an API Token, log in to the Proxmox web interface, and navigate to
`Datacenter` > `Permissions` > `API Tokens`. Click on `Add` to create a new
token. You can then use the `token` field in the `provider` block to provide
the token:

```terraform
provider "proxmox" {
  endpoint = var.virtual_environment_endpoint
  api_token = "root@pam!for-terraform-provider=xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx"
  insecure = true
  ssh {
    agent = true
    username = "root"
  }
}
```

> Note1: The `username` field in the `ssh` block is required when using API
> Token authentication. This is because provider needs to know which user to use
> for the SSH connection.

> Note2: Not all Proxmox API operations are supported via API Token. You may see errors like
> `error creating container: received an HTTP 403 response - Reason: Permission check failed (changing feature flags for privileged container is only allowed for root@pam)` or `error creating VM: received an HTTP 500 response - Reason: only root can set 'arch' config`
> when using API Token authentication, even when `Administrator` role or the `root@pam` user is used with the token.

## Argument Reference

In addition
to [generic provider arguments](https://www.terraform.io/docs/configuration/providers.html) (
e.g. `alias` and `version`), the following arguments are supported in the
Proxmox `provider` block:

- `endpoint` - (Required) The endpoint for the Proxmox Virtual Environment
  API (can also be sourced from `PROXMOX_VE_ENDPOINT`). Usually this is
  `https://<your-cluster-endpoint>:8006/`.
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
  a
  block, whose fields are documented below.
    - `username` - (Optional) The username to use for the SSH connection.
      Defaults to the username used for the Proxmox API connection. Can also be
      sourced from `PROXMOX_VE_SSH_USERNAME`.
    - `password` - (Optional) The password to use for the SSH connection.
      Defaults to the password used for the Proxmox API connection. Can also be
      sourced from `PROXMOX_VE_SSH_PASSWORD`.
    - `agent` - (Optional) Whether to use the SSH agent for the SSH
      authentication. Defaults to `false`. Can also be sourced
      from `PROXMOX_VE_SSH_AGENT`.
    - `agent_socket` - (Optional) The path to the SSH agent socket.
      Defaults to the value of the `SSH_AUTH_SOCK` environment variable. Can
      also be sourced from `PROXMOX_VE_SSH_AUTH_SOCK`.
