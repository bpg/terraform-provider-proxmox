---
layout: page
page_title: "Clone a VM"
subcategory: Guides
description: |-
    This guide explains how to create a VM template and then clone it to another VM.
---

# Clone a VM

## Create a VM template

VM templates in Proxmox provide an efficient way to create multiple identical VMs. Templates act as a base image that can be cloned to create new VMs, ensuring consistency and reducing the time needed to provision new instances. When a VM is created as a template, it is read-only and can't be started, but can be cloned multiple times to create new VMs.

You can create a template directly in Proxmox by setting the `template` attribute to `true` when creating the VM resource:

```terraform
resource "proxmox_virtual_environment_vm" "ubuntu_template" {
  name      = "ubuntu-template"
  node_name = "pve"

  template = true

  machine     = "q35"
  bios        = "ovmf"
  description = "Managed by Terraform"

  agent {
    enabled = true
  }

  cpu {
    cores = 2
  }

  memory {
    dedicated = 2048
  }

  efi_disk {
    datastore_id = "local"
    file_format  = "raw"
    type         = "4m"
  }

  disk {
    datastore_id = "local-lvm"
    file_id      = proxmox_virtual_environment_download_file.ubuntu_cloud_image.id
    interface    = "virtio0"
    iothread     = true
    discard      = "on"
    size         = 20
  }

  initialization {
    ip_config {
      ipv4 {
        address = "dhcp"
      }
    }

    user_data_file_id = proxmox_virtual_environment_file.user_data_cloud_config.id
  }

  network_device {
    bridge = "vmbr0"
  }

}

resource "proxmox_virtual_environment_download_file" "ubuntu_cloud_image" {
  content_type = "iso"
  datastore_id = "local"
  node_name    = "pve"

  url = "https://cloud-images.ubuntu.com/jammy/current/jammy-server-cloudimg-amd64.img"
}
```

Once you have a template, you can clone it to create new VMs. The cloned VMs will inherit all the configuration from the template but can be customized further as needed.

```terraform
resource "proxmox_virtual_environment_vm" "ubuntu_clone" {
  name      = "ubuntu-clone"
  node_name = "pve"

  clone {
    vm_id = proxmox_virtual_environment_vm.ubuntu_template.id
  }

  agent {
    enabled = true
  }

  memory {
    dedicated = 768
  }

  initialization {
    dns {
      servers = ["1.1.1.1"]
    }
    ip_config {
      ipv4 {
        address = "dhcp"
      }
    }
  }
}

output "vm_ipv4_address" {
  value = proxmox_virtual_environment_vm.ubuntu_clone.ipv4_addresses[1][0]
}
```
