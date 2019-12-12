resource "proxmox_virtual_environment_file" "alpine_template" {
  datastore_id = "${element(data.proxmox_virtual_environment_datastores.example.datastore_ids, index(data.proxmox_virtual_environment_datastores.example.datastore_ids, "local"))}"
  node_name    = "${data.proxmox_virtual_environment_datastores.example.node_name}"
  source       = "http://download.proxmox.com/images/system/alpine-3.10-default_20190626_amd64.tar.xz"
  template     = true
}

output "resource_proxmox_virtual_environment_file_alpine_template_datastore_id" {
  value = "${proxmox_virtual_environment_file.alpine_template.datastore_id}"
}

output "resource_proxmox_virtual_environment_file_alpine_template_file_name" {
  value = "${proxmox_virtual_environment_file.alpine_template.file_name}"
}

output "resource_proxmox_virtual_environment_file_alpine_template_id" {
  value = "${proxmox_virtual_environment_file.alpine_template.id}"
}

output "resource_proxmox_virtual_environment_file_alpine_template_node_name" {
  value = "${proxmox_virtual_environment_file.alpine_template.node_name}"
}

output "resource_proxmox_virtual_environment_file_alpine_template_source" {
  value = "${proxmox_virtual_environment_file.alpine_template.source}"
}

output "resource_proxmox_virtual_environment_file_alpine_template_template" {
  value = "${proxmox_virtual_environment_file.alpine_template.template}"
}
