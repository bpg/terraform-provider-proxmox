data "proxmox_files" "iso_files" {
  node_name    = "pve"
  datastore_id = "local"
  content_type = "iso"
}

# Check if a specific image already exists
locals {
  image_exists = anytrue([
    for f in data.proxmox_files.iso_files.files :
    f.file_name == "noble-server-cloudimg-amd64.img"
  ])
}

# Only download if the image doesn't already exist
resource "proxmox_virtual_environment_download_file" "ubuntu_noble" {
  count = local.image_exists ? 0 : 1

  datastore_id = "local"
  node_name    = "pve"
  content_type = "iso"
  url          = "https://cloud-images.ubuntu.com/noble/current/noble-server-cloudimg-amd64.img"
}

# List all files without filtering
data "proxmox_files" "all_files" {
  node_name    = "pve"
  datastore_id = "local"
}

output "iso_file_count" {
  value = length(data.proxmox_files.iso_files.files)
}

output "all_file_names" {
  value = [for f in data.proxmox_files.all_files.files : f.file_name]
}
