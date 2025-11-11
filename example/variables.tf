variable "virtual_environment_endpoint" {
  description = "The endpoint for the Proxmox Virtual Environment API (example: https://host:port)"
  type        = string
}

variable "virtual_environment_username" {
  description = "Proxmox User for API Access"
  type        = string
  default     = "root@pam"
}

variable "virtual_environment_password" {
  description = "Password for Proxmox API User"
  type        = string
  sensitive   = true
  default     = "don not use default passwords!"
}

variable "virtual_environment_api_token" {
  description = "The API token for the Proxmox Virtual Environment API"
  type        = string
  sensitive   = true
}

variable "virtual_environment_ssh_username" {
  description = "The username for the Proxmox Virtual Environment API"
  type        = string
  default     = "root"
}

variable "virtual_environment_node_name" {
  description = "Name of the Proxmox node"
  type        = string
  default     = "pve"
}

variable "virtual_environment_insecure" {
  description = "Self Signed Certificates Used"
  type        = bool
  default     = true
}

variable "virtual_environment_storage" {
  description = "Name of the Proxmox storage"
  type        = string
  default     = "local-lvm"
}

variable "latest_debian_12_bookworm_qcow2_img_url" {
  description = "The URL for the latest Debian 12 Bookworm qcow2 image"
  type        = string
  default     = "https://cloud.debian.org/images/cloud/bookworm/latest/debian-12-generic-amd64.qcow2"
}

variable "ubuntu_24_04_lxc_img_url" {
  description = "The URL for the Ubuntu 24.10 LXC image"
  type        = string
  default     = "https://cloud-images.ubuntu.com/minimal/releases/noble/release/ubuntu-24.04-minimal-cloudimg-amd64-root.tar.xz"
}

variable "ubuntu_24_04_lxc_img_checksum" {
  description = "The checksum for the Ubuntu 24.10 LXC image"
  type        = string
  default     = "bdc90f9c0908e45dc602208da8fbd9e38db6483a1483fc7042332f5cd4e9aa22"
}
