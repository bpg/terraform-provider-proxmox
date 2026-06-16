#!/usr/bin/env sh
# A download file can be imported using its identifier in the format: node_name/datastore_id:content_type/file_name, e.g.:
terraform import proxmox_virtual_environment_download_file.ubuntu_iso pve/local:iso/ubuntu-24.04-server.iso
