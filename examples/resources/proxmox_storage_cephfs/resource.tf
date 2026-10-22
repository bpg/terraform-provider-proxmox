resource "proxmox_storage_cephfs" "example" {
  id      = "cephfs-templates"
  fs_name = "cephfs"

  content = ["backup", "iso", "vztmpl", "snippets"]
  options = "noatime"

  backups {
    keep_last = 3
  }
}
