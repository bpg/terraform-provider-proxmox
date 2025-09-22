#!/usr/bin/env sh
# SDN subnet can be imported using its unique identifier (subnet ID)
# The subnet ID format is: <vnet>-<cidr> (e.g., "vnet1-192.168.1.0-24")
terraform import proxmox_virtual_environment_sdn_subnet.basic_subnet vnet1-192.168.1.0-24
terraform import proxmox_virtual_environment_sdn_subnet.dhcp_subnet vnet2-192.168.2.0-24
