---
layout: page
title: proxmox_virtual_environment_file
parent: Resources
subcategory: Virtual Environment
---

# Resource: proxmox_virtual_environment_file

Use this resource to upload files to a Proxmox VE node. The file can be a backup, an ISO image, a Disk Image, a snippet, or a container template depending on the `content_type` attribute.

## Example Usage

### Backups (`backup`)

-> The resource with this content type uses SSH access to the node. You might need to configure the [`ssh` option in the `provider` section](../index.md#node-ip-address-used-for-ssh-connection).

~> The provider currently does not support restoring backups. You can use the Proxmox VE web interface or the `qmrestore` / `pct restore` command to restore VM / Container from a backup.

```hcl
resource "proxmox_virtual_environment_file" "backup" {
  content_type = "backup"
  datastore_id = "local"
  node_name    = "pve"

  source_file {
    path = "vzdump-lxc-100-2023_11_08-23_10_05.tar.zst"
  }
}
```

### Images

-> Consider using `proxmox_virtual_environment_download_file` resource instead. Using this resource for images is less efficient (requires to transfer uploaded image to node) though still supported.

-> Importing Disks is not enabled by default in new Proxmox installations. You need to enable them in the 'Datacenter>Storage' section of the proxmox interface before first using this resource with `content_type = "import"`.

```hcl
resource "proxmox_virtual_environment_file" "ubuntu_container_template" {
  content_type = "iso"
  datastore_id = "local"
  node_name    = "pve"

  source_file {
    path = "https://cloud-images.ubuntu.com/jammy/20230929/jammy-server-cloudimg-amd64-disk-kvm.img"
  }
}
```

```hcl
resource "proxmox_virtual_environment_file" "ubuntu_container_template" {
  content_type = "import"
  datastore_id = "local"
  node_name    = "pve"

  source_file {
    path = "https://cloud-images.ubuntu.com/jammy/20230929/jammy-server-cloudimg-amd64-disk-kvm.img"
  }
}
```

### Snippets

-> Snippets are not enabled by default in new Proxmox installations. You need to enable them in the 'Datacenter>Storage' section of the proxmox interface before first using this resource.

-> The resource with this content type uses SSH access to the node. You might need to configure the [`ssh` option in the `provider` section](../index.md#node-ip-address-used-for-ssh-connection).

```hcl
resource "proxmox_virtual_environment_file" "cloud_config" {
  content_type = "snippets"
  datastore_id = "local"
  node_name    = "pve"

  source_raw {
    data = <<-EOF
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

The `file_mode` attribute can be used to make a script file executable, e.g. when referencing the file in the `hook_script_file_id` attribute of [a container](https://registry.terraform.io/providers/bpg/proxmox/latest/docs/resources/virtual_environment_container#hook_script_file_id) or [a VM](https://registry.terraform.io/providers/bpg/proxmox/latest/docs/resources/virtual_environment_vm#hook_script_file_id) resource which is a requirement enforced by the Proxmox VE API.

```hcl
resource "proxmox_virtual_environment_file" "hook_script" {
  content_type = "snippets"
  datastore_id = "local"
  node_name    = "pve"
  # Hook scripts must be executable, otherwise the Proxmox VE API will reject the configuration for the VM/CT.
  file_mode    = "0700"

  source_raw {
    data      = <<-EOF
      #!/usr/bin/env bash

      echo "Running hook script"
      EOF
    file_name = "prepare-hook.sh"
  }
}
```

### Container Template (`vztmpl`)

-> Consider using `proxmox_virtual_environment_download_file` resource instead. Using this resource for container images is less efficient (requires to transfer uploaded image to node) though still supported.

```hcl
resource "proxmox_virtual_environment_file" "ubuntu_container_template" {
  content_type = "vztmpl"
  datastore_id = "local"
  node_name    = "first-node"

  source_file {
    path = "http://download.proxmox.com/images/system/ubuntu-20.04-standard_20.04-1_amd64.tar.gz"
  }
}
```

## Argument Reference

- `content_type` - (Optional) The content type. If not specified, the content
    type will be inferred from the file extension. Valid values are:
    - `backup` (allowed extensions: `.vzdump`, `.tar.gz`, `.tar.xz`, `tar.zst`)
    - `iso` (allowed extensions: `.iso`, `.img`)
    - `snippets` (allowed extensions: any)
    - `import` (allowed extensions: `.raw`, `.qcow2`, `.vmdk`)
    - `vztmpl` (allowed extensions: `.tar.gz`, `.tar.xz`, `tar.zst`)
- `datastore_id` - (Required) The datastore id.
- `file_mode` - The file mode in octal format, e.g. `0700` or `600`. Note that the prefixes `0o` and `0x` is not supported! Setting this attribute is also only allowed for `root@pam` authenticated user.
- `node_name` - (Required) The node name.
- `overwrite` - (Optional) Whether to overwrite an existing file (defaults to
    `true`).
- `source_file` - (Optional) The source file (conflicts with `source_raw`),
    could be a local file or a URL. If the source file is a URL, the file will
    be downloaded and stored locally before uploading it to Proxmox VE.
    - `checksum` - (Optional) The SHA256 checksum of the source file.
    - `file_name` - (Optional) The file name to use instead of the source file
        name. Useful when the source file does not have a valid file extension,
        for example when the source file is a URL referencing a `.qcow2` image.
    - `insecure` - (Optional) Whether to skip the TLS verification step for
        HTTPS sources (defaults to `false`).
    - `min_tls` - (Optional) The minimum required TLS version for HTTPS
        sources. "Supported values: `1.0|1.1|1.2|1.3` (defaults to `1.3`).
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
terraform import proxmox_virtual_environment_file.cloud_config pve/local:snippets/example.cloud-config.yaml
```
