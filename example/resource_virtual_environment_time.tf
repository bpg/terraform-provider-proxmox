resource "proxmox_virtual_environment_time" "example" {
  node_name = "${data.proxmox_virtual_environment_time.example.node_name}"
  time_zone = "${data.proxmox_virtual_environment_time.example.time_zone}"
}

output "resource_proxmox_virtual_environment_time" {
  value = "${map(
    "local_time", data.proxmox_virtual_environment_time.example.local_time,
    "time_zone", data.proxmox_virtual_environment_time.example.time_zone,
    "utc_time", data.proxmox_virtual_environment_time.example.utc_time,
  )}"
}
