resource "proxmox_virtual_environment_hardware_mapping_usb" "example" {
  comment = "This is a comment"
  name    = "example"
  # The actual map of devices.
  map = [
    {
      comment = "This is a device specific comment"
      id      = "8087:0a2b"
      node    = "pve"
      # This attribute is optional, but can be used to map the device based on its port instead of only the device ID.
      path = "1-8.2"
    },
  ]
}
