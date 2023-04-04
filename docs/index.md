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
