resource "proxmox_virtual_environment_cluster_firewall" "cluster_options" {
  enabled = false

  ebtables      = false
  input_policy  = "ACCEPT"
  output_policy = "REJECT"
  log_ratelimit {
    enabled = false
    burst   = 20
    rate    = "5/second"
  }
}


resource "proxmox_virtual_environment_firewall_options" "vm_options" {
  depends_on = [proxmox_virtual_environment_vm.example]

  node_name = proxmox_virtual_environment_vm.example.node_name
  vm_id     = proxmox_virtual_environment_vm.example.vm_id

  dhcp          = true
  enabled       = false
  ipfilter      = true
  log_level_in  = "info"
  log_level_out = "info"
  macfilter     = false
  ndp           = true
  input_policy  = "REJECT"
  output_policy = "REJECT"
  radv          = true
}


resource "proxmox_virtual_environment_firewall_options" "container_options" {
  depends_on = [proxmox_virtual_environment_container.example]

  node_name    = proxmox_virtual_environment_container.example.node_name
  container_id = proxmox_virtual_environment_container.example.vm_id

  dhcp          = false
  enabled       = false
  ipfilter      = false
  log_level_in  = "alert"
  log_level_out = "alert"
  macfilter     = true
  ndp           = false
  input_policy  = "ACCEPT"
  output_policy = "DROP"
  radv          = false
}

output "resource_proxmox_virtual_environment_firewall_options_cluster" {
  value = proxmox_virtual_environment_cluster_firewall.cluster_options
}

output "resource_proxmox_virtual_environment_firewall_options_vm" {
  value = proxmox_virtual_environment_firewall_options.vm_options
}

output "resource_proxmox_virtual_environment_firewall_options_container" {
  value = proxmox_virtual_environment_firewall_options.container_options
}
