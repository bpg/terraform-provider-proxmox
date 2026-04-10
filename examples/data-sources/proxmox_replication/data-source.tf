data "proxmox_replication" "example" {
  id = "100-0"
}

output "data_proxmox_replication" {
  value = {
    id     = data.proxmox_replication.example.id
    target = data.proxmox_replication.example.target
    type   = data.proxmox_replication.example.type
    jobnum = data.proxmox_replication.example.jobnum
    guest  = data.proxmox_replication.example.guest
  }
}
