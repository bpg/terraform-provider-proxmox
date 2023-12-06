resource "proxmox_virtual_environment_container" "example_template" {
  description = "Managed by Terraform"

  disk {
    datastore_id = element(data.proxmox_virtual_environment_datastores.example.datastore_ids, index(data.proxmox_virtual_environment_datastores.example.datastore_ids, "local-lvm"))
    size         = 10
  }

  mount_point {
    // volume mount
    volume = element(data.proxmox_virtual_environment_datastores.example.datastore_ids, index(data.proxmox_virtual_environment_datastores.example.datastore_ids, "local-lvm"))
    size   = "4G"
    path   = "mnt/local"
  }

  initialization {
    dns {
      server = "1.1.1.1"
    }

    hostname = "terraform-provider-proxmox-example-lxc-template"

    ip_config {
      ipv4 {
        address = "dhcp"
      }
    }

    user_account {
      keys     = [trimspace(tls_private_key.example.public_key_openssh)]
      password = "example"
    }
  }

  network_interface {
    name = "veth0"
    mtu  = 1450
  }

  node_name = data.proxmox_virtual_environment_nodes.example.names[0]

  operating_system {
    template_file_id = proxmox_virtual_environment_file.ubuntu_container_template.id
    type             = "ubuntu"
  }

  pool_id  = proxmox_virtual_environment_pool.example.id
  template = true

  // use auto-generated vm_id

  tags = [
    "container",
    "example",
    "terraform",
  ]
}

resource "proxmox_virtual_environment_container" "example" {
  disk {
    datastore_id = element(data.proxmox_virtual_environment_datastores.example.datastore_ids, index(data.proxmox_virtual_environment_datastores.example.datastore_ids, "local-lvm"))
  }

  clone {
    vm_id = proxmox_virtual_environment_container.example_template.id
  }

  initialization {
    hostname = "terraform-provider-proxmox-example-lxc"
  }

  mount_point {
    // bind mount, requires root@pam
    volume = "/mnt/bindmounts/shared"
    path    = "/shared"
  }

  node_name = data.proxmox_virtual_environment_nodes.example.names[0]
  pool_id   = proxmox_virtual_environment_pool.example.id
  vm_id     = 2043
}

output "resource_proxmox_virtual_environment_container_example_id" {
  value = proxmox_virtual_environment_container.example.id
}
