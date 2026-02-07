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
    disk_size = "20G"
  }

  network_device {
    bridge = "vmbr0"
  }
}

