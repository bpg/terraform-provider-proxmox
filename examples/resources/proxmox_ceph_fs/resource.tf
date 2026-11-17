resource "proxmox_ceph_fs" "shared" {
  node_name = "pve"
  name      = "shared"
  pg_num    = 32

  metadata_pool = {
    size     = 3
    min_size = 2
  }

  data_pool = {
    crush_rule = "replicated_hdd"
    size       = 3
    min_size   = 2
  }
}

resource "proxmox_storage_cephfs" "shared" {
  id      = "shared"
  fs_name = proxmox_ceph_fs.shared.name
  content = ["backup", "iso", "vztmpl", "snippets"]
}
