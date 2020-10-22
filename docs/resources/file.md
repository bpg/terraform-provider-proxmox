---
layout: page
title: File
permalink: /resources/file
nav_order: 4
parent: Resources
---

# Resource: File

Manages a file.

## Example Usage

```
resource "proxmox_virtual_environment_file" "ubuntu_container_template" {
  content_type = "vztmpl"
  datastore_id = "local"
  node_name    = "first-node"

  source_file {
    path = "http://download.proxmox.com/images/system/ubuntu-18.04-standard_18.04.1-1_amd64.tar.gz"
  }
}
```

## Arguments Reference

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

## Attributes Reference

* `file_modification_date` - The file modification date (RFC 3339).
* `file_name` - The file name.
* `file_size` - The file size in bytes.
* `file_tag` - The file tag.

## Important Notes

The Proxmox VE API endpoint for file uploads does not support chunked transfer encoding, which means that we must first store the source file as a temporary file locally before uploading it.

You must ensure that you have at least `Size-in-MB * 2 + 1` MB of storage space available (twice the size plus overhead because a multipart payload needs to be created as another temporary file).
