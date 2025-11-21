## Debian and ubuntu image download

resource "proxmox_virtual_environment_download_file" "ubuntu_24_04_lxc_img" {
  content_type        = "vztmpl"
  datastore_id        = "local"
  node_name           = data.proxmox_virtual_environment_nodes.example.names[0]
  url                 = var.ubuntu_24_04_lxc_img_url
  upload_timeout      = 4444
  overwrite_unmanaged = true
}

resource "proxmox_virtual_environment_download_file" "latest_debian_12_bookworm_qcow2_img" {
  content_type        = "import"
  datastore_id        = "local"
  file_name           = "debian-12-generic-amd64.qcow2"
  node_name           = data.proxmox_virtual_environment_nodes.example.names[0]
  url                 = var.latest_debian_12_bookworm_qcow2_img_url
  overwrite           = true
  overwrite_unmanaged = true
}
