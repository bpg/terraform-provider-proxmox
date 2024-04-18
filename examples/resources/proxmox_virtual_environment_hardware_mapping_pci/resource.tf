resource "proxmox_virtual_environment_hardware_mapping_pci" "example" {
  comment = "This is a comment"
  name    = "example"
  # The actual map of devices.
  map = [
    {
      comment = "This is a device specific comment"
      id      = "8086:5916"
      # This is an optional attribute, but causes a mapping to be incomplete when not defined.
      iommu_group = 0
      node        = "pve"
      path        = "0000:00:02.0"
      # This is an optional attribute, but causes a mapping to be incomplete when not defined.
      subsystem_id = "8086:2068"
    },
  ]
  mediated_devices = true
}
