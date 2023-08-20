// This will fetch the set of all HA resource identifiers.
data "proxmox_virtual_environment_haresources" "example_all" {}

// This will fetch the set of HA resource identifiers that correspond to virtual machines.
data "proxmox_virtual_environment_haresources" "example_vm" {
  type = "vm"
}

output "data_proxmox_virtual_environment_haresources" {
  value = {
    all = data.proxmox_virtual_environment_haresources.example_all.resource_ids
    vms = data.proxmox_virtual_environment_haresources.example_vm.resource_ids
  }
}
