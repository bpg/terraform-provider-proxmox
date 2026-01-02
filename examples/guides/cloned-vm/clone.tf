# Using the new cloned_vm resource (recommended)
# This resource provides explicit opt-in management of devices and configuration
resource "proxmox_virtual_environment_cloned_vm" "ubuntu_clone" {
  node_name = var.virtual_environment_node_name
  name      = "ubuntu-clone"

  clone = {
    source_vm_id = proxmox_virtual_environment_vm.ubuntu_template.vm_id
    full         = true
  }

  # Only explicitly listed devices are managed
  # Network device inherited from template is preserved but not managed
  # To manage it, explicitly list it here:
  network = {
    net0 = {
      bridge = "vmbr0"
      model  = "virtio"
    }
  }

  # Memory configuration using new terminology
  memory = {
    size    = 2048
    balloon = 512
  }

  cpu = {
    cores = 2
  }
}

# NOTE: Initialization and agent settings are inherited from the template
# The template already has cloud-init configured (see cloud-config.tf)
# If you need to customize initialization, use the legacy VM resource with clone block
# See clone-legacy.tf for an example using proxmox_virtual_environment_vm

output "vm_id" {
  value = proxmox_virtual_environment_cloned_vm.ubuntu_clone.id
}

