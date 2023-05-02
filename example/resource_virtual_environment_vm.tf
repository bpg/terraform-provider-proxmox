locals {
  datastore_id = element(data.proxmox_virtual_environment_datastores.example.datastore_ids, index(data.proxmox_virtual_environment_datastores.example.datastore_ids, "local-lvm"))
}

resource "proxmox_virtual_environment_vm" "example_template" {
  agent {
    enabled = true
  }

  description = "Managed by Terraform"

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

    user_data_file_id = proxmox_virtual_environment_file.user_config.id
    vendor_data_file_id = proxmox_virtual_environment_file.vendor_config.id
  }

  name = "terraform-provider-proxmox-example-template"

  network_device {
    mtu  = 1450
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
  vm_id    = 2040
}

resource "proxmox_virtual_environment_vm" "example" {
  name      = "terraform-provider-proxmox-example"
  node_name = data.proxmox_virtual_environment_nodes.example.names[0]
  pool_id   = proxmox_virtual_environment_pool.example.id
  vm_id     = 2041
  tags        = ["terraform", "ubuntu"]

  clone {
    vm_id = proxmox_virtual_environment_vm.example_template.id
  }

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

# While overwriting the initialization block when cloning a template is possible, it is not recommended.
# This will cause the coned VM to be reinitialized each time on re-apply.
#  initialization {
#    dns {
#      server = "8.8.8.8"
#    }
#  }

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
