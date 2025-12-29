#!/usr/bin/env sh
# SDN vnet can be imported using its unique identifier (vnet ID)
terraform import proxmox_virtual_environment_sdn_vnet.basic_vnet vnet1
terraform import proxmox_virtual_environment_sdn_vnet.isolated_vnet vnet2
