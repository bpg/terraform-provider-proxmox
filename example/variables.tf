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

variable "latest_debian_12_bookworm_qcow2_img_url" {
  type        = string
  description = "The URL for the latest Debian 12 Bookworm qcow2 image"
  default     = "https://cloud.debian.org/images/cloud/bookworm/latest/debian-12-generic-amd64.qcow2"
}

variable "release_20240416_ubuntu_22_jammy_lxc_img_url" {
  type        = string
  description = "The URL for the Ubuntu 22.04 LXC image"
  default     = "https://cloud-images.ubuntu.com/releases/22.04/release-20240416/ubuntu-22.04-server-cloudimg-amd64-root.tar.xz"
}

variable "release_20240416_ubuntu_22_jammy_lxc_img_checksum" {
  type        = string
  description = "The checksum for the Ubuntu 22.04 LXC image"
  default     = "a362bf415ad2eae4854deda6237894a96178a2edbbd5a1956d6c55c5837a80d3"
}
