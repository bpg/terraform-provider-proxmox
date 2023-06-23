resource "proxmox_virtual_environment_firewall_ipset" "cluster_ipset" {
  name    = "cluster-ipset"
  comment = "Managed by Terraform"

  cidr {
    name    = "192.168.0.1"
    comment = "Server 1"
    nomatch = true
  }

  cidr {
    name    = "192.168.0.2"
    comment = "Server 2"
  }
}

resource "proxmox_virtual_environment_firewall_ipset" "vm_ipset" {
  depends_on = [proxmox_virtual_environment_vm.example]

  node_name = proxmox_virtual_environment_vm.example.node_name
  vm_id     = proxmox_virtual_environment_vm.example.vm_id

  name    = "vm-ipset"
  comment = "Managed by Terraform"

  cidr {
    name    = "192.168.1.1"
    comment = "Server 1"
    nomatch = true
  }

  cidr {
    name    = "192.168.1.2"
    comment = "Server 2"
  }
}

resource "proxmox_virtual_environment_firewall_ipset" "container_ipset" {
  depends_on = [proxmox_virtual_environment_container.example]

  node_name    = proxmox_virtual_environment_container.example.node_name
  container_id = proxmox_virtual_environment_container.example.vm_id

  name    = "container-ipset"
  comment = "Managed by Terraform"

  cidr {
    name    = "192.168.2.1"
    comment = "Server 1"
    nomatch = true
  }

  cidr {
    name    = "192.168.2.2"
    comment = "Server 2"
  }
}


output "resource_proxmox_virtual_environment_firewall_ipset_cluster" {
  value = proxmox_virtual_environment_firewall_ipset.cluster_ipset
}

output "resource_proxmox_virtual_environment_firewall_ipset_vm" {
  value = proxmox_virtual_environment_firewall_ipset.vm_ipset
}

output "resource_proxmox_virtual_environment_firewall_ipset_container" {
  value = proxmox_virtual_environment_firewall_ipset.container_ipset
}
