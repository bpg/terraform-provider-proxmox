resource "proxmox_virtual_environment_hosts" "example" {
  node_name = "${data.proxmox_virtual_environment_nodes.example.names[0]}"

  dynamic "entry" {
    for_each = "${data.proxmox_virtual_environment_hosts.example.entries}"

    content {
      address   = "${entry.value.address}"
      hostnames = "${entry.value.hostnames}"
    }
  }
}

output "resource_proxmox_virtual_environment_hosts_example_addresses" {
  value = "${proxmox_virtual_environment_hosts.example.addresses}"
}

output "resource_proxmox_virtual_environment_hosts_example_digest" {
  value = "${proxmox_virtual_environment_hosts.example.digest}"
}

output "resource_proxmox_virtual_environment_hosts_example_entries" {
  value = "${proxmox_virtual_environment_hosts.example.entries}"
}

output "resource_proxmox_virtual_environment_hosts_example_hostnames" {
  value = "${proxmox_virtual_environment_hosts.example.hostnames}"
}
