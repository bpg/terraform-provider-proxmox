resource "proxmox_virtual_environment_cluster_ipset" "example" {
	name    = "local_network"
	comment = "Managed by Terraform"

    cidr {
        name = "192.168.0.0/23"
        comment = "Local network 1"
    }

    cidr {
        name = "192.168.0.1"
        comment = "Server 1"
        nomatch = true
    }

    cidr {
        name = "192.168.2.1"
        comment = "Server 1"
    }
}

output "resource_proxmox_virtual_environment_cluster_ipset" {
  value = "${proxmox_virtual_environment_cluster_ipset.example.name}"
}

