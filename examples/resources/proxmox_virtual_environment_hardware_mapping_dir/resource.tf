resource "proxmox_virtual_environment_hardware_mapping_dir" "example" {
  comment = "This is a comment"
  name    = "example"
  # The actual map of devices.
  map = [
    {
      node = "pve"
      path = "/mnt/data"
    },
  ]
}
