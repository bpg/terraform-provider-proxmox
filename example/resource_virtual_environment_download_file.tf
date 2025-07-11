## Debian and ubuntu image download

resource "proxmox_virtual_environment_download_file" "release_20250610_ubuntu_24_noble_lxc_img" {
  content_type        = "vztmpl"
  datastore_id        = "local"
  node_name           = var.virtual_environment_node_name
  url                 = var.release_20250610_ubuntu_24_noble_lxc_img_url
  checksum            = var.release_20250610_ubuntu_24_noble_lxc_img_checksum
  checksum_algorithm  = "sha256"
  upload_timeout      = 4444
  overwrite_unmanaged = true
}

resource "proxmox_virtual_environment_download_file" "latest_debian_12_bookworm_qcow2_img" {
  content_type        = "iso"
  datastore_id        = "local"
  file_name           = "debian-12-generic-amd64.img"
  node_name           = var.virtual_environment_node_name
  url                 = var.latest_debian_12_bookworm_qcow2_img_url
  overwrite           = true
  overwrite_unmanaged = true
}
