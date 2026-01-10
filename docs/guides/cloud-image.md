---
layout: page
page_title: "Create a VM from a Cloud Image"
subcategory: Guides
description: |-
    This guide explains how to create a VM from a cloud image.
---

# Create a VM from a Cloud Image

## Download a public cloud image from URL

Example of how to create a CentOS 8 VM from a "generic cloud" `qcow2` image. CentOS 8 images are available at [cloud.centos.org](https://cloud.centos.org/centos/8-stream/x86_64/images/):

```terraform
resource "proxmox_virtual_environment_vm" "centos_vm" {
  name      = "test-centos"
  node_name = "pve"

  # should be true if qemu agent is not installed / enabled on the VM
  stop_on_destroy = true

  initialization {
    user_account {
      # do not use this in production, configure your own ssh key instead!
      username = "user"
      password = "password"
    }
  }

  disk {
    datastore_id = "local-lvm"
    import_from  = proxmox_virtual_environment_download_file.centos_cloud_image.id
    interface    = "virtio0"
    iothread     = true
    discard      = "on"
    size         = 20
  }
}

resource "proxmox_virtual_environment_download_file" "centos_cloud_image" {
  content_type = "import"
  datastore_id = "local"
  node_name    = "pve"
  url          = "https://cloud.centos.org/centos/8-stream/x86_64/images/CentOS-Stream-GenericCloud-8-latest.x86_64.qcow2"
}
```

Ubuntu cloud images are available at [cloud-images.ubuntu.com](https://cloud-images.ubuntu.com/). Ubuntu cloud images are in `qcow2` format as well, but stored with `.img` extension, so they can be directly uploaded to Proxmox without renaming.

```terraform
resource "proxmox_virtual_environment_vm" "ubuntu_vm" {
  name      = "test-ubuntu"
  node_name = "pve"

  # should be true if qemu agent is not installed / enabled on the VM
  stop_on_destroy = true

  initialization {
    user_account {
      # do not use this in production, configure your own ssh key instead!
      username = "user"
      password = "password"
    }
  }

  disk {
    datastore_id = "local-lvm"
    import_from  = proxmox_virtual_environment_download_file.ubuntu_cloud_image.id
    interface    = "virtio0"
    iothread     = true
    discard      = "on"
    size         = 20
  }
}

resource "proxmox_virtual_environment_download_file" "ubuntu_cloud_image" {
  content_type = "import"
  datastore_id = "local"
  node_name    = "pve"
  url          = "https://cloud-images.ubuntu.com/jammy/current/jammy-server-cloudimg-amd64.img"
  # need to rename the file to *.qcow2 to indicate the actual file format for import
  file_name = "jammy-server-cloudimg-amd64.qcow2"
}
```

## Create a VM from a compressed cloud image

Some distributions, like Fedora CoreOS, only provide cloud images in compressed formats (e.g., `.qcow2.xz`). Proxmox does not directly support XZ-compressed images, but its ZST decompression can handle XZ archives.

~> **Important:** Compressed images cannot be used with `import_from`. You must use `file_id` with `content_type = "iso"` instead. This requires SSH configuration in the provider.

Example of how to create a Fedora CoreOS VM from an XZ-compressed `qcow2` image. CoreOS images are available at [fedoraproject.org/coreos](https://fedoraproject.org/coreos/download):

```terraform
data "http" "coreos_stable_metadata" {
  url = "https://builds.coreos.fedoraproject.org/streams/stable.json"
}

locals {
  coreos_metadata     = jsondecode(data.http.coreos_stable_metadata.response_body)
  coreos_qemu_stable  = local.coreos_metadata.architectures.x86_64.artifacts.qemu.formats["qcow2.xz"].disk
  coreos_download_url = local.coreos_qemu_stable.location
  coreos_checksum     = local.coreos_qemu_stable.sha256
}

resource "proxmox_virtual_environment_download_file" "coreos_image" {
  content_type = "iso"
  datastore_id = "local"
  node_name    = "pve"

  url                = local.coreos_download_url
  checksum           = local.coreos_checksum
  checksum_algorithm = "sha256"

  # use zst decompression for xz-compressed images
  decompression_algorithm = "zst"
  # rename to .img to avoid Proxmox file extension validation errors
  file_name = "fedora-coreos-stable-qemu.qcow2.img"
}

resource "proxmox_virtual_environment_vm" "coreos_vm" {
  name      = "test-coreos"
  node_name = "pve"

  # CoreOS does not have qemu-guest-agent by default
  stop_on_destroy = true

  cpu {
    cores = 2
  }

  memory {
    dedicated = 2048
  }

  disk {
    datastore_id = "local-lvm"
    # use file_id instead of import_from for decompressed images
    file_id   = proxmox_virtual_environment_download_file.coreos_image.id
    interface = "virtio0"
    iothread  = true
    discard   = "on"
    size      = 20
  }

  network_device {
    bridge = "vmbr0"
  }
}
```

## Create a VM from an existing image on Proxmox

If you already have a cloud image on Proxmox, you can use it to create a VM:

```terraform
resource "proxmox_virtual_environment_vm" "debian_vm" {
  name      = "test-debian"
  node_name = "pve"

  # should be true if qemu agent is not installed / enabled on the VM
  stop_on_destroy = true

  initialization {
    user_account {
      # do not use this in production, configure your own ssh key instead!
      username = "user"
      password = "password"
    }
  }

  disk {
    datastore_id = "local-lvm"
    # qcow2 image downloaded from https://cloud.debian.org/images/cloud/bookworm/latest/ and renamed to *.img
    # the image is not of import type, so provider will use SSH client to import it
    file_id   = "local:iso/debian-12-genericcloud-amd64.img"
    interface = "virtio0"
    iothread  = true
    discard   = "on"
    size      = 20
  }
}
```
