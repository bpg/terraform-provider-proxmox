resource "proxmox_virtual_environment_file" "alpine_template" {
  datastore_id = "${element(data.proxmox_virtual_environment_datastores.example.datastore_ids, index(data.proxmox_virtual_environment_datastores.example.datastore_ids, "local"))}"
  node_name    = "${data.proxmox_virtual_environment_datastores.example.node_name}"
  source       = "${path.module}/assets/alpine-3.10-default_20190626_amd64.tar.xz"
  template     = true
}
