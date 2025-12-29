terraform {
  required_providers {
    proxmox = {
      source  = "bpg/proxmox"
      version = ">= 0.60.0"
    }
  }
}

provider "proxmox" {
  # Configuration options
}

data "proxmox_virtual_environment_file" "ubuntu_iso" {
  node_name    = "pve"
  datastore_id = "local"
  content_type = "iso"
  file_name    = "ubuntu-22.04.3-live-server-amd64.iso"
}

data "proxmox_virtual_environment_file" "ubuntu_container_template" {
  node_name    = "pve"
  datastore_id = "local"
  content_type = "vztmpl"
  file_name    = "ubuntu-22.04-standard_22.04-1_amd64.tar.zst"
}

data "proxmox_virtual_environment_file" "cloud_init_snippet" {
  node_name    = "pve"
  datastore_id = "local"
  content_type = "snippets"
  file_name    = "cloud-init-config.yaml"
}

data "proxmox_virtual_environment_file" "imported_file" {
  node_name    = "pve"
  datastore_id = "local"
  content_type = "import"
  file_name    = "imported-config.yaml"
}

output "ubuntu_iso_id" {
  value = data.proxmox_virtual_environment_file.ubuntu_iso.id
}

output "ubuntu_iso_size" {
  value = data.proxmox_virtual_environment_file.ubuntu_iso.file_size
}

output "container_template_format" {
  value = data.proxmox_virtual_environment_file.ubuntu_container_template.file_format
}

resource "proxmox_virtual_environment_vm" "example" {
  node_name = "pve"
  vm_id     = 100

  cdrom {
    file_id = data.proxmox_virtual_environment_file.ubuntu_iso.id
  }

  cpu {
    cores = 2
  }

  memory {
    dedicated = 2048
  }

  disk {
    datastore_id = "local-lvm"
    file_format  = "qcow2"
    size         = 20
  }

  network_device {
    bridge = "vmbr0"
  }
}
