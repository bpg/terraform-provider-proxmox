data "proxmox_virtual_environment_pools" "example" {
  depends_on = ["proxmox_virtual_environment_pool.example"]
}

output "data_proxmox_virtual_environment_pools_example" {
  value = "${map(
    "pool_ids", data.proxmox_virtual_environment_pools.example.pool_ids,
  )}"
}
