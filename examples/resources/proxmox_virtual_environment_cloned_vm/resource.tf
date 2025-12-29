# Example 1: Basic clone with minimal management
resource "proxmox_virtual_environment_cloned_vm" "basic_clone" {
  node_name = "pve"
  name      = "basic-clone"

  clone = {
    source_vm_id = 100  # Template VM ID
    full         = true # Perform full clone (not linked)
  }

  # Only manage CPU, inherit everything else from template
  cpu = {
    cores = 4
  }
}

# Example 2: Clone with explicit network management
resource "proxmox_virtual_environment_cloned_vm" "network_managed" {
  node_name = "pve"
  name      = "network-clone"

  clone = {
    source_vm_id = 100
  }

  # Map-based network devices - manage specific interfaces
  network = {
    net0 = {
      bridge = "vmbr0"
      model  = "virtio"
      tag    = 100 # VLAN tag
    }

    net1 = {
      bridge      = "vmbr1"
      model       = "virtio"
      firewall    = true
      mac_address = "BC:24:11:2E:C5:00"
    }
  }

  cpu = {
    cores = 2
  }
}

# Example 3: Clone with disk management
resource "proxmox_virtual_environment_cloned_vm" "disk_managed" {
  node_name = "pve"
  name      = "disk-clone"

  clone = {
    source_vm_id     = 100
    target_datastore = "local-lvm"
  }

  # Map-based disk management
  disk = {
    scsi0 = {
      # Resize the cloned boot disk
      datastore_id = "local-lvm"
      size_gb      = 50
      discard      = "on"
      ssd          = true
    }

    scsi1 = {
      # Add a new data disk
      datastore_id = "local-lvm"
      size_gb      = 100
      backup       = false
    }
  }
}

# Example 4: Clone with explicit device deletion
resource "proxmox_virtual_environment_cloned_vm" "selective_delete" {
  node_name = "pve"
  name      = "minimal-clone"

  clone = {
    source_vm_id = 100
  }

  # Only manage net0
  network = {
    net0 = {
      bridge = "vmbr0"
      model  = "virtio"
    }
  }

  # Explicitly delete inherited net1 and net2
  delete = {
    network = ["net1", "net2"]
  }
}

# Example 5: Full-featured clone with multiple settings
resource "proxmox_virtual_environment_cloned_vm" "full_featured" {
  node_name   = "pve"
  name        = "production-vm"
  description = "Production VM cloned from template"
  tags        = ["production", "web"]

  clone = {
    source_vm_id     = 100
    source_node_name = "pve" # Source node (defaults to target node if omitted)
    full             = true
    target_datastore = "local-lvm"
    retries          = 3
  }

  cpu = {
    cores        = 8
    sockets      = 1
    architecture = "x86_64"
    type         = "host"
  }

  memory = {
    size    = 8192
    balloon = 2048
    shares  = 2000
  }

  network = {
    net0 = {
      bridge     = "vmbr0"
      model      = "virtio"
      tag        = 100
      firewall   = true
      rate_limit = 100.0 # MB/s
    }
  }

  disk = {
    scsi0 = {
      datastore_id = "local-lvm"
      size_gb      = 100
      discard      = "on"
      iothread     = true
      ssd          = true
      cache        = "writethrough"
    }
  }

  vga = {
    type   = "std"
    memory = 16
  }

  # Option 1: Remove inherited CD-ROM using delete block
  delete = {
    disk = ["ide2"] # Remove inherited CD-ROM device
  }

  # Option 2: Alternatively, manage the CD-ROM and ensure it's empty
  # cdrom = {
  #   ide2 = {
  #     file_id = "none" # Ensure CD-ROM is empty
  #   }
  # }

  # Lifecycle options
  stop_on_destroy                      = false # Shutdown gracefully instead of force stop
  purge_on_destroy                     = true
  delete_unreferenced_disks_on_destroy = false # Safety: don't delete unmanaged disks

  timeouts = {
    create = "30m"
    update = "30m"
    delete = "10m"
  }
}

# Example 6: Linked clone for testing
resource "proxmox_virtual_environment_cloned_vm" "test_clone" {
  node_name = "pve"
  name      = "test-vm"

  clone = {
    source_vm_id = 100
    full         = false # Linked clone - faster, uses less space
  }

  cpu = {
    cores = 2
  }

  network = {
    net0 = {
      bridge = "vmbr0"
      model  = "virtio"
    }
  }
}

# Example 7: Clone with pool assignment
resource "proxmox_virtual_environment_cloned_vm" "pooled_clone" {
  node_name = "pve"
  name      = "pooled-vm"

  clone = {
    source_vm_id = 100
    pool_id      = "production" # Assign to pool during clone
  }

  cpu = {
    cores = 4
  }
}

# Example 8: Import existing cloned VM
resource "proxmox_virtual_environment_cloned_vm" "imported" {
  # Import with: terraform import proxmox_virtual_environment_cloned_vm.imported pve/123

  id        = 123 # VM ID to manage
  node_name = "pve"

  # After import, define managed configuration
  clone = {
    source_vm_id = 100 # Must match original clone source
  }

  cpu = {
    cores = 4
  }
}
