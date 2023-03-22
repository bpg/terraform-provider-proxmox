data "proxmox_virtual_environment_vms" "example" {
  depends_on = [proxmox_virtual_environment_vm.example]
  tags = ["ubuntu"]

  lifecycle {
    postcondition {
      condition     = length(self.vms) == 1
      error_message = "Only 1 vm should have this tag"
    }
  }
}

output "proxmox_virtual_environment_vms_example" {
  value = data.proxmox_virtual_environment_vms.example.vms
}
