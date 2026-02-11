data "proxmox_virtual_environment_replication" "example" {
  id = "100-0"
}

output "data_proxmox_virtual_environment_replication" {
  value = {
    id     = data.proxmox_virtual_environment_replication.example.id
    target = data.proxmox_virtual_environment_replication.example.target
    type   = data.proxmox_virtual_environment_replication.example.type
    jobnum = data.proxmox_virtual_environment_replication.example.jobnum
    guest  = data.proxmox_virtual_environment_replication.example.guest
  }
}
