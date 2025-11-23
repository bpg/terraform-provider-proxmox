resource "proxmox_virtual_environment_oci_image" "ubuntu_latest" {
  node_name    = "pve"
  datastore_id = "local"
  reference    = "docker.io/library/ubuntu:latest"
}

# Pull an OCI image with a custom filename
resource "proxmox_virtual_environment_oci_image" "nginx" {
  node_name    = "pve"
  datastore_id = "local"
  reference    = "docker.io/library/nginx:alpine"
  file_name    = "nginx_alpine.tar"
}

# Pull an OCI image with custom timeout and overwrite settings
resource "proxmox_virtual_environment_oci_image" "debian" {
  node_name           = "pve"
  datastore_id        = "local"
  reference           = "docker.io/library/debian:bookworm"
  upload_timeout      = 900 # 15 minutes
  overwrite           = false
  overwrite_unmanaged = true
}
