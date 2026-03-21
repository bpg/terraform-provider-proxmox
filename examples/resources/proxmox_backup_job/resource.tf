resource "proxmox_backup_job" "daily_backup" {
  id       = "daily-backup"
  schedule = "*-*-* 02:00"
  storage  = "local"
  all      = true
  mode     = "snapshot"
  compress = "zstd"
}
