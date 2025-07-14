variable "virtual_environment_endpoint" {
  type        = string
  description = "The endpoint for the Proxmox Virtual Environment API (example: https://host:port)"
}

variable "virtual_environment_api_token" {
  type        = string
  description = "The API token for the Proxmox Virtual Environment API"
}

variable "virtual_environment_ssh_username" {
  type        = string
  description = "The username for the Proxmox Virtual Environment API"
}

variable "virtual_environment_node_name" {
  description = "Name of the Proxmox node"
  type        = string
  default = "pve"
}

variable "virtual_environment_storage" {
  description = "Name of the Proxmox storage"
  type        = string
  default = "local-lvm"
}

variable "latest_debian_12_bookworm_qcow2_img_url" {
  type        = string
  description = "The URL for the latest Debian 12 Bookworm qcow2 image"
  default     = "https://cloud.debian.org/images/cloud/bookworm/latest/debian-12-generic-amd64.qcow2"
}

variable "release_20250610_ubuntu_24_noble_lxc_img_url" {
  type        = string
  description = "The URL for the Ubuntu 24.04 LXC image"
  default     = "https://mirrors.servercentral.com/ubuntu-cloud-images/releases/24.04/release-20250610/ubuntu-24.04-server-cloudimg-amd64-root.tar.xz"
}

variable "release_20250610_ubuntu_24_noble_lxc_img_checksum" {
  type        = string
  description = "The checksum for the Ubuntu 24.04 LXC image"
  default     = "ae1fc4b5f020e6f1f2048beb5a7635f7bce4d72723239b7dea86af062cc1ab79"
}
