resource "proxmox_ceph_pool" "example" {
  node_name = "pve"
  name      = "tank"

  application       = "rbd"
  size              = 3
  min_size          = 2
  pg_num            = 128
  pg_autoscale_mode = "warn"
}

# An erasure-coded pool: Proxmox creates the data pool "archive-data" and the metadata pool "archive-metadata".
resource "proxmox_ceph_pool" "erasure_coded" {
  node_name = "pve"
  name      = "archive"

  # With erasure_coding, size, min_size and crush_rule apply to the metadata pool.
  size       = 3
  min_size   = 2
  crush_rule = "replicated_rule"

  # The data pool.
  pg_num = 32

  erasure_coding = {
    k              = 4
    m              = 2
    failure_domain = "host"
    device_class   = "hdd"
  }
}
