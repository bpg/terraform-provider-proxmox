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

## Argument Reference

* `content_type` - (Optional) The content type.
    * `backup`
    * `iso`
    * `snippets`
    * `vztmpl`
* `datastore_id` - (Required) The datastore id.
* `node_name` - (Required) The node name.
* `source_file` - (Optional) The source file (conflicts with `source_raw`).
    * `checksum` - (Optional) The SHA256 checksum of the source file.
    * `file_name` - (Optional) The file name to use instead of the source file name.
    * `insecure` - (Optional) Whether to skip the TLS verification step for HTTPS sources (defaults to `false`).
    * `path` - (Required) A path to a local file or a URL.
* `source_raw` - (Optional) The raw source (conflicts with `source_file`).
    * `data` - (Required) The raw data.
    * `file_name` - (Required) The file name.
    * `resize` - (Optional) The number of bytes to resize the file to.

## Attribute Reference

* `file_modification_date` - The file modification date (RFC 3339).
* `file_name` - The file name.
* `file_size` - The file size in bytes.
* `file_tag` - The file tag.

## Important Notes

The Proxmox VE API endpoint for file uploads does not support chunked transfer encoding, which means that we must first
store the source file as a temporary file locally before uploading it.

You must ensure that you have at least `Size-in-MB * 2 + 1` MB of storage space available (twice the size plus overhead
because a multipart payload needs to be created as another temporary file).
