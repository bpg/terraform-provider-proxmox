resource "proxmox_virtual_environment_oci_image" "ubuntu_latest" {
  node_name    = "pve"
  datastore_id = "local"
  reference    = "docker.io/library/ubuntu:latest"
}

resource "proxmox_virtual_environment_oci_image" "nginx" {
  node_name    = "pve"
  datastore_id = "local"
  reference    = "docker.io/library/nginx:alpine"
  file_name    = "custom_image_name.tar"
}

resource "proxmox_virtual_environment_oci_image" "debian" {
  node_name           = "pve"
  datastore_id        = "local"
  reference           = "docker.io/library/debian:bookworm"
  upload_timeout      = 900
  overwrite           = false
  overwrite_unmanaged = true
}
