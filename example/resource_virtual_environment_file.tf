#===============================================================================
# Cloud Config (cloud-init)
#===============================================================================

resource "proxmox_virtual_environment_file" "user_config" {
  content_type = "snippets"
  datastore_id = element(data.proxmox_virtual_environment_datastores.example.datastore_ids, index(data.proxmox_virtual_environment_datastores.example.datastore_ids, "local"))
  node_name    = data.proxmox_virtual_environment_datastores.example.node_name

  source_raw {
    data = <<EOF
#cloud-config
chpasswd:
  list: |
    ubuntu:example
  expire: false
hostname: terraform-provider-proxmox-example
users:
  - default
  - name: ubuntu
    groups: sudo
    shell: /bin/bash
    ssh-authorized-keys:
      - ${trimspace(tls_private_key.example.public_key_openssh)}
    sudo: ALL=(ALL) NOPASSWD:ALL
    EOF

    file_name = "terraform-provider-proxmox-example-user-config.yaml"
  }
}

resource "proxmox_virtual_environment_file" "vendor_config" {
  content_type = "snippets"
  datastore_id = element(data.proxmox_virtual_environment_datastores.example.datastore_ids, index(data.proxmox_virtual_environment_datastores.example.datastore_ids, "local"))
  node_name    = data.proxmox_virtual_environment_datastores.example.node_name

  source_raw {
    data = <<EOF
#cloud-config
runcmd:
    - apt update
    - apt install -y qemu-guest-agent
    - systemctl enable qemu-guest-agent
    - systemctl start qemu-guest-agent
    - echo "done" > /tmp/vendor-cloud-init-done
    EOF

    file_name = "terraform-provider-proxmox-example-vendor-config.yaml"
  }
}


#===============================================================================
# Ubuntu Cloud Image
#===============================================================================

resource "proxmox_virtual_environment_file" "ubuntu_cloud_image" {
  content_type = "iso"
  datastore_id = element(data.proxmox_virtual_environment_datastores.example.datastore_ids, index(data.proxmox_virtual_environment_datastores.example.datastore_ids, "local"))
  node_name    = data.proxmox_virtual_environment_datastores.example.node_name

  source_file {
    path = "https://cloud-images.ubuntu.com/bionic/current/bionic-server-cloudimg-amd64.img"
  }
}

output "resource_proxmox_virtual_environment_file_ubuntu_cloud_image_content_type" {
  value = proxmox_virtual_environment_file.ubuntu_cloud_image.content_type
}

output "resource_proxmox_virtual_environment_file_ubuntu_cloud_image_datastore_id" {
  value = proxmox_virtual_environment_file.ubuntu_cloud_image.datastore_id
}

output "resource_proxmox_virtual_environment_file_ubuntu_cloud_image_file_modification_date" {
  value = proxmox_virtual_environment_file.ubuntu_cloud_image.file_modification_date
}

output "resource_proxmox_virtual_environment_file_ubuntu_cloud_image_file_name" {
  value = proxmox_virtual_environment_file.ubuntu_cloud_image.file_name
}

output "resource_proxmox_virtual_environment_file_ubuntu_cloud_image_file_size" {
  value = proxmox_virtual_environment_file.ubuntu_cloud_image.file_size
}

output "resource_proxmox_virtual_environment_file_ubuntu_cloud_image_file_tag" {
  value = proxmox_virtual_environment_file.ubuntu_cloud_image.file_tag
}

output "resource_proxmox_virtual_environment_file_ubuntu_cloud_image_id" {
  value = proxmox_virtual_environment_file.ubuntu_cloud_image.id
}

output "resource_proxmox_virtual_environment_file_ubuntu_cloud_image_node_name" {
  value = proxmox_virtual_environment_file.ubuntu_cloud_image.node_name
}

output "resource_proxmox_virtual_environment_file_ubuntu_cloud_image_source_file" {
  value = proxmox_virtual_environment_file.ubuntu_cloud_image.source_file
}

#===============================================================================
# Ubuntu Container Template
#===============================================================================

resource "proxmox_virtual_environment_file" "ubuntu_container_template" {
  content_type = "vztmpl"
  datastore_id = element(data.proxmox_virtual_environment_datastores.example.datastore_ids, index(data.proxmox_virtual_environment_datastores.example.datastore_ids, "local"))
  node_name    = data.proxmox_virtual_environment_datastores.example.node_name

  source_file {
    path = "http://download.proxmox.com/images/system/ubuntu-18.04-standard_18.04.1-1_amd64.tar.gz"
  }
}

output "resource_proxmox_virtual_environment_file_ubuntu_container_template_content_type" {
  value = proxmox_virtual_environment_file.ubuntu_container_template.content_type
}

output "resource_proxmox_virtual_environment_file_ubuntu_container_template_datastore_id" {
  value = proxmox_virtual_environment_file.ubuntu_container_template.datastore_id
}

output "resource_proxmox_virtual_environment_file_ubuntu_container_template_file_modification_date" {
  value = proxmox_virtual_environment_file.ubuntu_container_template.file_modification_date
}

output "resource_proxmox_virtual_environment_file_ubuntu_container_template_file_name" {
  value = proxmox_virtual_environment_file.ubuntu_container_template.file_name
}

output "resource_proxmox_virtual_environment_file_ubuntu_container_template_file_size" {
  value = proxmox_virtual_environment_file.ubuntu_container_template.file_size
}

output "resource_proxmox_virtual_environment_file_ubuntu_container_template_file_tag" {
  value = proxmox_virtual_environment_file.ubuntu_container_template.file_tag
}

output "resource_proxmox_virtual_environment_file_ubuntu_container_template_id" {
  value = proxmox_virtual_environment_file.ubuntu_container_template.id
}

output "resource_proxmox_virtual_environment_file_ubuntu_container_template_node_name" {
  value = proxmox_virtual_environment_file.ubuntu_container_template.node_name
}

output "resource_proxmox_virtual_environment_file_ubuntu_container_template_source_file" {
  value = proxmox_virtual_environment_file.ubuntu_container_template.source_file
}
