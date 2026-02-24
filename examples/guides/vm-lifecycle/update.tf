resource "proxmox_virtual_environment_vm" "hotplug_example" {
  name      = "hotplug-example"
  node_name = var.virtual_environment_node_name

  started         = true
  stop_on_destroy = true

  # reboot_after_update defaults to true.
  # When a non-hotpluggable attribute changes (e.g. cpu.cores),
  # Terraform will automatically reboot the VM to apply it.
  # Set to false if you prefer to reboot manually.
  reboot_after_update = true

  cpu {
    # Changing cores or sockets requires a reboot.
    cores   = 4
    sockets = 1
  }

  memory {
    # Memory is hotpluggable (increase only) when the VM's hotplug
    # setting includes "memory".
    dedicated = 4096
  }

  disk {
    datastore_id = var.datastore_id
    interface    = "virtio0"
    iothread     = true
    discard      = "on"
    # Disks can only grow. Shrinking produces an error.
    size = 40
  }

  network_device {
    bridge = "vmbr0"
  }

  # Adding a second NIC is hotpluggable.
  network_device {
    bridge  = "vmbr0"
    vlan_id = 100
  }
}
