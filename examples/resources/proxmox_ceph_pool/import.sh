#!/usr/bin/env sh
# Ceph pools can be imported using the format `node_name/pool_name`, e.g.:
terraform import proxmox_ceph_pool.example pve/tank

# An erasure-coded pool is imported by the name it was created with (without `-data` or `-metadata`); its
# configuration must then include a matching `erasure_coding` attribute.
terraform import proxmox_ceph_pool.erasure_coded pve/archive
