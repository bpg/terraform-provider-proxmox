#!/usr/bin/env sh
# SDN subnet can be imported using its unique identifier in the format: <vnet>/<subnet-id>
# The <subnet-id> is the canonical ID from Proxmox, e.g., "zone1-192.168.1.0-24"
terraform import proxmox_virtual_environment_sdn_subnet.basic_subnet vnet1/zone1-192.168.1.0-24
terraform import proxmox_virtual_environment_sdn_subnet.dhcp_subnet vnet2/zone2-192.168.2.0-24
