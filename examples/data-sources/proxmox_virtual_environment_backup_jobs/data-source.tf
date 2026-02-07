data "proxmox_virtual_environment_backup_jobs" "all" {}

output "backup_job_ids" {
  value = [for job in data.proxmox_virtual_environment_backup_jobs.all.jobs : job.id]
}

output "enabled_backup_jobs" {
  value = [for job in data.proxmox_virtual_environment_backup_jobs.all.jobs : job.id if job.enabled]
}
