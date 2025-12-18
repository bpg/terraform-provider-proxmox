# Example 1: Partial management - only manage net0, preserve net1 and net2
resource "proxmox_virtual_environment_cloned_vm" "partial_management" {
  node_name = var.virtual_environment_node_name
  name      = "partial-managed-clone"

  clone = {
    source_vm_id = proxmox_virtual_environment_vm.multi_nic_template.vm_id
    full         = true
  }

  # Only manage the first network interface
  # net1 and net2 are inherited from template but not tracked in Terraform
  network = {
    net0 = {
      bridge = "vmbr0"
      model  = "virtio"
      tag    = 100 # Add VLAN tag to net0
    }
  }

  cpu = {
    cores = 4 # Override CPU count
  }

  agent {
    enabled = true
  }
}

# Example 2: Selective deletion - manage net0, delete net1 and net2
resource "proxmox_virtual_environment_cloned_vm" "selective_deletion" {
  node_name = var.virtual_environment_node_name
  name      = "selective-delete-clone"

  clone = {
    source_vm_id = proxmox_virtual_environment_vm.multi_nic_template.vm_id
    full         = true
  }

  # Manage only net0
  network = {
    net0 = {
      bridge = "vmbr0"
      model  = "virtio"
    }
  }

  # Explicitly delete the inherited net1 and net2
  delete = {
    network = ["net1", "net2"]
  }

  cpu = {
    cores = 2
  }

  agent {
    enabled = true
  }
}

# Example 3: Full management - manage all three NICs with different configs
resource "proxmox_virtual_environment_cloned_vm" "full_management" {
  node_name   = var.virtual_environment_node_name
  name        = "full-managed-clone"
  description = "Clone with all NICs explicitly managed"

  clone = {
    source_vm_id = proxmox_virtual_environment_vm.multi_nic_template.vm_id
    full         = true
  }

  # Explicitly manage all three network interfaces
  network = {
    net0 = {
      bridge   = "vmbr0"
      model    = "virtio"
      tag      = 100 # Management VLAN
      firewall = true
    }

    net1 = {
      bridge   = "vmbr0"
      model    = "virtio"
      tag      = 200 # Application VLAN
      firewall = true
    }

    net2 = {
      bridge     = "vmbr0"
      model      = "virtio"
      tag        = 300 # Storage VLAN
      rate_limit = 100.0
    }
  }

  cpu = {
    cores = 8
  }

  memory = {
    dedicated = 4096
  }

  agent {
    enabled = true
  }
}

# Example 4: Disk management - resize and add disks
resource "proxmox_virtual_environment_cloned_vm" "disk_management" {
  node_name = var.virtual_environment_node_name
  name      = "disk-managed-clone"

  clone = {
    source_vm_id     = proxmox_virtual_environment_vm.multi_nic_template.vm_id
    full             = true
    target_datastore = var.datastore_id
  }

  # Manage disks by slot
  disk = {
    # Resize the boot disk inherited from template
    virtio0 = {
      datastore_id = var.datastore_id
      size_gb      = 50 # Expand from 20GB to 50GB
      discard      = "on"
      iothread     = true
      ssd          = true
    }

    # Add a new data disk
    virtio1 = {
      datastore_id = var.datastore_id
      size_gb      = 100
      backup       = false # Don't include in backups
      cache        = "writethrough"
    }
  }

  network = {
    net0 = {
      bridge = "vmbr0"
      model  = "virtio"
    }
  }

  # Delete unneeded NICs from template
  delete = {
    network = ["net1", "net2"]
  }

  cpu = {
    cores = 4
  }

  agent {
    enabled = true
  }
}

# Outputs to show VM information
output "partial_management_id" {
  value       = proxmox_virtual_environment_cloned_vm.partial_management.id
  description = "VM ID of partially managed clone"
}

output "partial_management_ipv4" {
  value       = try(proxmox_virtual_environment_cloned_vm.partial_management.ipv4_addresses[1][0], "N/A")
  description = "IPv4 address of partially managed clone"
}

output "selective_deletion_id" {
  value       = proxmox_virtual_environment_cloned_vm.selective_deletion.id
  description = "VM ID of selective deletion clone"
}

output "full_management_id" {
  value       = proxmox_virtual_environment_cloned_vm.full_management.id
  description = "VM ID of fully managed clone"
}

output "disk_management_id" {
  value       = proxmox_virtual_environment_cloned_vm.disk_management.id
  description = "VM ID of disk managed clone"
}
