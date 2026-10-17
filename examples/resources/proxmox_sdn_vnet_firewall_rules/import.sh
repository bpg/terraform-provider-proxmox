#!/usr/bin/env sh
# The complete firewall ruleset is imported using its VNet identifier.
terraform import proxmox_sdn_vnet_firewall_rules.application appnet
