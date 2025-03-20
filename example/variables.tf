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

variable "release_20240725_ubuntu_24_noble_lxc_img_url" {
  type        = string
  description = "The URL for the Ubuntu 24.04 LXC image"
  default     = "https://mirrors.servercentral.com/ubuntu-cloud-images/releases/24.04/release-20240725/ubuntu-24.04-server-cloudimg-amd64-root.tar.xz"
}

variable "release_20240725_ubuntu_24_noble_lxc_img_checksum" {
  type        = string
  description = "The checksum for the Ubuntu 24.04 LXC image"
  default     = "10331782a01cd2348b421a261f0e15ba041358bd540f66f2432b162e70b90d28"
}
