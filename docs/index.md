---
layout: home
title: "Provider: Proxmox Virtual Environment"
---

# Proxmox Provider

This provider for [Terraform](https://www.terraform.io/) / [OpenTofu](https://opentofu.org/) is used for interacting with resources supported by [Proxmox VE](https://www.proxmox.com/en/products/proxmox-virtual-environment/overview).
The provider needs to be configured with the proper endpoint and credentials before it can be used.

Use the navigation to the left to read about the available resources.

## Getting Started

To use this provider, you only need:

1. **API access** to your Proxmox VE server (endpoint URL + username/password or API token)
2. That's it for most use cases!

-> **SSH access is optional.** Most resources work entirely through the Proxmox API. See [When is SSH Required?](#when-is-ssh-required) for specific cases that need SSH.

## Table of Contents

- [Example Usage](#example-usage)
- [Authentication](#authentication)
  - [Authentication Methods Comparison](#authentication-methods-comparison)
  - [Quick Examples](#quick-examples)
  - [Security Best Practices](#security-best-practices)
  - [Environment Variables](#environment-variables)
  - [API Token Authentication](#api-token-authentication)
  - [Pre-Authentication](#pre-authentication)
- [SSH Connection](#ssh-connection) *(optional)*
  - [When is SSH Required?](#when-is-ssh-required)
  - [SSH Configuration](#ssh-configuration)
  - [SSH Agent](#ssh-agent)
  - [SSH Private Key](#ssh-private-key)
  - [SSH User](#ssh-user)
  - [Node IP address used for SSH connection](#node-ip-address-used-for-ssh-connection)
  - [SSH Connection via SOCKS5 Proxy](#ssh-connection-via-socks5-proxy)
- [VM and Container ID Assignment](#vm-and-container-id-assignment)
- [Temporary Directory](#temporary-directory)
- [Environment Variables Summary](#environment-variables-summary)
- [Argument Reference](#argument-reference)

## Example Usage

**Minimal configuration (no SSH):**

```hcl
provider "proxmox" {
  endpoint = "https://10.0.0.2:8006/"

  # TODO: use terraform variable or remove the line, and use PROXMOX_VE_USERNAME environment variable
  username = "root@pam"
  # TODO: use terraform variable or remove the line, and use PROXMOX_VE_PASSWORD environment variable
  password = "the-password-set-during-installation-of-proxmox-ve"

  # because self-signed TLS certificate is in use
  insecure = true
}
```

**With SSH access (only if needed for snippets or certain file uploads):**

```hcl
provider "proxmox" {
  endpoint = "https://10.0.0.2:8006/"
  username = "root@pam"
  password = "the-password-set-during-installation-of-proxmox-ve"
  insecure = true

  ssh {
    agent = true
    # username = "root"  # required when using api_token
  }
}
```

## Authentication

The provider supports three authentication methods (in order of precedence):

1. **API Token** — recommended for production and CI/CD
2. **Auth Ticket** — for automated scripts with TOTP support
3. **Username/Password** — simplest, good for development

!> Hard-coding credentials into any Terraform configuration is not recommended. Use environment variables or a `.tfvars` file (add to `.gitignore`) instead.

### Authentication Methods Comparison

| Method                                                                                   | Use Case             | Pros                                                              | Cons                                                              | Security Level |
|------------------------------------------------------------------------------------------|----------------------|-------------------------------------------------------------------|-------------------------------------------------------------------|----------------|
| [API Token](#api-token-authentication)                                                   | Production, CI/CD    | - No password needed<br>- Fine-grained permissions<br>- Revocable | - Some operations not supported<br>- Requires SSH username config | High           |
| [Auth Ticket](#pre-authentication-or-passing-an-authentication-ticket-into-the-provider) | Automated scripts    | - Short-lived<br>- No password storage<br>- TOTP support          | - More complex setup<br>- Needs periodic renewal                  | High           |
| Username/Password                                                                        | Development, Testing | - Full API support<br>- Simple setup                              | - Password in config/env<br>- Not revocable individually          | Medium         |

### Quick Examples

Here are examples for each authentication method:

**API Token (Recommended for Production):**

```hcl
provider "proxmox" {
  endpoint  = "https://10.0.0.2:8006/"
  api_token = "terraform@pve!provider=xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx"
}
```

**Username/Password (Development/Testing):**

```hcl
provider "proxmox" {
  endpoint = "https://10.0.0.2:8006/"
  insecure = true
  username = "username@realm"
  password = "a-strong-password"
}
```

**Auth Ticket (Automated Scripts):**

```hcl
provider "proxmox" {
  endpoint              = "https://10.0.0.2:8006/"
  auth_ticket          = "PVE:username@realm:12345678::some_base64_payload=="
  csrf_prevention_token = "12345678:some_blob"
}
```

A better approach is to extract these values into Terraform variables and reference them instead:

```hcl
provider "proxmox" {
  endpoint = var.virtual_environment_endpoint
  
  # Choose one authentication method:
  api_token = var.virtual_environment_api_token
  # OR
  username  = var.virtual_environment_username
  password  = var.virtual_environment_password
  # OR
  auth_ticket           = var.virtual_environment_auth_ticket
  csrf_prevention_token = var.virtual_environment_csrf_prevention_token
}
```

The variable values can be provided via a separate `.tfvars` file (add it to `.gitignore`).
See the [Terraform documentation](https://developer.hashicorp.com/terraform/language/values/variables#input-variables) for more information.

### Security Best Practices

- **Use API tokens** in production — they're revocable and support fine-grained permissions
- **Never commit credentials** to version control — use environment variables or `.tfvars` files (in `.gitignore`)
- **Use HTTPS with valid certificates** — only set `insecure = true` in development environments
- **Apply least privilege** — create tokens/users with minimal required permissions
- **Rotate credentials** regularly

### Environment Variables

Credentials can also be provided via environment variables instead of static arguments.
For example:

```hcl
provider "proxmox" {
  endpoint = "https://10.0.0.2:8006/"
}
```

```sh
export PROXMOX_VE_USERNAME="username@realm"
export PROXMOX_VE_PASSWORD='a-strong-password'
terraform plan
```

See the [Argument Reference](#argument-reference) section for the supported variable names and use cases.

### API Token Authentication

API tokens allow password-less authentication with the Proxmox API. If you already have a token, use it like this:

```hcl
provider "proxmox" {
  endpoint  = "https://10.0.0.2:8006/"
  api_token = "user@realm!tokenid=xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx"
}
```

#### Creating an API Token on the Proxmox Server

You can create an API Token via the Proxmox UI or the command line on the Proxmox host:

- Create a user:

    ```sh
    pveum user add terraform@pve
    ```

- Create a role for the user (you can skip this step if you want to use any of the existing roles):

    ```sh
    pveum role add Terraform -privs "Realm.AllocateUser, VM.PowerMgmt, VM.GuestAgent.Unrestricted, Sys.Console, Sys.Audit, Sys.AccessNetwork, VM.Config.Cloudinit, VM.Replicate, Pool.Allocate, SDN.Audit, Realm.Allocate, SDN.Use, Mapping.Modify, VM.Config.Memory, VM.GuestAgent.FileSystemMgmt, VM.Allocate, SDN.Allocate, VM.Console, VM.Clone, VM.Backup, Datastore.AllocateTemplate, VM.Snapshot, VM.Config.Network, Sys.Incoming, Sys.Modify, VM.Snapshot.Rollback, VM.Config.Disk, Datastore.Allocate, VM.Config.CPU, VM.Config.CDROM, Group.Allocate, Datastore.Audit, VM.Migrate, VM.GuestAgent.FileWrite, Mapping.Use, Datastore.AllocateSpace, Sys.Syslog, VM.Config.Options, Pool.Audit, User.Modify, VM.Config.HWType, VM.Audit, Sys.PowerMgmt, VM.GuestAgent.Audit, Mapping.Audit, VM.GuestAgent.FileRead, Permissions.Modify"
    ```

  ~> The list of available privileges has changed in PVE 9.0. The above list is only an example (and likely too permissive for most use cases). Please review and adjust to your needs.
  Refer to the [privileges documentation](https://pve.proxmox.com/pve-docs/pveum.1.html#_privileges) for more details.

- Assign the role to the previously created user:

    ```sh
    pveum aclmod / -user terraform@pve -role Terraform
    ```

- Create an API token for the user:

    ```sh
    pveum user token add terraform@pve provider --privsep=0
    ```

    -> Make sure you copy the token value, as it will not be displayed again.

Refer to the [PVE User Management](https://pve.proxmox.com/wiki/User_Management) documentation for more details.

The command outputs a table with the token ID and secret. Concatenate them into a single string (e.g., `user@realm!tokenid=secret`) for the `api_token` field or the `PROXMOX_VE_API_TOKEN` environment variable:

```hcl
provider "proxmox" {
  endpoint  = var.virtual_environment_endpoint
  api_token = "terraform@pve!provider=xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx"
  insecure  = true
  ssh {
    agent    = true
    username = "terraform"
  }
}
```

-> Not all Proxmox API operations are supported via API Token.
You may see errors like
`error creating container: received an HTTP 403 response - Reason: Permission check failed (changing feature flags for privileged container is only allowed for root@pam)` or
`error creating VM: received an HTTP 500 response - Reason: only root can set 'arch' config` or
`Permission check failed (user != root@pam)` when using API Token authentication, even when `Administrator` role or the `root@pam` user is used with the token.
The workaround is to use password authentication for those operations.

-> You can also configure additional Proxmox users and roles using [`virtual_environment_user`](https://registry.terraform.io/providers/bpg/proxmox/latest/docs/resources/virtual_environment_user) and [`virtual_environment_role`](https://registry.terraform.io/providers/bpg/proxmox/latest/docs/resources/virtual_environment_role) resources of the provider.

### Pre-Authentication, or Passing an Authentication Ticket into the provider

It is possible to generate a session ticket with the API, and to pass the ticket and csrf_prevention_token into the provider using environment variables `PROXMOX_VE_AUTH_TICKET` and `PROXMOX_VE_CSRF_PREVENTION_TOKEN` (or provider's arguments `auth_ticket` and `csrf_prevention_token`). See more details in the [Proxmox Wiki](https://pve.proxmox.com/wiki/Proxmox_VE_API#Ticket_Cookie).

An example of using `curl` and `jq` to query the Proxmox API to get a Proxmox session ticket; it is also very easy to pass in a TOTP password this way:

```hcl
provider "proxmox" {
  endpoint = "https://10.0.0.2:8006/"
}
```

```bash
#!/usr/bin/bash

## assume vars are set: PROXMOX_VE_ENDPOINT, PROXMOX_VE_USERNAME, PROXMOX_VE_PASSWORD
## end-goal: automatically set PROXMOX_VE_AUTH_TICKET and PROXMOX_VE_CSRF_PREVENTION_TOKEN

_user_totp_password='123456' ## optional TOTP password


proxmox_api_ticket_path='api2/json/access/ticket' ## cannot have double "//" - ensure endpoint ends with a "/" and this string does not begin with a "/", or vice-versa

## call the auth api endpoint
resp=$( curl -q -s -k --data-urlencode "username=${PROXMOX_VE_USERNAME}"  --data-urlencode "password=${PROXMOX_VE_PASSWORD}"  "${PROXMOX_VE_ENDPOINT}${proxmox_api_ticket_path}" )
auth_ticket=$( jq -r '.data.ticket' <<<"${resp}" )
resp_csrf=$( jq -r '.data.CSRFPreventionToken' <<<"${resp}" )

## check if the response payload needs a TFA (totp) passed, call the auth-api endpoint again
if [[ $(jq -r '.data.NeedTFA' <<<"${resp}") == 1 ]]; then
  resp=$( curl -q -s -k  -H "CSRFPreventionToken: ${resp_csrf}" --data-urlencode  "username=${PROXMOX_VE_USERNAME}" --data-urlencode "tfa-challenge=${auth_ticket}" --data-urlencode "password=totp:${_user_totp_password}"  "${PROXMOX_VE_ENDPOINT}${proxmox_api_ticket_path}" )
  auth_ticket=$( jq -r '.data.ticket' <<<"${resp}" )
  resp_csrf=$( jq -r '.data.CSRFPreventionToken' <<<"${resp}" )
fi


export PROXMOX_VE_AUTH_TICKET="${auth_ticket}"
export PROXMOX_VE_CSRF_PREVENTION_TOKEN="${resp_csrf}"

terraform plan
```

## SSH Connection

-> **SSH is optional for most users.** The provider primarily uses the Proxmox API. SSH is only needed for specific edge cases listed below.

### When is SSH Required?

SSH connection is **only** required for these specific operations:

| Operation | Resource | Why SSH is needed |
| --------- | -------- | ----------------- |
| Upload snippets | `proxmox_virtual_environment_file` | Proxmox API doesn't support snippet uploads |
| Upload certain file types | `proxmox_virtual_environment_file` | Some content types require direct node access |
| Import disks via `source_file.path` | `proxmox_virtual_environment_vm` | Local file transfer to node |
| Configure `idmap` entries | `proxmox_virtual_environment_container` | Proxmox API doesn't support `lxc[n]` parameters |

**SSH is NOT required for:**

- Creating, modifying, or deleting VMs and Containers
- Managing storage, networks, pools, users, or any other resources
- Importing disks using `import_from` attribute (uses API)
- Downloading files using `proxmox_virtual_environment_download_file` (uses API)

If you don't need the operations listed above, you can skip the SSH configuration entirely.

### SSH Configuration

If you need SSH access, the connection is configured via the optional `ssh` block in the `provider` block:

```hcl
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

If no `ssh` block is provided, the provider will attempt to connect to the target node using the credentials provided in the `username` and `password` arguments (or `PROXMOX_VE_USERNAME` and `PROXMOX_VE_PASSWORD` environment variables).
Note that the target node is identified by the `node` argument in the resource, and may be different from the Proxmox API endpoint.
Please refer to the [Argument Reference](#argument-reference) section to view the available arguments of the `ssh` block.

### SSH Agent

The provider does not use OS-specific SSH configuration files, such as `~/.ssh/config`.
Instead, it uses the SSH protocol directly, and supports the `SSH_AUTH_SOCK` environment variable (or `agent_socket` argument) to connect to the SSH agent.
This allows the provider to use the SSH agent configured by the user, and to support multiple SSH agents running on the same machine.
You can find more details on the SSH Agent [here](https://www.digitalocean.com/community/tutorials/ssh-essentials-working-with-ssh-servers-clients-and-keys#adding-your-ssh-keys-to-an-ssh-agent-to-avoid-typing-the-passphrase).
The SSH agent authentication takes precedence over the `private_key` and `password` authentication.

-> By default on Windows, the provider will assume the SSH agent is at `\\.\pipe\openssh-ssh-agent`.

### SSH Private Key

When an SSH agent is not available (for example, in CI/CD pipelines without SSH agent forwarding), you can use the `private_key` argument in the `ssh` block (or the `PROXMOX_VE_SSH_PRIVATE_KEY` environment variable) to provide the private key directly.

The private key must not be encrypted, and must be in PEM format.

You can provide the private key from a file:

```hcl
provider "proxmox" {
  // ...
  ssh {
    agent       = false
    private_key = file("~/.ssh/id_rsa")
  }
}
```

Alternatively, heredoc syntax can supply the private key as a string (not recommended due to security risks). The `<<-` format ignores indentation:

```hcl
provider "proxmox" {
  // ...

  ssh {
    agent       = false
    private_key = <<-EOF
    -----BEGIN OPENSSH PRIVATE KEY-----
    b3BlbnNzaC1rZXktdjEAAAAABG5vbmUAAAAEbm9uZQAAAAAAAAABAAAAMwAAAAtzc2gtZW
    <SKIPPED>
    DMUWUEaH7yMCKl7uCZ9xAAAAAAECAwQF
    -----END OPENSSH PRIVATE KEY-----
    EOF
  }
}
```

### SSH User

By default, the provider will use the same username for the SSH connection as the one used for the Proxmox API connection (when using PAM authentication).
This can be overridden by specifying the `username` argument in the `ssh` block (or alternatively a username in the `PROXMOX_VE_SSH_USERNAME` environment variable):

```hcl
provider "proxmox" {
  // ...

  ssh {
    agent    = true
    username = "terraform"
  }
}
```

-> When using API Token or non-PAM authentication for Proxmox API, the `username` field in the `ssh` block (or alternatively a username in `PROXMOX_VE_USERNAME` or `PROXMOX_VE_SSH_USERNAME` environment variable) is **required**.
This is because the provider needs to know which PAM user to use for the SSH connection.

When using a non-root user for the SSH connection, the user **must** have the `sudo` privilege on the target node without requiring a password.

-> If you run clustered Proxmox VE, you will need to configure the `sudo` privilege for the user on all nodes in the cluster.

-> `sudo` may not be installed by default on Proxmox VE nodes. You can install it via the command line on the Proxmox host: `apt install sudo`

You can configure the `sudo` privilege for the user via the command line on the Proxmox host.
In the example below, we create a user `terraform` and assign the `sudo` privilege to it. Run the following commands on the Proxmox node in the root shell:

- Create a new system user:

    ```sh
    useradd -m terraform
    ```

- Configure the `sudo` privilege for the user, by adding a new sudoers file to the `/etc/sudoers.d` directory:

    ```sh
    visudo -f /etc/sudoers.d/terraform
    ```

  Add the following lines to the file:

    ```text
    terraform ALL=(root) NOPASSWD: /usr/sbin/pvesm
    terraform ALL=(root) NOPASSWD: /usr/sbin/qm
    terraform ALL=(root) NOPASSWD: /usr/bin/tee /var/lib/vz/snippets/[a-zA-Z0-9_][a-zA-Z0-9_.-]*
    ```

  If you use the `idmap` attribute on `proxmox_virtual_environment_container`, the provider edits the container configuration file via SSH. Add the following rules to allow `sed` and `tee` access to the LXC configuration directory:

    ```text
    terraform ALL=(root) NOPASSWD: /usr/bin/sed -i * /etc/pve/lxc/*.conf
    terraform ALL=(root) NOPASSWD: /usr/bin/tee -a /etc/pve/lxc/*.conf
    ```

  If you're using a different datastore for snippets, not the default `local`, you should add the datastore's mount point to the sudoers file as well, for example:

    ```text
    terraform ALL=(root) NOPASSWD: /usr/bin/tee /mnt/pve/cephfs/snippets/[a-zA-Z0-9_][a-zA-Z0-9_.-]*
    ```

  You can find the mount point of the datastore by running `pvesh get /storage/<name>` on the Proxmox node.

  ~> **Security Warning:** Do not use wildcard patterns like `/var/lib/vz/*` in sudoers rules for `tee`. Such patterns allow path traversal attacks (e.g., `/var/lib/vz/../../../etc/sudoers.d/malicious`) that can lead to privilege escalation. Always restrict to specific subdirectories with strict filename patterns as shown above.

- Copy your SSH public key to the `~/.ssh/authorized_keys` file of the `terraform` user on the target node.

- Test the SSH connection and password-less `sudo`:
  
    ```sh
    ssh terraform@<target-node> sudo pvesm apiinfo 
    ```

  You should be able to connect to the target node and see the output containing `APIVER <number>` on the screen without being prompted for your password.

Alternatively, if `pam_ssh_agent_auth` is configured on the target node, the `agent_forwarding` option can forward the SSH agent to the remote server. This allows `sudo` without a password by validating the public SSH key configured for `pam_ssh_agent_auth`.

### Node IP address used for SSH connection

To make the SSH connection, the provider needs to resolve the target node name to an IP address.
The following methods are used to resolve the node name, in the specified order:

1. Enumerate the node's network interfaces via the Proxmox API, and identify the first interface that:
   1. Has an IPv4 address with IPv4 gateway configured, or
   2. Has an IPv6 address with IPv6 gateway configured, or
   3. Has an IPv4 address
2. Resolve the Proxmox node name (usually a shortname) via DNS using the system DNS resolver of the machine running Terraform.

In some cases, this may not be the desired behavior — for example, when the node has multiple network interfaces and the one that should be used for SSH is not the first one.

To override the node IP address used for SSH connection, you can use the optional `node` blocks in the `ssh` block, and specify the desired IP address (or FQDN) for each node.
For example:

```hcl
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

### SSH Connection via SOCKS5 Proxy

The provider supports SSH connection to the target node via a SOCKS5 proxy.

To enable the SOCKS5 proxy, specify the `socks5_server` argument in the `ssh` block:

```hcl
provider "proxmox" {
  // ...
  ssh {
    // ...
    socks5_server     = "ip-or-fqdn-of-socks5-server:port"
    socks5_username   = "username"  # optional  
    socks5_password   = "password"  # optional
  }
}
```

If enabled, this method will be used for all SSH connections to the target nodes in the cluster.

## VM and Container ID Assignment

When creating VMs and Containers, you can specify the optional `vm_id` attribute to set the ID. If omitted, the provider generates a unique ID automatically.

The Proxmox API requires unique IDs within the cluster but doesn't support reserving IDs before resource creation. The provider uses file-based locking to prevent duplicates, but conflicts can still occur when multiple provider instances create resources simultaneously.

To reduce conflicts, set `random_vm_ids = true` in the provider block. This generates random IDs (checked for uniqueness via the API) instead of sequential ones.

## Temporary Directory

Using `proxmox_virtual_environment_file` with `.iso` files or disk images can require a large amount of space in the temporary directory of the computer running Terraform.

Consider pointing `tmp_dir` to a directory with enough space, especially if the default temporary directory is limited by the system memory (e.g. `tmpfs` mounted on `/tmp`).

A better approach is to use the `proxmox_virtual_environment_download_file` resource to download files directly to the target node without buffering to the local machine.

## Environment Variables Summary

All provider arguments can be configured via environment variables. This is the recommended approach for credentials.

**API Connection (required):**

| Environment Variable | Description |
| -------------------- | ----------- |
| `PROXMOX_VE_ENDPOINT` | API endpoint URL (e.g., `https://pve.example.com:8006/`) |

**Authentication (one method required):**

| Environment Variable | Description |
| -------------------- | ----------- |
| `PROXMOX_VE_API_TOKEN` | API token (recommended for production) |
| `PROXMOX_VE_USERNAME` | Username with realm (e.g., `root@pam`) |
| `PROXMOX_VE_PASSWORD` | Password for username/password auth |
| `PROXMOX_VE_AUTH_TICKET` | Pre-authenticated session ticket |
| `PROXMOX_VE_CSRF_PREVENTION_TOKEN` | CSRF token (used with auth ticket) |

**API Options (optional):**

| Environment Variable | Description |
| -------------------- | ----------- |
| `PROXMOX_VE_INSECURE` | Skip TLS verification (`true`/`false`) |
| `PROXMOX_VE_MIN_TLS` | Minimum TLS version (`1.0`, `1.1`, `1.2`, `1.3`) |
| `PROXMOX_VE_TMPDIR` | Custom temporary directory |

**SSH Connection (optional — only if [SSH is required](#when-is-ssh-required)):**

| Environment Variable | Description |
| -------------------- | ----------- |
| `PROXMOX_VE_SSH_USERNAME` | SSH username |
| `PROXMOX_VE_SSH_PASSWORD` | SSH password |
| `PROXMOX_VE_SSH_PRIVATE_KEY` | SSH private key (PEM format) |
| `PROXMOX_VE_SSH_AGENT` | Use SSH agent (`true`/`false`) |
| `PROXMOX_VE_SSH_AUTH_SOCK` | SSH agent socket path |
| `PROXMOX_VE_SSH_AGENT_FORWARDING` | Enable SSH agent forwarding |
| `PROXMOX_VE_SSH_SOCKS5_SERVER` | SOCKS5 proxy server address |
| `PROXMOX_VE_SSH_SOCKS5_USERNAME` | SOCKS5 proxy username |
| `PROXMOX_VE_SSH_SOCKS5_PASSWORD` | SOCKS5 proxy password |

## Argument Reference

In addition to [generic provider arguments](https://developer.hashicorp.com/terraform/language/providers/configuration#provider-configuration-1) (e.g. `alias` and `version`), the following arguments are supported in the Proxmox `provider` block:

- `endpoint` - (Required) The endpoint for the Proxmox Virtual Environment API (can also be sourced from `PROXMOX_VE_ENDPOINT`). Usually this is `https://<your-cluster-endpoint>:8006/`. **Do not** include `/api2/json` at the end.
- `insecure` - (Optional) Whether to skip the TLS verification step (can also be sourced from `PROXMOX_VE_INSECURE`). If omitted, defaults to `false`.
- `min_tls` - (Optional) The minimum required TLS version for API calls (can also be sourced from `PROXMOX_VE_MIN_TLS`). Supported values: `1.0|1.1|1.2|1.3`. If omitted, defaults to `1.3`.

- `auth_ticket` - (Optional) The auth ticket from an external auth call (can also be sourced from `PROXMOX_VE_AUTH_TICKET`). To be used in conjunction with `csrf_prevention_token`, takes precedence over `api_token` and `username` with `password`. For example, `PVE:username@realm:12345678::some_base64_payload==`.
- `csrf_prevention_token` - (Optional) The CSRF Prevention Token from an external auth call (can also be sourced from `PROXMOX_VE_CSRF_PREVENTION_TOKEN`). For example, `12345678:some_blob`.

- `api_token` - (Optional) The API Token for the Proxmox Virtual Environment API (can also be sourced from `PROXMOX_VE_API_TOKEN`). Takes precedence over `username` with `password`. For example, `username@realm!for-terraform-provider=xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx`.

- `otp` - (Optional, Deprecated) The one-time password for the Proxmox Virtual Environment API (can also be sourced from `PROXMOX_VE_OTP`).

- `username` - (Required) The username and realm for the Proxmox Virtual Environment API (can also be sourced from `PROXMOX_VE_USERNAME`). For example, `root@pam`.
- `password` - (Required) The password for the Proxmox Virtual Environment API (can also be sourced from `PROXMOX_VE_PASSWORD`).

- `ssh` - (Optional) The SSH connection configuration to a Proxmox node. This is a block, whose fields are documented below.
    - `username` - (Optional) The username to use for the SSH connection. Defaults to the username used for the Proxmox API connection. Can also be sourced from `PROXMOX_VE_SSH_USERNAME`. Required when using API Token.
    - `password` - (Optional) The password to use for the SSH connection. Defaults to the password used for the Proxmox API connection. Can also be sourced from `PROXMOX_VE_SSH_PASSWORD`.
    - `agent` - (Optional) Whether to use the SSH agent for the SSH authentication. Defaults to `false`. Can also be sourced from `PROXMOX_VE_SSH_AGENT`.
    - `agent_socket` - (Optional) The path to the SSH agent socket. Defaults to the value of the `SSH_AUTH_SOCK` environment variable. Can also be sourced from `PROXMOX_VE_SSH_AUTH_SOCK`.
    - `agent_forwarding` - (Optional) Whether to enable SSH agent forwarding. Defaults to the value of the `PROXMOX_VE_SSH_AGENT_FORWARDING` environment variable, or `false` if not set.
    - `private_key` - (Optional) The private key to use for the SSH connection. Can also be sourced from `PROXMOX_VE_SSH_PRIVATE_KEY`. The private key must be in PEM format.
    - `socks5_server` - (Optional) The address of the SOCKS5 proxy server to use for the SSH connection. Can also be sourced from `PROXMOX_VE_SSH_SOCKS5_SERVER`.
    - `socks5_username` - (Optional) The username to use for the SOCKS5 proxy server. Can also be sourced from `PROXMOX_VE_SSH_SOCKS5_USERNAME`.
    - `socks5_password` - (Optional) The password to use for the SOCKS5 proxy server. Can also be sourced from `PROXMOX_VE_SSH_SOCKS5_PASSWORD`.
    - `node` - (Optional) The node configuration for the SSH connection. Can be specified multiple times to provide configuration for multiple nodes.
        - `name` - (Required) The name of the node.
        - `address` - (Required) The FQDN/IP address of the node.
        - `port` - (Optional) SSH port of the node. Defaults to 22.
- `tmp_dir` - (Optional) Use a custom temporary directory. (can also be sourced from `PROXMOX_VE_TMPDIR`)
- `random_vm_ids` - (Optional) Use random VM IDs for VMs and Containers when `vm_id` attribute is not specified. Defaults to `false`.
- `random_vm_id_start` - (Optional) The start of the range for random VM IDs. Defaults to `10000`.
- `random_vm_id_end` - (Optional) The end of the range for random VM IDs. Defaults to `99999`.
