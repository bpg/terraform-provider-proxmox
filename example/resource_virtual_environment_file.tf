resource "proxmox_virtual_environment_file" "ubuntu_cloud_image" {
  content_type    = "iso"
  datastore_id    = "${element(data.proxmox_virtual_environment_datastores.example.datastore_ids, index(data.proxmox_virtual_environment_datastores.example.datastore_ids, "local"))}"
  node_name       = "${data.proxmox_virtual_environment_datastores.example.node_name}"
  source          = "https://cloud-images.ubuntu.com/bionic/current/bionic-server-cloudimg-amd64.img"
}

output "resource_proxmox_virtual_environment_file_ubuntu_cloud_image_content_type" {
  value = "${proxmox_virtual_environment_file.ubuntu_cloud_image.content_type}"
}

output "resource_proxmox_virtual_environment_file_ubuntu_cloud_image_datastore_id" {
  value = "${proxmox_virtual_environment_file.ubuntu_cloud_image.datastore_id}"
}

output "resource_proxmox_virtual_environment_file_ubuntu_cloud_image_file_modification_date" {
  value = "${proxmox_virtual_environment_file.ubuntu_cloud_image.file_modification_date}"
}

output "resource_proxmox_virtual_environment_file_ubuntu_cloud_image_file_name" {
  value = "${proxmox_virtual_environment_file.ubuntu_cloud_image.file_name}"
}

output "resource_proxmox_virtual_environment_file_ubuntu_cloud_image_file_size" {
  value = "${proxmox_virtual_environment_file.ubuntu_cloud_image.file_size}"
}

output "resource_proxmox_virtual_environment_file_ubuntu_cloud_image_file_tag" {
  value = "${proxmox_virtual_environment_file.ubuntu_cloud_image.file_tag}"
}

output "resource_proxmox_virtual_environment_file_ubuntu_cloud_image_id" {
  value = "${proxmox_virtual_environment_file.ubuntu_cloud_image.id}"
}

output "resource_proxmox_virtual_environment_file_ubuntu_cloud_image_node_name" {
  value = "${proxmox_virtual_environment_file.ubuntu_cloud_image.node_name}"
}

output "resource_proxmox_virtual_environment_file_ubuntu_cloud_image_source" {
  value = "${proxmox_virtual_environment_file.ubuntu_cloud_image.source}"
}
