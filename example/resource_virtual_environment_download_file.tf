## Debian and ubuntu image download

resource "proxmox_virtual_environment_download_file" "release_20250701_ubuntu_24_10_lxc_img" {
  content_type        = "vztmpl"
  datastore_id        = "local"
  node_name           = var.virtual_environment_node_name
  url                 = var.release_20250701_ubuntu_24_10_lxc_img_url
  checksum            = var.release_20250701_ubuntu_24_10_lxc_img_checksum
  checksum_algorithm  = "sha256"
  upload_timeout      = 4444
  overwrite_unmanaged = true
}

resource "proxmox_virtual_environment_download_file" "latest_debian_12_bookworm_qcow2_img" {
  content_type        = "import"
  datastore_id        = "local"
  file_name           = "debian-12-generic-amd64.qcow2"
  node_name           = var.virtual_environment_node_name
  url                 = var.latest_debian_12_bookworm_qcow2_img_url
  overwrite           = true
  overwrite_unmanaged = true
}
