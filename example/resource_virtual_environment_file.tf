resource "proxmox_virtual_environment_file" "ubuntu_cloud_image" {
  content_type = "iso"
  datastore_id = "${element(data.proxmox_virtual_environment_datastores.example.datastore_ids, index(data.proxmox_virtual_environment_datastores.example.datastore_ids, "local"))}"
  node_name    = "${data.proxmox_virtual_environment_datastores.example.node_name}"

  source_file {
    path = "https://cloud-images.ubuntu.com/bionic/current/bionic-server-cloudimg-amd64.img"
  }
}

resource "proxmox_virtual_environment_file" "cloud_init_config" {
  content_type = "snippets"
  datastore_id = "${element(data.proxmox_virtual_environment_datastores.example.datastore_ids, index(data.proxmox_virtual_environment_datastores.example.datastore_ids, "local"))}"
  node_name    = "${data.proxmox_virtual_environment_datastores.example.node_name}"

  source_raw {
    data = <<EOF
#cloud-config
chpasswd:
  list: |
    ubuntu:example
  expire: False
hostname: terraform-provider-proxmox-example
packages:
  - qemu-guest-agent
users:
  - default
  - name: ubuntu
    groups: sudo
    shell: /bin/bash
    ssh-authorized-keys:
      - ${trimspace(tls_private_key.example.public_key_openssh)}
    sudo: ALL=(ALL) NOPASSWD:ALL
    EOF

    file_name = "terraform-provider-proxmox-example-cloud-init.yaml"
  }
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

output "resource_proxmox_virtual_environment_file_ubuntu_cloud_image_source_file" {
  value = "${proxmox_virtual_environment_file.ubuntu_cloud_image.source_file}"
}
