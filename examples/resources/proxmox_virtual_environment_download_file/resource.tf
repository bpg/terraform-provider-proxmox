resource "proxmox_virtual_environment_download_file" "release_20231228_debian_12_bookworm_qcow2_img" {
  content_type       = "iso"
  datastore_id       = "local"
  file_name          = "debian-12-generic-amd64-20231228-1609.img"
  node_name          = "pve"
  url                = "https://cloud.debian.org/images/cloud/bookworm/20231228-1609/debian-12-generic-amd64-20231228-1609.qcow2"
  checksum           = "d2fbcf11fb28795842e91364d8c7b69f1870db09ff299eb94e4fbbfa510eb78d141e74c1f4bf6dfa0b7e33d0c3b66e6751886feadb4e9916f778bab1776bdf1b"
  checksum_algorithm = "sha512"
}

resource "proxmox_virtual_environment_download_file" "release_20231228_debian_12_bookworm_qcow2" {
  content_type       = "import"
  datastore_id       = "local"
  file_name          = "debian-12-generic-amd64-20231228-1609.qcow2"
  node_name          = "pve"
  url                = "https://cloud.debian.org/images/cloud/bookworm/20231228-1609/debian-12-generic-amd64-20231228-1609.qcow2"
  checksum           = "d2fbcf11fb28795842e91364d8c7b69f1870db09ff299eb94e4fbbfa510eb78d141e74c1f4bf6dfa0b7e33d0c3b66e6751886feadb4e9916f778bab1776bdf1b"
  checksum_algorithm = "sha512"
}

resource "proxmox_virtual_environment_download_file" "latest_debian_12_bookworm_qcow2_img" {
  content_type = "iso"
  datastore_id = "local"
  file_name    = "debian-12-generic-amd64.qcow2.img"
  node_name    = "pve"
  url          = "https://cloud.debian.org/images/cloud/bookworm/latest/debian-12-generic-amd64.qcow2"
}

resource "proxmox_virtual_environment_download_file" "latest_debian_12_bookworm_qcow2" {
  content_type = "import"
  datastore_id = "local"
  file_name    = "debian-12-generic-amd64.qcow2"
  node_name    = "pve"
  url          = "https://cloud.debian.org/images/cloud/bookworm/latest/debian-12-generic-amd64.qcow2"
}

resource "proxmox_virtual_environment_download_file" "latest_ubuntu_22_jammy_qcow2_img" {
  content_type = "iso"
  datastore_id = "local"
  node_name    = "pve"
  url          = "https://cloud-images.ubuntu.com/jammy/current/jammy-server-cloudimg-amd64.img"
}

resource "proxmox_virtual_environment_download_file" "latest_static_ubuntu_24_noble_qcow2_img" {
  content_type = "iso"
  datastore_id = "local"
  node_name    = "pve"
  url          = "https://cloud-images.ubuntu.com/noble/current/noble-server-cloudimg-amd64.img"
  overwrite    = false
}

resource "proxmox_virtual_environment_download_file" "release_20231211_ubuntu_22_jammy_lxc_img" {
  content_type       = "vztmpl"
  datastore_id       = "local"
  node_name          = "pve"
  url                = "https://cloud-images.ubuntu.com/releases/22.04/release-20231211/ubuntu-22.04-server-cloudimg-amd64-root.tar.xz"
  checksum           = "c9997dcfea5d826fd04871f960c513665f2e87dd7450bba99f68a97e60e4586e"
  checksum_algorithm = "sha256"
  upload_timeout     = 4444
}

resource "proxmox_virtual_environment_download_file" "latest_ubuntu_22_jammy_lxc_img" {
  content_type = "vztmpl"
  datastore_id = "local"
  node_name    = "pve"
  url          = "https://cloud-images.ubuntu.com/jammy/current/jammy-server-cloudimg-amd64.tar.gz"
}
