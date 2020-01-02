resource "proxmox_virtual_environment_container" "example" {
  description = "Managed by Terraform"

  initialization {
    dns {
      server = "1.1.1.1"
    }

    hostname = "terraform-provider-proxmox-example-lxc"

    ip_config {
      ipv4 {
        address = "dhcp"
      }
    }

    user_account {
      keys     = ["${trimspace(tls_private_key.example.public_key_openssh)}"]
      password = "example"
    }
  }

  network_interface {
    name = "veth0"
  }

  node_name = "${data.proxmox_virtual_environment_nodes.example.names[0]}"

  operating_system {
    template_file_id = "${proxmox_virtual_environment_file.ubuntu_container_template.id}"
    type             = "ubuntu"
  }

  pool_id = "${proxmox_virtual_environment_pool.example.id}"
  vm_id   = 2039
}

output "resource_proxmox_virtual_environment_container_example_id" {
  value = "${proxmox_virtual_environment_container.example.id}"
}
