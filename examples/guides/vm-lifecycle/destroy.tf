resource "proxmox_virtual_environment_vm" "destroy_example" {
  name      = "destroy-example"
  node_name = var.virtual_environment_node_name

  started = true

  # When false (default), Terraform sends an ACPI shutdown and waits
  # for the guest to power off gracefully. Setting to true force-stops
  # the VM immediately, which is safer for started VMs without a
  # guest agent.
  stop_on_destroy = true

  # purge_on_destroy = true (default): Removes backup jobs, replication
  # entries, and HA configuration for this VM on destroy.
  purge_on_destroy = true

  # delete_unreferenced_disks_on_destroy = true (default for vm resource):
  # Deletes any disks not tracked by Terraform on destroy.
  # The cloned_vm resource defaults to false for safety.
  delete_unreferenced_disks_on_destroy = true

  cpu {
    cores = 2
  }

  memory {
    dedicated = 2048
  }

  disk {
    datastore_id = var.datastore_id
    interface    = "virtio0"
    size         = 20
  }

  network_device {
    bridge = "vmbr0"
  }
}
