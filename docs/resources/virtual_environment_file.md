---
layout: page
title: proxmox_virtual_environment_file
permalink: /resources/virtual_environment_file
nav_order: 6
parent: Resources
subcategory: Virtual Environment
---

# Resource: proxmox_virtual_environment_file

Manages a file.

## Example Usage

### Backups

-> **Note:** The resource with this content type uses SSH access to the node. You might need to configure the [`ssh` option in the `provider` section](../index.md#node-ip-address-used-for-ssh-connection).

```terraform
resource "proxmox_virtual_environment_file" "backup" {
  content_type = "backup"
  datastore_id = "local"
  node_name    = "pve"

  source_file {
    path = "vzdump-lxc-100-2023_11_08-23_10_05.tar"
  }
}
```

### Images

```terraform
resource "proxmox_virtual_environment_file" "ubuntu_container_template" {
  content_type = "iso"
  datastore_id = "local"
  node_name    = "pve"

  source_file {
    path = "https://cloud-images.ubuntu.com/jammy/20230929/jammy-server-cloudimg-amd64-disk-kvm.img"
  }
}
```

### Snippets

-> **Note:**  Snippets are not enabled by default in new Proxmox installations. You need to enable them in the 'Datacenter>Storage' section of the proxmox interface before first using this resource.

-> **Note:** The resource with this content type uses SSH access to the node. You might need to configure the [`ssh` option in the `provider` section](../index.md#node-ip-address-used-for-ssh-connection).

```terraform
resource "proxmox_virtual_environment_file" "cloud_config" {
  content_type = "snippets"
  datastore_id = "local"
  node_name    = "pve"

  source_raw {
    data = <<EOF
#cloud-config
chpasswd:
  list: |
    ubuntu:example
  expire: false
hostname: example-hostname
packages:
  - qemu-guest-agent
users:
  - default
  - name: ubuntu
    groups: sudo
    shell: /bin/bash
    ssh-authorized-keys:
      - ${trimspace(tls_private_key.example.public_key_openssh)}
    sudo: ALL=(ALL) NOPASSWD:ALL
EOF

    file_name = "example.cloud-config.yaml"
  }
}
```

### Container Template (`vztmpl`)

```terraform
resource "proxmox_virtual_environment_file" "ubuntu_container_template" {
  content_type = "vztmpl"
  datastore_id = "local"
  node_name    = "first-node"

  source_file {
    path = "https://download.proxmox.com/images/system/ubuntu-20.04-standard_20.04-1_amd64.tar.gz"
  }
}
```


## Argument Reference

- `content_type` - (Optional) The content type. If not specified, the content type will be inferred from the file
  extension. Valid values are:
    - `backup` (allowed extensions: `.vzdump`)
    - `iso` (allowed extensions: `.iso`, `.img`)
    - `snippets` (allowed extensions: any)
    - `vztmpl` (allowed extensions: `.tar.gz`, `.tar.xz`, `tar.zst`)
- `datastore_id` - (Required) The datastore id.
- `node_name` - (Required) The node name.
- `overwrite` - (Optional) Whether to overwrite an existing file (defaults to
  `true`).
- `source_file` - (Optional) The source file (conflicts with `source_raw`), could be a
  local file or a URL. If the source file is a URL, the file will be downloaded
  and stored locally before uploading it to Proxmox VE.
    - `checksum` - (Optional) The SHA256 checksum of the source file.
    - `file_name` - (Optional) The file name to use instead of the source file
      name. Useful when the source file does not have a valid file extension, for example 
      when the source file is a URL referencing a `.qcow2` image.
    - `insecure` - (Optional) Whether to skip the TLS verification step for
      HTTPS sources (defaults to `false`).
    - `path` - (Required) A path to a local file or a URL.
- `source_raw` - (Optional) The raw source (conflicts with `source_file`).
    - `data` - (Required) The raw data.
    - `file_name` - (Required) The file name.
    - `resize` - (Optional) The number of bytes to resize the file to.
- `timeout_upload` - (Optional) Timeout for uploading ISO/VSTMPL files in
  seconds (defaults to 1800).

## Attribute Reference

- `file_modification_date` - The file modification date (RFC 3339).
- `file_name` - The file name.
- `file_size` - The file size in bytes.
- `file_tag` - The file tag.

## Important Notes

The Proxmox VE API endpoint for file uploads does not support chunked transfer
encoding, which means that we must first store the source file as a temporary
file locally before uploading it.

You must ensure that you have at least `Size-in-MB * 2 + 1` MB of storage space
available (twice the size plus overhead because a multipart payload needs to be
created as another temporary file).

By default, if the specified file already exists, the resource will
unconditionally replace it and take ownership of the resource. On destruction,
the file will be deleted as if it did not exist before. If you want to prevent
the resource from replacing the file, set `overwrite` to `false`.

## Import

Instances can be imported using the `node_name`, `datastore_id`, `content_type`
and the `file_name` in the following format:

```text
node_name:datastore_id/content_type/file_name
```

Example:

```bash
$ terraform import proxmox_virtual_environment_file.cloud_config pve/local:snippets/example.cloud-config.yaml
```
