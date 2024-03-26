resource "proxmox_virtual_environment_user" "operations_automation" {
  comment  = "Managed by Terraform"
  password = "a-strong-password"
  user_id  = "operations-automation@pve"
}

resource "proxmox_virtual_environment_role" "operations_monitoring" {
  role_id = "operations-monitoring"

  privileges = [
    "VM.Monitor",
  ]
}

resource "proxmox_virtual_environment_acl" "operations_automation_monitoring" {
  user_id = proxmox_virtual_environment_user.operations_automation.user_id
  role_id = proxmox_virtual_environment_role.operations_monitoring.role_id

  path      = "/vms/1234"
  propagate = true
}
