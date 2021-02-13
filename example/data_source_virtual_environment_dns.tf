data "proxmox_virtual_environment_dns" "example" {
  node_name = data.proxmox_virtual_environment_nodes.example.names[0]
}

output "data_proxmox_virtual_environment_dns_example_domain" {
  value = data.proxmox_virtual_environment_dns.example.domain
}

output "data_proxmox_virtual_environment_dns_example_servers" {
  value = data.proxmox_virtual_environment_dns.example.servers
}
