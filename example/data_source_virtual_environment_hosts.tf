data "proxmox_virtual_environment_hosts" "example" {
  node_name = "${data.proxmox_virtual_environment_nodes.example.names[0]}"
}

output "data_proxmox_virtual_environment_hosts_example_addresses" {
  value = "${data.proxmox_virtual_environment_hosts.example.addresses}"
}

output "data_proxmox_virtual_environment_hosts_example_digest" {
  value = "${data.proxmox_virtual_environment_hosts.example.digest}"
}

output "data_proxmox_virtual_environment_hosts_example_hostnames" {
  value = "${data.proxmox_virtual_environment_hosts.example.hostnames}"
}
