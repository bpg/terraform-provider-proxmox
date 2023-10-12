locals {
  datastore_id = element(data.proxmox_virtual_environment_datastores.example.datastore_ids, index(data.proxmox_virtual_environment_datastores.example.datastore_ids, "local-lvm"))
}

resource "proxmox_virtual_environment_vm" "example_template" {
  agent {
    enabled = true
  }

  bios        = "ovmf"
  description = "Managed by Terraform"

  cpu {
    cores = 2
    numa  = true
  }

  smbios {
    manufacturer = "Terraform"
    product      = "Terraform Provider Proxmox"
    version      = "0.0.1"
  }

  startup {
    order      = "3"
    up_delay   = "60"
    down_delay = "60"
  }

  efi_disk {
    datastore_id = local.datastore_id
    file_format  = "raw"
    type         = "4m"
  }

  #  disk {
  #    datastore_id = local.datastore_id
  #    file_id      = proxmox_virtual_environment_file.ubuntu_cloud_image.id
  #    interface    = "virtio0"
  #    iothread     = true
  #  }

  disk {
    datastore_id = local.datastore_id
    file_id      = proxmox_virtual_environment_file.ubuntu_cloud_image.id
    interface    = "scsi0"
    discard      = "on"
    cache        = "writeback"
    ssd          = true
  }

  #  disk {
  #    datastore_id = "nfs"
  #    interface    = "scsi1"
  #    discard      = "ignore"
  #    file_format  = "raw"
  #  }

  initialization {
    datastore_id = local.datastore_id
    interface    = "scsi4"

    dns {
      server = "1.1.1.1"
    }

    ip_config {
      ipv4 {
        address = "dhcp"
      }
      # ipv6 {
      #    address = "dhcp" 
      #}
    }

    user_data_file_id   = proxmox_virtual_environment_file.user_config.id
    vendor_data_file_id = proxmox_virtual_environment_file.vendor_config.id
    meta_data_file_id   = proxmox_virtual_environment_file.meta_config.id
  }

  machine = "q35"
  name    = "terraform-provider-proxmox-example-template"

  network_device {
    mtu    = 1450
    queues = 2
  }

  network_device {
    vlan_id = 1024
  }

  node_name = data.proxmox_virtual_environment_nodes.example.names[0]

  operating_system {
    type = "l26"
  }

  pool_id = proxmox_virtual_environment_pool.example.id

  serial_device {}

  template = true

  // use auto-generated vm_id
}

resource "proxmox_virtual_environment_vm" "example" {
  name      = "terraform-provider-proxmox-example"
  node_name = data.proxmox_virtual_environment_nodes.example.names[0]
  migrate   = true // migrate the VM on node change
  pool_id   = proxmox_virtual_environment_pool.example.id
  vm_id     = 2041
  tags      = ["terraform", "ubuntu"]

  clone {
    vm_id = proxmox_virtual_environment_vm.example_template.id
  }

  machine = "q35"

  memory {
    dedicated = 768
  }

  connection {
    type        = "ssh"
    agent       = false
    host        = element(element(self.ipv4_addresses, index(self.network_interface_names, "eth0")), 0)
    private_key = tls_private_key.example.private_key_pem
    user        = "ubuntu"
  }

  provisioner "remote-exec" {
    inline = [
      "echo Welcome to $(hostname)!",
    ]
  }

  initialization {
    // if unspecified:
    //   - autodetected if there is a cloud-init device on the template
    //   - otherwise defaults to ide2
    interface = "scsi4"

    dns {
      server = "8.8.8.8"
    }
    ip_config {
      ipv4 {
        address = "dhcp"
      }
    }
  }

  #hostpci {
  #  device = "hostpci0"
  #  id = "0000:00:1f.0"
  #  pcie = true
  #}

  #hostpci {
  #  device = "hostpci1"
  #  mapping = "gpu"
  #  pcie = true
  #}

  # attached disks from data_vm
  dynamic "disk" {
    for_each = {for idx, val in proxmox_virtual_environment_vm.data_vm.disk : idx => val}
    iterator = data_disk
    content {
      datastore_id      = data_disk.value["datastore_id"]
      path_in_datastore = data_disk.value["path_in_datastore"]
      file_format       = data_disk.value["file_format"]
      size              = data_disk.value["size"]
      # assign from scsi1 and up
      interface         = "scsi${data_disk.key + 1}"
    }
  }
}

resource "proxmox_virtual_environment_vm" "data_vm" {
  name      = "terraform-provider-proxmox-data-vm"
  node_name = data.proxmox_virtual_environment_nodes.example.names[0]
  started   = false
  on_boot   = false

  disk {
    datastore_id = local.datastore_id
    file_format  = "raw"
    interface    = "scsi0"
    size         = 1
  }
  disk {
    datastore_id = local.datastore_id
    file_format  = "raw"
    interface    = "scsi1"
    size         = 4
  }
}

output "resource_proxmox_virtual_environment_vm_example_id" {
  value = proxmox_virtual_environment_vm.example.id
}

output "resource_proxmox_virtual_environment_vm_example_ipv4_addresses" {
  value = proxmox_virtual_environment_vm.example.ipv4_addresses
}

output "resource_proxmox_virtual_environment_vm_example_ipv6_addresses" {
  value = proxmox_virtual_environment_vm.example.ipv6_addresses
}

output "resource_proxmox_virtual_environment_vm_example_mac_addresses" {
  value = proxmox_virtual_environment_vm.example.mac_addresses
}

output "resource_proxmox_virtual_environment_vm_example_network_interface_names" {
  value = proxmox_virtual_environment_vm.example.network_interface_names
}
