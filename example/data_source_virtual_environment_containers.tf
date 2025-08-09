data "proxmox_virtual_environment_containers" "example" {
  depends_on = [proxmox_virtual_environment_container.example]
  tags       = ["ubuntu"]

  lifecycle {
    postcondition {
      condition     = length(self.containers) == 1
      error_message = "Only 1 container should have this tag"
    }
  }
}

data "proxmox_virtual_environment_containers" "template_example" {
  depends_on = [proxmox_virtual_environment_container.example]
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

output "proxmox_virtual_environment_containers_example" {
  value = data.proxmox_virtual_environment_containers.example.containers
}

output "proxmox_virtual_environment_template_containers_example" {
  value = data.proxmox_virtual_environment_containers.template_example.containers
}