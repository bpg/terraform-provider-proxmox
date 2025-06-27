resource "proxmox_virtual_environment_vm" "trunks-example" {
  name        = "trunks-example"
  node_name   = data.proxmox_virtual_environment_nodes.example.names[0]
  description = "Example of a VM using trunks to pass multiple VLANs on a single network interface."

  disk {
    datastore_id = local.datastore_id
    file_id      = proxmox_virtual_environment_download_file.latest_debian_12_bookworm_qcow2_img.id
    interface    = "scsi0"
    discard      = "on"
    cache        = "writeback"
    ssd          = true
  }

  initialization {
    datastore_id = local.datastore_id
    interface    = "scsi4"

    dns {
      servers = ["1.1.1.1", "8.8.8.8"]
    }

    ip_config {
      ipv4 {
        address = "dhcp"
      }
    }
    user_data_file_id   = proxmox_virtual_environment_file.user_config.id
    vendor_data_file_id = proxmox_virtual_environment_file.vendor_config.id
    meta_data_file_id   = proxmox_virtual_environment_file.meta_config.id
  }

  memory {
    dedicated = 1024
  }

  cpu {
    cores = 2
  }

  agent {
    enabled = true
  }

  serial_device {}

  boot_order    = ["scsi0"]
  scsi_hardware = "virtio-scsi-pci"

  network_device {
    model  = "virtio"
    bridge = "vmbr0"
    trunks = "10;20;30"
  }
}
