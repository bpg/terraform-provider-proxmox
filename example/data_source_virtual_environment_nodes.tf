data "proxmox_virtual_environment_nodes" "example" {}

output "data_proxmox_virtual_environment_nodes_example_cpu_count" {
  value = "${data.proxmox_virtual_environment_nodes.example.cpu_count}"
}

output "data_proxmox_virtual_environment_nodes_example_cpu_utilization" {
  value = "${data.proxmox_virtual_environment_nodes.example.cpu_utilization}"
}

output "data_proxmox_virtual_environment_nodes_example_memory_available" {
  value = "${data.proxmox_virtual_environment_nodes.example.memory_available}"
}

output "data_proxmox_virtual_environment_nodes_example_memory_used" {
  value = "${data.proxmox_virtual_environment_nodes.example.memory_used}"
}

output "data_proxmox_virtual_environment_nodes_example_names" {
  value = "${data.proxmox_virtual_environment_nodes.example.names}"
}

output "data_proxmox_virtual_environment_nodes_example_online" {
  value = "${data.proxmox_virtual_environment_nodes.example.online}"
}

output "data_proxmox_virtual_environment_nodes_example_ssl_fingerprints" {
  value = "${data.proxmox_virtual_environment_nodes.example.ssl_fingerprints}"
}

output "data_proxmox_virtual_environment_nodes_example_support_levels" {
  value = "${data.proxmox_virtual_environment_nodes.example.support_levels}"
}

output "data_proxmox_virtual_environment_nodes_example_uptime" {
  value = "${data.proxmox_virtual_environment_nodes.example.uptime}"
}
