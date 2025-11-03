resource "proxmox_virtual_environment_vm" "test_vm1" {
  vm_id     = 1234
  node_name = "pve"
  started   = false
}

resource "proxmox_virtual_environment_pool" "test_pool" {
  pool_id = "test-pool"
}

resource "proxmox_virtual_environment_pool_membership" "vm_membership" {
  pool_id = proxmox_virtual_environment_pool.test_pool.id
  vm_id   = proxmox_virtual_environment_vm.test_vm1.id
}

resource "proxmox_virtual_environment_pool_membership" "storage_membership" {
  pool_id    = proxmox_virtual_environment_pool.test_pool.id
  storage_id = "local-lvm"
}