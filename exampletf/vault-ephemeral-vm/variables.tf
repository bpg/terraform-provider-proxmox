variable "vault_address" {
  description = "Address of the Vault server (e.g. https://vault.example.com:8200)"
  type        = string
}

variable "vault_proxmox_token_path" {
  description = "Vault KV v2 path for the Proxmox API token. The secret must have key 'api_token'."
  type        = string
  default     = "proxmox/api-token"
}

variable "vault_vm_credentials_path" {
  description = "Vault KV v2 path for the VM cloud-init credentials. The secret must have key 'password'."
  type        = string
  default     = "proxmox/vm-credentials"
}

variable "vault_kv_mount" {
  description = "Vault KV v2 secret engine mount path."
  type        = string
  default     = "secret"
}

variable "proxmox_endpoint" {
  description = "Proxmox VE API endpoint (e.g. https://pve.example.com:8006)"
  type        = string
}

variable "proxmox_insecure" {
  description = "Skip TLS certificate verification for the Proxmox API."
  type        = bool
  default     = false
}

variable "proxmox_node" {
  description = "Name of the Proxmox node to create the VM on."
  type        = string
}

variable "vm_id" {
  description = "VMID to assign to the cloned VM."
  type        = number
}

variable "template_vm_id" {
  description = "VMID of the source cloud-init template VM to clone from."
  type        = number
}

variable "vm_name" {
  description = "Name for the VM."
  type        = string
  default     = "vault-provisioned-vm"
}

variable "vm_username" {
  description = "Cloud-init default username. Not sensitive — stored in Terraform state."
  type        = string
  default     = "ubuntu"
}

variable "vm_cores" {
  description = "Number of CPU cores."
  type        = number
  default     = 2
}

variable "vm_memory_mb" {
  description = "Memory in MiB."
  type        = number
  default     = 2048
}

variable "vm_network_bridge" {
  description = "Linux bridge to attach the VM's network interface to (e.g. vmbr0)."
  type        = string
  default     = "vmbr0"
}
