# Example 1: Partial management - only manage net0
# NOTE: The template only has one NIC (net0), so this example manages it explicitly
resource "proxmox_virtual_environment_cloned_vm" "partial_management" {
  node_name = var.virtual_environment_node_name
  name      = "partial-managed-clone"

  clone = {
    source_vm_id = proxmox_virtual_environment_vm.ubuntu_template.vm_id
    full         = true
  }

  # Only manage the first network interface
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

  # NOTE: Agent settings are inherited from template
}

# Example 2: Memory configuration with new terminology
resource "proxmox_virtual_environment_cloned_vm" "selective_deletion" {
  node_name = var.virtual_environment_node_name
  name      = "memory-configured-clone"

  clone = {
    source_vm_id = proxmox_virtual_environment_vm.ubuntu_template.vm_id
    full         = true
  }

  # Manage network interface
  network = {
    net0 = {
      bridge = "vmbr0"
      model  = "virtio"
    }
  }

  # Memory configuration using new terminology
  memory = {
    size    = 3072 # Total memory available
    balloon = 1024 # Minimum guaranteed memory
  }

  cpu = {
    cores = 2
  }
}

# Example 3: Full management - manage network, CPU, and memory
resource "proxmox_virtual_environment_cloned_vm" "full_management" {
  node_name   = var.virtual_environment_node_name
  name        = "full-managed-clone"
  description = "Clone with all configuration explicitly managed"

  clone = {
    source_vm_id = proxmox_virtual_environment_vm.ubuntu_template.vm_id
    full         = true
  }

  # Manage network interface
  network = {
    net0 = {
      bridge   = "vmbr0"
      model    = "virtio"
      tag      = 100 # Management VLAN
      firewall = true
    }
  }

  cpu = {
    cores = 8
  }

  # Memory configuration using new terminology
  memory = {
    size    = 4096 # Total memory available
    balloon = 2048 # Minimum guaranteed memory
  }
}

# Example 4: Disk management - resize and add disks
resource "proxmox_virtual_environment_cloned_vm" "disk_management" {
  node_name = var.virtual_environment_node_name
  name      = "disk-managed-clone"

  clone = {
    source_vm_id     = proxmox_virtual_environment_vm.ubuntu_template.vm_id
    full             = true
    target_datastore = var.datastore_id
  }

  # Manage disks by slot
  disk = {
    # Resize the boot disk inherited from template
    virtio0 = {
      datastore_id = var.datastore_id
      disk_size    = "50G" # Expand from 20GB to 50GB
      discard      = "on"
      iothread     = true
      ssd          = true
    }

    # Add a new data disk
    virtio1 = {
      datastore_id = var.datastore_id
      disk_size    = "100G"
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

  cpu = {
    cores = 4
  }

  memory = {
    size    = 2048
    balloon = 512
  }
}

# Outputs to show VM information
output "partial_management_id" {
  value       = proxmox_virtual_environment_cloned_vm.partial_management.id
  description = "VM ID of partially managed clone"
}

# NOTE: ipv4_addresses is not available on cloned_vm resource
# Use proxmox_virtual_environment_vm datasource if you need IP addresses

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
