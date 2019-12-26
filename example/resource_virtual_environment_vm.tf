resource "proxmox_virtual_environment_vm" "example" {
  cloud_init {
    dns {
      server = "1.1.1.1"
    }

    ip_config {
      ipv4 {
        address = "dhcp"
      }
    }

    user_account {
      keys     = ["${trimspace(tls_private_key.example.public_key_openssh)}"]
      username = "ubuntu"
    }
  }

  network_device {}

  node_name = "${data.proxmox_virtual_environment_nodes.example.names[0]}"
}

resource "tls_private_key" "example" {
  algorithm = "RSA"
  rsa_bits  = 2048
}
