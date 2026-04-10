# List all Replications
data "proxmox_replications" "all" {}

output "data_proxmox_replications_all" {
  value = {
    replications = data.proxmox_replications.all.replications
  }
}
