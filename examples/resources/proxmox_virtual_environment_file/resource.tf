resource "proxmox_virtual_environment_file" "latest_debian_12_bookworm_qcow2" {
  content_type = "import"
  datastore_id = "local"
  node_name    = "pve"

  source_file {
    path = "https://cloud.debian.org/images/cloud/bookworm/latest/debian-12-generic-amd64.qcow2"
  }
}

resource "proxmox_virtual_environment_file" "release_20231228_debian_12_bookworm_qcow2" {
  content_type = "import"
  datastore_id = "local"
  node_name    = "pve"

  source_file {
    path = "https://cloud.debian.org/images/cloud/bookworm/20231228-1609/debian-12-generic-amd64-20231228-1609.qcow2"
  }
}
