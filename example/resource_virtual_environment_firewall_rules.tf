resource "proxmox_virtual_environment_firewall_rules" "cluster_rules" {
  rule {
    type    = "in"
    action  = "ACCEPT"
    comment = "Allow FTP"
    dest    = "192.168.0.5"
    dport   = "21"
    proto   = "tcp"
    log     = "info"
  }

  rule {
    type    = "out"
    action  = "DROP"
    comment = "Allow SSH"
    dest    = "192.168.0.5"
    dport   = "22"
    proto   = "tcp"
  }
}

resource "proxmox_virtual_environment_firewall_rules" "vm_rules" {
  depends_on = [proxmox_virtual_environment_vm.example]

  node_name = proxmox_virtual_environment_vm.example.node_name
  vm_id     = proxmox_virtual_environment_vm.example.vm_id

  rule {
    type    = "in"
    action  = "ACCEPT"
    comment = "Allow FTP"
    dest    = "192.168.1.5"
    dport   = "21"
    proto   = "tcp"
    log     = "info"
  }

  rule {
    type    = "out"
    action  = "DROP"
    comment = "Allow SSH"
    dest    = "192.168.1.5"
    dport   = "22"
    proto   = "tcp"
  }
}

resource "proxmox_virtual_environment_firewall_rules" "container_rules" {
  depends_on = [proxmox_virtual_environment_container.example]

  node_name = proxmox_virtual_environment_container.example.node_name
  container_id    = proxmox_virtual_environment_container.example.vm_id

  rule {
    type    = "in"
    action  = "ACCEPT"
    comment = "Allow FTP"
    dest    = "192.168.2.5"
    dport   = "21"
    proto   = "tcp"
    log     = "info"
  }

  rule {
    type    = "out"
    action  = "DROP"
    comment = "Allow SSH"
    dest    = "192.168.2.5"
    dport   = "22"
    proto   = "tcp"
  }
}

output "resource_proxmox_virtual_environment_firewall_rules_cluster" {
  value = proxmox_virtual_environment_firewall_rules.cluster_rules
}

output "resource_proxmox_virtual_environment_firewall_rules_vm" {
  value = proxmox_virtual_environment_firewall_rules.vm_rules
}

output "resource_proxmox_virtual_environment_firewall_rules_container" {
  value = proxmox_virtual_environment_firewall_rules.container_rules
}
