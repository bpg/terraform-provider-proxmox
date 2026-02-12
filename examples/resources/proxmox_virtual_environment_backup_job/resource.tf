# Simple daily backup job for specific VMs
resource "proxmox_virtual_environment_backup_job" "daily_vm_backup" {
  id       = "daily-vm-backup"
  schedule = "02:00"
  storage  = "local"
  vmid     = "100,101,102"
  mode     = "snapshot"
  compress = "zstd"
  mailto   = "admin@example.com"
}

# Weekly backup of all VMs in a pool with retention policy
resource "proxmox_virtual_environment_backup_job" "weekly_pool_backup" {
  id            = "weekly-pool-backup"
  schedule      = "sun 03:00"
  storage       = "pbs"
  pool          = "production"
  mode          = "snapshot"
  compress      = "zstd"
  prune_backups = "keep-last=3,keep-weekly=4,keep-monthly=6"
  remove        = true
  enabled       = true
}

# Backup all VMs on a specific node
resource "proxmox_virtual_environment_backup_job" "node_backup" {
  id       = "node-backup"
  schedule = "01:00"
  storage  = "local"
  node     = "pve"
  all      = true
  mode     = "stop"
  compress = "gzip"
}
