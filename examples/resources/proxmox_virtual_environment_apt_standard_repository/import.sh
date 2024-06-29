#!/usr/bin/env sh
# An APT standard repository can be imported using a comma-separated list consisting of the name of the Proxmox VE node,
# and the standard repository handle in the exact same order, e.g.:
terraform import proxmox_virtual_environment_apt_standard_repository.example pve,no-subscription
