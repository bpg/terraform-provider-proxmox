## Debian and ubuntu image download

resource "proxmox_virtual_environment_download_file" "debian12_image" {
  content_type            = "iso"
  datastore_id            = "local"
  node_name               = "pve"
  download_url            = "https://cloud.debian.org/images/cloud/bookworm/20230802-1460/debian-12-generic-amd64-20230802-1460.qcow2"
  allow_unsupported_types = true
  upload_timeout          = 4444
}

resource "proxmox_virtual_environment_download_file" "ubuntu_noble_image" {
  content_type       = "iso"
  datastore_id       = "local"
  node_name          = "pve"
  download_url       = "https://cloud-images.ubuntu.com/noble/20231218/noble-server-cloudimg-amd64.img"
  upload_timeout     = 5555
  checksum           = "ee57ee52e80d1cdd4121043dd5109c7869433e5e2a78bfa6058881fc1d535e03"
  checksum_algorithm = "sha256"
}
