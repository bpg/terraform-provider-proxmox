resource "proxmox_ceph_pool" "vm" {
  node_name = "pve"
  name      = "vm"

  application = "rbd"
  size        = 3
  min_size    = 2
}

resource "proxmox_storage_rbd" "example" {
  id      = "ceph-vm"
  pool    = proxmox_ceph_pool.vm.name
  content = ["images", "rootdir"]
  krbd    = true
}

# A storage on an erasure-coded pool keeps image metadata in the metadata pool and image data in the data pool.
resource "proxmox_ceph_pool" "archive" {
  node_name = "pve"
  name      = "archive"

  erasure_coding = {
    k = 4
    m = 2
  }
}

resource "proxmox_storage_rbd" "erasure_coded" {
  id        = "ceph-archive"
  pool      = proxmox_ceph_pool.archive.metadata_pool_name
  data_pool = proxmox_ceph_pool.archive.data_pool.name
  content   = ["images"]
}
