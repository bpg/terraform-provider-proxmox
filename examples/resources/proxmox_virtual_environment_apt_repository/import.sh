#!/usr/bin/env sh
# An APT repository can be imported using a comma-separated list consisting of the name of the Proxmox VE node,
# the absolute source list file path, and the index in the exact same order, e.g.:
terraform import proxmox_virtual_environment_apt_repository.example pve,/etc/apt/sources.list,0
