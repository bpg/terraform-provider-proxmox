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
* `virtual_environment` - (Optional) This is the configuration block for the Proxmox Virtual Environment
    * `endpoint` - (Required) The endpoint for the Proxmox Virtual Environment API
    * `insecure` - (Optional) Whether to skip the TLS verification step (defaults to `false`)
    * `password` - (Required) The password for the Proxmox Virtual Environment API
    * `username` - (Required) The username for the Proxmox Virtual Environment API

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
    * `id` - The member id
    * `node` - The node id
    * `storage` - The storage id
    * `type` - The member type
    * `virtual_machine_id` - The virtual machine id

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
* `expiration_date` - The user account's expiration date
* `first_name` - The user's first name
* `groups` - The user's groups
* `keys` - The user's keys
* `last_name` - The user's last name

##### Users (proxmox_virtual_environment_user)

###### Arguments
This data source doesn't accept arguments.

###### Attributes
* `comments` - The user comments
* `emails` - The users' email addresses
* `enabled` - Whether a user account is enabled
* `expiration_dates` - The user accounts' expiration dates
* `first_names` - The users' first names
* `groups` - The users' groups
* `keys` - The users' keys
* `last_names` - The users' last names
* `user_ids` - The user ids

##### Version (proxmox_virtual_environment_version)

###### Arguments
This data source doesn't accept arguments.

###### Attributes
* `keyboard` - The keyboard layout
* `release` - The release number
* `repository_id` - The repository id
* `version` - The version string

### Resources

#### Virtual Environment

##### File (proxmox_virtual_environment_file)

###### Arguments
* `datastore_id` - (Required) The datastore id
* `node_name` - (Required) The node name
* `override_file_name` - (Optional) The file name to use in the datastore (leave undefined to use source file name)
* `source` - (Required) A path to a local file or a URL
* `template` - (Required) Whether this is a container template (`vztmpl` instead of `iso`)

###### Attributes
* `file_name` - The datastore file name

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
    * `id` - The member id
    * `node` - The node id
    * `storage` - The storage id
    * `type` - The member type
    * `virtual_machine_id` - The virtual machine id

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
* `expiration_date` - (Optional) The user account's expiration date
* `first_name` - (Optional) The user's first name
* `groups` - (Optional) The user's groups
* `keys` - (Optional) The user's keys
* `last_name` - (Optional) The user's last name
* `password` - (Required) The user's password
* `user_id` - (Required) The user id

###### Attributes
This resource doesn't expose any additional attributes.

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
