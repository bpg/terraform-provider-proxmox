# HOW-TO Create a VM from a Cloud Image

> [!NOTE]
> Examples below use the following defaults:
>
> - a single Proxmox node named `pve`
> - local storages named `local` and `local-lvm`

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
    file_id      = proxmox_virtual_environment_file.centos_cloud_image.id
    interface    = "virtio0"
    iothread     = true
    discard      = "on"
    size         = 20
  }
}

resource "proxmox_virtual_environment_file" "centos_cloud_image" {
  content_type = "iso"
  datastore_id = "local"
  node_name    = "pve"

  source_file {
    # you may download this image locally on your workstation and then use the local path instead of the remote URL
    path      = "https://cloud.centos.org/centos/8-stream/x86_64/images/CentOS-Stream-GenericCloud-8-20231113.0.x86_64.qcow2"
    file_name = "centos8.img"

    # you may also use the SHA256 checksum of the image to verify its integrity
    checksum = "b9ba602de681e493b020825db0ee30602a46ef92"
  }
}
```

Ubuntu cloud images are available at [cloud-images.ubuntu.com](https://cloud-images.ubuntu.com/). Ubuntu cloud images are in `qcow2` format as well, but stored with `.img` extension, so they can be directly uploaded to Proxmox without renaming.

Just update the `source_file` block in the example above to use the Ubuntu image URL:

```terraform
 source_file {
    path     = "https://cloud-images.ubuntu.com/jammy/current/jammy-server-cloudimg-amd64.img"
    ...
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
