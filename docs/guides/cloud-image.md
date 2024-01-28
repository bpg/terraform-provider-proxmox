---
layout: page
page_title: "Create a VM from a Cloud Image"
subcategory: Guides
description: |-
    This guide explains how to create a VM from a cloud image.
---

# Create a VM from a Cloud Image

## Download a public cloud image from URL

Proxmox does not natively support QCOW2 images, but provider can do the conversion for you.

Example of how to create a CentOS 8 VM from a "generic cloud" `qcow2` image. CentOS 8 images are available at [cloud.centos.org](https://cloud.centos.org/centos/8-stream/x86_64/images/):

```terraform
resource "proxmox_virtual_environment_vm" "centos_vm" {
  name      = "test-centos"
  node_name = "pve"

  initialization {
    user_account {
      # do not use this in production, configure your own ssh key instead!
      username = "user"
      password = "password"
    }
  }

  disk {
    datastore_id = "local-lvm"
    file_id      = proxmox_virtual_environment_download_file.centos_cloud_image.id
    interface    = "virtio0"
    iothread     = true
    discard      = "on"
    size         = 20
  }
}

resource "proxmox_virtual_environment_download_file" "centos_cloud_image" {
  content_type = "iso"
  datastore_id = "local"
  node_name    = "pve"
  url          = "https://cloud.centos.org/centos/8-stream/x86_64/images/CentOS-Stream-GenericCloud-8-20231113.0.x86_64.qcow2"
  file_name    = "centos8.img"
}
```

Ubuntu cloud images are available at [cloud-images.ubuntu.com](https://cloud-images.ubuntu.com/). Ubuntu cloud images are in `qcow2` format as well, but stored with `.img` extension, so they can be directly uploaded to Proxmox without renaming.

```terraform
resource "proxmox_virtual_environment_vm" "ubuntu_vm" {
  name      = "test-ubuntu"
  node_name = "pve"

  initialization {
    user_account {
      # do not use this in production, configure your own ssh key instead!
      username = "user"
      password = "password"
    }
  }

  disk {
    datastore_id = "local-lvm"
    file_id      = proxmox_virtual_environment_download_file.ubuntu_cloud_image.id
    interface    = "virtio0"
    iothread     = true
    discard      = "on"
    size         = 20
  }
}

resource "proxmox_virtual_environment_download_file" "ubuntu_cloud_image" {
  content_type = "iso"
  datastore_id = "local"
  node_name    = "pve"
  url          = "https://cloud-images.ubuntu.com/jammy/current/jammy-server-cloudimg-amd64.img"

}
```

For [large images](https://registry.terraform.io/providers/bpg/proxmox/latest/docs/resources/virtual_environment_file#important-notes), you may want to use a dedicated temporary directory [configured](https://registry.terraform.io/providers/bpg/proxmox/latest/docs#tmp_dir) for provider via `tmp_dir` attribute, instead of system's default temporary directory. This is especially useful if you are deploying from a container with limited disk space.

## Create a VM from an exiting image on Proxmox

If you already have a cloud image on Proxmox, you can use it to create a VM:

```terraform
resource "proxmox_virtual_environment_vm" "debian_vm" {
  name      = "test-debian"
  node_name = "pve"

  initialization {
    user_account {
      # do not use this in production, configure your own ssh key instead!
      username = "user"
      password = "password"
    }
  }

  disk {
    datastore_id = "local-lvm"
    file_id      = "local:iso/debian-12-genericcloud-amd64.img"
    interface    = "virtio0"
    iothread     = true
    discard      = "on"
    size         = 20
  }
}
```
