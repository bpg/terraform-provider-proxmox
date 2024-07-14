---
layout: home
title: "Provider: Proxmox Virtual Environment"
---

# Proxmox Provider

This provider for [Terraform](https://www.terraform.io/) is used for interacting with resources supported by [Proxmox](https://www.proxmox.com/en/).
The provider needs to be configured with the proper endpoints and credentials before it can be used.

Use the navigation to the left to read about the available resources.

## Example Usage

```hcl
provider "proxmox" {
  endpoint = "https://10.0.0.2:8006/"

  # TODO: use terraform variable or remove the line, and use PROXMOX_VE_AUTH_PAYLOAD environment variable
  # auth_payload = '{"data":{"CSRFPreventionToken":"12345678:some_blob","username":"username@realm","ticket":"PVE:username@realmy:12345678::some_base64_payload==","cap":{"access":{"User.Modify":1}}}}'

  # TODO: use terraform variable or remove the line, and use PROXMOX_VE_AUTH_TICKET environment variable
  # auth_ticket           = "PVE:username@realm:12345678::some_base64_payload=="
  # TODO: use terraform variable or remove the line, and use PROXMOX_VE_CSRF_PREVENTION_TOKEN= environment variable
  # csrf_prevention_token = "12345678:some_blob"

  # TODO: use terraform variable or remove the line, and use PROXMOX_VE_API_TOKEN environment variable
  # api_token = "root@pam!for-terraform-provider=xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx"

  # TODO: use terraform variable or remove the line, and use PROXMOX_VE_USERNAME environment variable
  username = "root@pam"
  # TODO: use terraform variable or remove the line, and use PROXMOX_VE_PASSWORD environment variable
  password = "the-password-set-during-installation-of-proxmox-ve"

  # because self-signed TLS certificate is in use
  insecure = true
  # uncomment (unless on Windows...)
  # tmp_dir  = "/var/tmp"

  ssh {
    agent = true
    # TODO: uncomment and configure if using api_token instead of password
    # username = "root"
  }
}
```

## Authentication

The Proxmox provider offers a flexible means of providing credentials for authentication.
Static credentials and pre-authenticated session-ticket can be provided to the `proxmox` block through one the choices of arguments below, ordered by precedence:
- `auth_payload`
- `auth_ticket` and `csrf_prevention_token`
- `api_token`
- `username` and `password`

!> Hard-coding credentials into any Terraform configuration is not recommended, and risks secret leakage should this file ever be committed to a public version control system.

Static credentials can be provided in-line in the Proxmox provider block, by adding one of the arguments above:

```hcl
provider "proxmox" {
  endpoint = "https://10.0.0.2:8006/"

  auth_payload = "{\"data\":{\"CSRFPreventionToken\":\"12345678:some_blob\",\"username\":\"username@realm\",\"ticket\":\"PVE:username@realmn:12345678::some_base64_payload==\",\"cap\":{\"access\":{\"User.Modify\":1}}}}"
  # setting the auth_payload into terraform will need (at least) its double-quotes escaped; even better to use as an environment-variable

  auth_ticket = "PVE:username@realm:12345678::some_base64_payload=="
  csrf_prevention_token = "12345678:some_blob"

  api_token = "username@realm!for-terraform-provider=xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx"

  username = "username@realm"
  password = "a-strong-password"
}
```

A better approach is to extract these values into Terraform variables, and reference the variables instead:

```hcl
provider "proxmox" {
  endpoint = var.virtual_environment_endpoint

  auth_payload = var.virtual_environment_auth_payload

  auth_ticket           = var.virtual_environment_auth_ticket
  csrf_prevention_token = var.virtual_environment_csrf_prevention_token

  api_token = var.virtual_environment_api_token

  username = var.virtual_environment_username
  password = var.virtual_environment_password
}
```

The variable values can be provided via a separate `.tfvars` file that should be gitignored.
See the [Terraform documentation](https://www.terraform.io/docs/configuration/variables.html) for more information.


### Environment variables

Instead of using static arguments, credentials can be handled through the use of environment variables.
For example:

```hcl
provider "proxmox" {
  endpoint = "https://10.0.0.2:8006/"
}
```

```sh
export PROXMOX_VE_AUTH_PAYLOAD='{"data":{"CSRFPreventionToken":"12345678:some_blob","username":"username@realm","ticket":"PVE:username@realm:12345678::some_base64_payload==","cap":{"access":{"User.Modify":1}}}}'

export PROXMOX_VE_AUTH_TICKET='PVE:username@realm:12345678::some_base64_payload=='
export PROXMOX_VE_CSRF_PREVENTION_TOKEN='12345678:some_blob'

export PROXMOX_VE_API_TOKEN='username@realm!for-terraform-provider=xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx'

export PROXMOX_VE_USERNAME="username@realm"
export PROXMOX_VE_PASSWORD='a-strong-password'
terraform plan
```

See the [Argument Reference](#argument-reference) section for the supported variable names and use cases.

## Pre-Authentication, or Passing an Authentication Ticket into the provider

It is possible to generate a session ticket with the api, and to pass the ticket-json into the provider with 1 of 2 methods:
  - env-var `PROXMOX_VE_AUTH_PAYLOAD` (or auth_payload)
  - env-vars `PROXMOX_VE_AUTH_TICKET` with `PROXMOX_VE_CSRF_PREVENTION_TOKEN` (auth_ticket with csrf_prevention_token)

An example `bash` of using `curl` and `jq` to query the proxmox api to get a proxmox session ticket; it is also very easy to pass in a TOTP password this way
- Note: the export'ing of the 3x vars are redundant and shown for demonstration purposes, pick either: `PROXMOX_VE_AUTH_PAYLOAD` or `PROXMOX_VE_AUTH_TICKET` with `PROXMOX_VE_CSRF_PREVENTION_TOKEN`

```bash
#!/usr/bin/bash

## assume vars are set: PROXMOX_VE_ENDPOINT, PROXMOX_VE_USERNAME, PROXMOX_VE_PASSWORD
## end-goal: automatically set PROXMOX_VE_AUTH_PAYLOAD

_user_totp_password='123456' ## optional TOTP password

proxmox_api_ticket_path='api2/json/access/ticket' ## cannot have double "//" - ensure endpoint ends with a "/" and this string does not begin with a "/", or vice-versa


resp=$( curl -q -s -k --data-urlencode "username=${PROXMOX_VE_USERNAME}"  --data-urlencode "password=${PROXMOX_VE_PASSWORD}"  "${PROXMOX_VE_ENDPOINT}${proxmox_api_ticket_path}" )
auth_ticket=$( jq -r '.data.ticket' <<<"${resp}" )
resp_csrf=$( jq -r '.data.CSRFPreventionToken' <<<"${resp}" )


if [[ $(jq -r '.data.NeedTFA' <<<"${resp}") == 1 ]]; then
  resp=$( curl -q -s -k  -H "CSRFPreventionToken: ${resp_csrf}" --data-urlencode  "username=${PROXMOX_VE_USERNAME}" --data-urlencode "tfa-challenge=${auth_ticket}" --data-urlencode "password=totp:${_user_totp_password}"  "${PROXMOX_VE_ENDPOINT}${proxmox_api_ticket_path}" )
  auth_ticket=$( jq -r '.data.ticket' <<<"${resp}" )
  resp_csrf=$( jq -r '.data.CSRFPreventionToken' <<<"${resp}" )
fi



## PROXMOX_VE_AUTH_PAYLOAD
export PROXMOX_VE_AUTH_PAYLOAD="${resp}"


## PROXMOX_VE_AUTH_TICKET  with  PROXMOX_VE_CSRF_PREVENTION_TOKEN
export PROXMOX_VE_AUTH_TICKET="${auth_ticket}"
export PROXMOX_VE_CSRF_PREVENTION_TOKEN="${resp_csrf}"

```

## SSH Connection

~> Please read if you are using VMs with custom disk images, or uploading snippets.

The Proxmox provider can connect to a Proxmox node via SSH.
This is used in the `proxmox_virtual_environment_vm` or `proxmox_virtual_environment_file` resource to execute commands on the node to perform actions that are not supported by Proxmox API.
For example, to import VM disks, or to uploading certain type of resources, such as snippets.

The SSH connection configuration is provided via the optional `ssh` block in the `provider` block:

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
Instead, it uses the SSH protocol directly, and supports the `SSH_AUTH_SOCK` environment variable (or `agent_socket` argument) to connect to the `ssh-agent`.
This allows the provider to use the SSH agent configured by the user, and to support multiple SSH agents running on the same machine.
You can find more details on the SSH Agent [here](https://www.digitalocean.com/community/tutorials/ssh-essentials-working-with-ssh-servers-clients-and-keys#adding-your-ssh-keys-to-an-ssh-agent-to-avoid-typing-the-passphrase).
The SSH agent authentication takes precedence over the `private_key` and `password` authentication.

-> By default on Windows, the provider will assume the SSH agent is at `\\.\pipe\openssh-ssh-agent`.

### SSH Private Key

In some cases where SSH agent is not available, for example when using a CI/CD pipeline that does not support SSH agent forwarding, 
you can use the `private_key` argument in the `ssh` block (or alternatively `PROXMOX_VE_SSH_PRIVATE_KEY` environment variable) to provide the private key for the SSH connection.

The private key mut not be encrypted, and must be in PEM format.

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

Alternatively, although not recommended due to the increased risk of exposing an unprotected key, heredoc syntax can be used to supply the private key as a string.
Note that the content of the private key is injected using `<<-` format to ignore indentation:
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

~> The `root` user on the Proxmox node must be configured with `bash` as the default shell.

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
    terraform ALL=(root) NOPASSWD: /sbin/pvesm
    terraform ALL=(root) NOPASSWD: /sbin/qm
    terraform ALL=(root) NOPASSWD: /usr/bin/tee /var/lib/vz/*
    ```

  If you're using a different datastore for snippets, not the default `local`, you should add the datastore's mount point to the sudoers file as well, for example: 
  
    ```text
    terraform ALL=(root) NOPASSWD: /usr/bin/tee /mnt/pve/cephfs/*
    ```
  You can find the mount point of the datastore by running `pvesh get /storage/<name>` on the Proxmox node.

- Copy your SSH public key to the `~/.ssh/authorized_keys` file of the `terraform` user on the target node.

- Test the SSH connection and password-less `sudo`:
  
    ```sh
    ssh terraform@<target-node> sudo pvesm apiinfo 
    ```

  You should be able to connect to the target node and see the output containing `APIVER <number>` on the screen without being prompted for your password.

### Node IP address used for SSH connection

In order to make the SSH connection, the provider needs to be able to resolve the target node name to an IP.
The following methods are used to resolve the node name, in the specified order:

1. Enumerate the node's network interfaces via the Proxmox API, and identify the first interface that:
   1. Has an IPv4 address with IPv4 gateway configured, or
   2. Has an IPv6 address with IPv6 gateway configured, or
   3. Has an IPv4 address
2. Resolve the Proxmox node name (usually a shortname) via DNS using the system DNS resolver of the machine running Terraform.

In some cases this may not be the desired behavior, for example, when the node has multiple network interfaces, and the one that should be used for SSH is not the first one.

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

To enable the SOCKS5 proxy, you need to configure the `ssh` block in the `provider` block, and specify the `socks5_server` argument:

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

## API Token Authentication

API Token authentication can be used to authenticate with the Proxmox API without the need to provide a password.
In combination with the `ssh` block and `ssh-agent` support, this allows for a fully password-less authentication.

You can create an API Token for a user via the Proxmox UI, or via the command line on the Proxmox host or cluster:

- Create a user:

    ```sh
    sudo pveum user add terraform@pve
    ```

- Create a role for the user (you can skip this step if you want to use any of the existing roles):

    ```sh
    sudo pveum role add Terraform -privs "Datastore.Allocate Datastore.AllocateSpace Datastore.AllocateTemplate Datastore.Audit Pool.Allocate Sys.Audit Sys.Console Sys.Modify SDN.Use VM.Allocate VM.Audit VM.Clone VM.Config.CDROM VM.Config.Cloudinit VM.Config.CPU VM.Config.Disk VM.Config.HWType VM.Config.Memory VM.Config.Network VM.Config.Options VM.Migrate VM.Monitor VM.PowerMgmt User.Modify"
    ```

  ~> The list of privileges above is only an example, please review it and adjust to your needs.
  Refer to the [privileges documentation](https://pve.proxmox.com/pve-docs/pveum.1.html#_privileges) for more details.

- Assign the role to the previously created user:

    ```sh
    sudo pveum aclmod / -user terraform@pve -role Terraform
    ```

- Create an API token for the user:

    ```sh
    sudo pveum user token add terraform@pve provider --privsep=0
    ```

Refer to the upstream docs as needed for additional details concerning [PVE User Management](https://pve.proxmox.com/wiki/User_Management).

Generating the token will output a table containing the token's ID and secret which are meant to be concatenated into a single string for use with either the `api_token` field of the `provider` block (fine for testing but should be avoided) or sourced from the `PROXMOX_VE_API_TOKEN` environment variable.

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

-> The token authentication is taking precedence over the password authentication.

-> Not all Proxmox API operations are supported via API Token.
You may see errors like 
`error creating container: received an HTTP 403 response - Reason: Permission check failed (changing feature flags for privileged container is only allowed for root@pam)` or 
`error creating VM: received an HTTP 500 response - Reason: only root can set 'arch' config` or
`Permission check failed (user != root@pam)` when using API Token authentication, even when `Administrator` role or the `root@pam` user is used with the token.
The workaround is to use password authentication for those operations.

-> You can also configure additional Proxmox users and roles using [`virtual_environment_user`](https://registry.terraform.io/providers/bpg/proxmox/latest/docs/data-sources/virtual_environment_user) and [`virtual_environment_role`](https://registry.terraform.io/providers/bpg/proxmox/latest/docs/data-sources/virtual_environment_role) resources of the provider.

## Temporary Directory

Using `proxmox_virtual_environment_file` with `.iso` files or disk images can require large amount of space in the temporary directory of the computer running terraform.

Consider pointing `tmp_dir` to a directory with enough space, especially if the default temporary directory is limited by the system memory (e.g. `tmpfs` mounted on `/tmp`).

## Argument Reference

In addition to [generic provider arguments](https://www.terraform.io/docs/configuration/providers.html) ( e.g. `alias` and `version`), the following arguments are supported in the Proxmox `provider` block:

- `endpoint` - (Required) The endpoint for the Proxmox Virtual Environment API (can also be sourced from `PROXMOX_VE_ENDPOINT`). Usually this is `https://<your-cluster-endpoint>:8006/`. **Do not** include `/api2/json` at the end.
- `insecure` - (Optional) Whether to skip the TLS verification step (can also be sourced from `PROXMOX_VE_INSECURE`). If omitted, defaults to `false`.
- `min_tls` - (Optional) The minimum required TLS version for API calls (can also be sourced from `PROXMOX_VE_MIN_TLS`). Supported values: `1.0|1.1|1.2|1.3`. If omitted, defaults to `1.3`.

- `auth_payload` - (Optional) The full authentication ticket json from an external auth call (can also be sourced from `PROXMOX_VE_AUTH_PAYLOAD`). Takes precedence over other api auth-methods. For example, `{"data":{"CSRFPreventionToken":"12345678:some_blob","username":"username@realm","ticket":"PVE:username@realm:12345678::some_base64_payload==","cap":{"access":{"User.Modify":1}}}}`.

- `auth_ticket` - (Optional) The auth ticket from an external auth call (can also be sourced from `PROXMOX_VE_AUTH_TICKET`). To be used in conjunction with `csrf_prevention_token`, takes precedence over `api_token` and `username` with `password`. For example, `PVE:username@realm:12345678::some_base64_payload==`.
- `csrf_prevention_token` - (Optional) The CSRF Prevention Token from an external auth call (can also be sourced from `PROXMOX_VE_CSRF_PREVENTION_TOKEN`). For example, `12345678:some_blob`.

- `api_token` - (Optional) The API Token for the Proxmox Virtual Environment API (can also be sourced from `PROXMOX_VE_API_TOKEN`). Takes precedence over `username` with `password. For example, `username@realm!for-terraform-provider=xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx`.

- `otp` - (Optional, Deprecated) The one-time password for the Proxmox Virtual Environment API (can also be sourced from `PROXMOX_VE_OTP`).

- `username` - (Required) The username and realm for the Proxmox Virtual Environment API (can also be sourced from `PROXMOX_VE_USERNAME`). For example, `root@pam`.
- `password` - (Required) The password for the Proxmox Virtual Environment API (can also be sourced from `PROXMOX_VE_PASSWORD`).

- `ssh` - (Optional) The SSH connection configuration to a Proxmox node. This is a block, whose fields are documented below.
    - `username` - (Optional) The username to use for the SSH connection. Defaults to the username used for the Proxmox API connection. Can also be sourced from `PROXMOX_VE_SSH_USERNAME`. Required when using API Token.
    - `password` - (Optional) The password to use for the SSH connection. Defaults to the password used for the Proxmox API connection. Can also be sourced from `PROXMOX_VE_SSH_PASSWORD`.
    - `agent` - (Optional) Whether to use the SSH agent for the SSH authentication. Defaults to `false`. Can also be sourced from `PROXMOX_VE_SSH_AGENT`.
    - `agent_socket` - (Optional) The path to the SSH agent socket. Defaults to the value of the `SSH_AUTH_SOCK` environment variable. Can also be sourced from `PROXMOX_VE_SSH_AUTH_SOCK`.
    - `private_key` - (Optional) The private key to use for the SSH connection. Can also be sourced from `PROXMOX_VE_SSH_PRIVATE_KEY`. The private key must be in PEM format.
    - `socks5_server` - (Optional) The address of the SOCKS5 proxy server to use for the SSH connection. Can also be sourced from `PROXMOX_VE_SSH_SOCKS5_SERVER`.
    - `socks5_username` - (Optional) The username to use for the SOCKS5 proxy server. Can also be sourced from `PROXMOX_VE_SSH_SOCKS5_USERNAME`.
    - `socks5_password` - (Optional) The password to use for the SOCKS5 proxy server. Can also be sourced from `PROXMOX_VE_SSH_SOCKS5_PASSWORD`.
    - `node` - (Optional) The node configuration for the SSH connection. Can be specified multiple times to provide configuration fo multiple nodes.
        - `name` - (Required) The name of the node.
        - `address` - (Required) The FQDN/IP address of the node.
        - `port` - (Optional) SSH port of the node. Defaults to 22.
- `tmp_dir` - (Optional) Use custom temporary directory. (can also be sourced from `PROXMOX_VE_TMPDIR`)
