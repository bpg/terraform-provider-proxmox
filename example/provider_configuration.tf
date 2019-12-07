provider "proxmox" {
  virtual_environment {
    endpoint = "${var.virtual_environment_endpoint}"
    username = "${var.virtual_environment_username}"
    password = "${var.virtual_environment_password}"
    insecure = true
  }
}

variable "virtual_environment_endpoint" {
  type        = "string"
  description = "The endpoint for the Proxmox Virtual Environment API (example: https://host:port)"
}

variable "virtual_environment_password" {
  type        = "string"
  description = "The password for the Proxmox Virtual Environment API"
}

variable "virtual_environment_username" {
  type        = "string"
  description = "The username and realm for the Proxmox Virtual Environment API (example: root@pam)"
}
