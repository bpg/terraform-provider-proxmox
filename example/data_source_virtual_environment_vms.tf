data "proxmox_virtual_environment_vms" "example" {
  depends_on = [proxmox_virtual_environment_vm.example]
  tags       = ["ubuntu"]

  lifecycle {
    postcondition {
      condition     = length(self.vms) == 1
      error_message = "Only 1 vm should have this tag"
    }
  }
}

data "proxmox_virtual_environment_vms" "template_example" {
  depends_on = [proxmox_virtual_environment_vm.example]
  tags       = ["ubuntu"]

  filter {
    name   = "template"
    values = [false]
  }

  filter {
    name   = "status"
    values = ["running"]
  }

  filter {
    name   = "name"
    regex  = true
    values = [".*ubuntu.*"]
  }
}

output "proxmox_virtual_environment_vms_example" {
  value = data.proxmox_virtual_environment_vms.example.vms
}

output "proxmox_virtual_environment_template_vms_example" {
  value = data.proxmox_virtual_environment_vms.template_example.vms
}
