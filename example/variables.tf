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