# List all Replications
data "proxmox_virtual_environment_replications" "all" {}

output "data_proxmox_virtual_environment_replications_all" {
  value = {
    replications = data.proxmox_virtual_environment_replications.all.replications
  }
}
