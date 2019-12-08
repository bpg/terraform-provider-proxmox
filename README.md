# Terraform Provider for Proxmox
A Terraform Provider which adds support for Proxmox solutions.

## Requirements
- [Terraform](https://www.terraform.io/downloads.html) 0.11+
- [Go](https://golang.org/doc/install) 1.12 (to build the provider plugin)

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
* `virtual_environment` - (Optional) This is the configuration block for the Proxmox Virtual Environment.
    * `endpoint` - (Required) The endpoint for the Proxmox Virtual Environment API.
    * `insecure` - (Optional) Whether to skip the TLS verification step. Defaults to `false`.
    * `password` - (Required) The password for the Proxmox Virtual Environment API.
    * `username` - (Required) The username for the Proxmox Virtual Environment API.

### Data Sources

#### Virtual Environment

##### Access Group (proxmox_virtual_environment_access_group)

###### Arguments
* `group_id` - (Required) The group id.

###### Attributes
* `comment` - The group comment.
* `members` - The group members as a list with `username@realm` entries.

##### Access Groups (proxmox_virtual_environment_access_groups)

###### Arguments
This data source doesn't accept arguments.

###### Attributes
* `comments` - The group comments.
* `group_ids` - The group ids.

##### Access Role (proxmox_virtual_environment_access_role)

###### Arguments
* `role_id` - (Required) The role id.

###### Attributes
* `privileges` - The role privileges.

##### Access Roles (proxmox_virtual_environment_access_roles)

###### Arguments
This data source doesn't accept arguments.

###### Attributes
* `privileges` - The role privileges.
* `role_ids` - The role ids.
* `special` - Whether the role is special (built-in).

##### Version (proxmox_virtual_environment_version)

###### Arguments
This data source doesn't accept arguments.

###### Attributes
* `keyboard` - The keyboard layout.
* `release` - The release number.
* `repository_id` - The repository id.
* `version` - The version string.

### Resources

#### Virtual Environment

##### Access Group (proxmox_virtual_environment_access_group)

###### Arguments
* `comment` - (Optional) The group comment.
* `group_id` - (Required) The group id.

###### Attributes
* `members` - The group members as a list with `username@realm` entries.

## Developing the Provider
If you wish to work on the provider, you'll first need [Go](http://www.golang.org) installed on your machine (version 1.12+ is *required*). You'll also need to correctly setup a [GOPATH](http://golang.org/doc/code.html#GOPATH), as well as adding `$GOPATH/bin` to your `$PATH`.

To compile the provider, run `make build`. This will build the provider and put the provider binary in the `$GOPATH/bin` directory.

```sh
$ make build
...
$ $GOPATH/bin/terraform-provider-proxmox
...
```

If you wish to contribute to the provider, the following requirements must be met,

* All tests must pass using `make test`
* The Go code must be formatted using Gofmt
* Dependencies are installed by `make init`

## Testing the Provider
In order to test the provider, you can simply run `make test`.

```sh
$ make test
```

Tests are limited to regression tests, ensuring backwards compability.
