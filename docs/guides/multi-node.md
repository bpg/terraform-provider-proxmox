---
layout: page
page_title: "Multi-Node Cluster Management"
subcategory: Guides
description: |-
    Managing resources across multiple nodes in a Proxmox VE cluster.
---

# Multi-Node Cluster Management

This guide covers working with multiple nodes within a Proxmox VE cluster — SSH configuration, cross-node cloning, VM migration, HA drift handling, and node-scoped resources.

## SSH node configuration

The provider uses SSH for file uploads (cloud-init snippets, ISOs) and certain node-level commands. By default, it resolves node SSH addresses from the Proxmox API by querying each node's network interfaces. Use `ssh { node { } }` blocks when the default resolution does not work — for example, when nodes are behind NAT, use non-standard SSH ports, or have hostnames that don't resolve from the Terraform host:

```terraform
provider "proxmox" {
  endpoint  = "https://pve.example.com:8006/"
  api_token = var.api_token

  ssh {
    agent    = true
    username = "terraform"

    node {
      name    = "pve1"
      address = "10.0.1.10"
    }

    node {
      name    = "pve2"
      address = "10.0.1.11"
      port    = 2222
    }
  }
}
```

Each `node` block maps a Proxmox node name (as shown in the PVE web UI) to its SSH address and optional port.

-> **Note:** SSH is needed for snippet/file uploads, VM disk imports (when using `file_id` on a disk), and container `idmap` entries. Most other operations (VM/container CRUD, migration, downloads) use the HTTP API and do not require SSH.

## Cross-node cloning

To clone a VM template that lives on a different node than the target, use `clone { node_name }` to specify the source node:

```terraform
resource "proxmox_virtual_environment_vm" "from_remote_template" {
  name      = "app-server"
  node_name = "pve2"

  clone {
    vm_id     = 100
    node_name = "pve1"
  }
}
```

The provider handles the two cases automatically:

- **Shared storage** — clones directly to the target node in one API call.
- **Local (non-shared) storage** — clones on the source node first, then migrates the new VM to the target node. A `clone { datastore_id }` can be set to specify the target datastore for the migration.

~> **Important:** Cross-node cloning only works within the same Proxmox cluster. The Proxmox API does not support cloning across cluster boundaries.

## VM migration

By default, changing a VM's `node_name` forces Terraform to destroy and recreate it. Set `migrate = true` to migrate the VM in-place instead:

```terraform
resource "proxmox_virtual_environment_vm" "example" {
  name      = "app-server"
  node_name = "pve2"
  migrate   = true

  # ...
}
```

Changing `node_name` from `"pve1"` to `"pve2"` now migrates the VM in-place instead of a destroy/create cycle.

The provider handles HA-managed VMs automatically:

- **Running HA VM** — uses the HA migrate endpoint, which sequences the migration correctly.
- **Stopped HA VM** — temporarily removes the VM from HA, migrates it, then re-adds it with the original HA configuration.
- **Non-HA VM** — uses the standard migration API with local disk transfer.

The migration timeout defaults to 1800 seconds (30 minutes) and can be adjusted with `timeout_migrate`.

## HA clusters and `node_name` drift

When Proxmox HA moves a VM to another node (due to node failure, maintenance, or rebalancing), Terraform detects a drift on `node_name` and wants to move it back or recreate it. The provider does not perform automatic load balancing or node selection — Proxmox VE's API requires a single target node, and the provider follows the API.

If you use HA and want Terraform to tolerate VMs being on any node, use `lifecycle`:

```terraform
resource "proxmox_virtual_environment_vm" "ha_managed" {
  name      = "ha-server"
  node_name = "pve1"

  lifecycle {
    ignore_changes = [node_name]
  }
}

resource "proxmox_virtual_environment_haresource" "ha_managed" {
  resource_id = "vm:${proxmox_virtual_environment_vm.ha_managed.vm_id}"
  group       = "my-ha-group"
  state       = "started"
}
```

With `ignore_changes = [node_name]`, Terraform will not attempt to move the VM back to `pve1` if HA migrated it elsewhere. Updates to other VM attributes still work — the Proxmox API accepts configuration changes for a VM from any node in the cluster.

~> **Caveat:** With `ignore_changes`, you lose the ability to deliberately move a VM by changing `node_name` in Terraform. If you need both HA tolerance and Terraform-driven migration, manage them as separate operational concerns.

## Node-scoped resources

Many resources are tied to a specific node. When deploying across multiple nodes, you need one resource instance per node:

```terraform
locals {
  nodes = toset(["pve1", "pve2", "pve3"])
}

# Download the cloud image to every node
resource "proxmox_virtual_environment_download_file" "ubuntu_image" {
  for_each     = local.nodes
  node_name    = each.key
  content_type = "iso"
  datastore_id = "local"
  url          = "https://cloud-images.ubuntu.com/jammy/current/jammy-server-cloudimg-amd64.img"
}
```

Node-scoped resources include: `download_file`, `file`, `oci_image`, `certificate`, `dns`, `hosts`, `time`, `network_linux_bridge`, and `network_linux_vlan`. The `apt_repository` and `apt_standard_repository` resources are also node-scoped but use `node` (not `node_name`) as their attribute name.

-> **Tip:** If you use shared storage (NFS, Ceph, GlusterFS), a file downloaded to one node is accessible from all nodes. You only need one `download_file` resource per shared datastore, not per node.

## Dynamic node discovery

Use the `proxmox_virtual_environment_nodes` data source to discover nodes dynamically instead of hardcoding names:

```terraform
data "proxmox_virtual_environment_nodes" "all" {}

output "node_names" {
  value = data.proxmox_virtual_environment_nodes.all.names
}
```

For selecting a node for new VMs, combine with `random_shuffle` from the [`hashicorp/random` provider](https://registry.terraform.io/providers/hashicorp/random/latest):

```terraform
data "proxmox_virtual_environment_nodes" "all" {}

locals {
  online_nodes = [
    for i, name in data.proxmox_virtual_environment_nodes.all.names :
    name if data.proxmox_virtual_environment_nodes.all.online[i]
  ]
}

resource "random_shuffle" "node" {
  input        = local.online_nodes
  result_count = 1
}

resource "proxmox_virtual_environment_vm" "example" {
  name      = "app-server"
  node_name = random_shuffle.node.result[0]

  lifecycle {
    ignore_changes = [node_name]
  }

  # ...
}
```

~> **Important:** The provider does not perform automatic node selection or load balancing. The Proxmox API requires a specific node name for VM creation. Node selection logic belongs in your Terraform configuration (as shown above), not in the provider.

## SOCKS5 proxy

For nodes reachable only through a bastion host, configure a SOCKS5 proxy for the SSH connection:

```terraform
provider "proxmox" {
  endpoint = "https://pve.example.com:8006/"

  ssh {
    agent          = true
    username       = "terraform"
    socks5_server  = "bastion.example.com:1080"
    socks5_username = var.socks5_user
    socks5_password = var.socks5_pass
  }
}
```

The SOCKS5 proxy applies only to SSH connections (file uploads, node commands). API calls go directly to the `endpoint` URL.
