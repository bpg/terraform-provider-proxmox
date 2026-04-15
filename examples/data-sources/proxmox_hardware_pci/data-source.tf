# List all PCI devices on a node (using default blacklist)
data "proxmox_hardware_pci" "example" {
  node_name = "pve"
}

# List all PCI devices including bridges and memory controllers
data "proxmox_hardware_pci" "all" {
  node_name           = "pve"
  pci_class_blacklist = []
}

# Find all NVIDIA GPUs (vendor ID 10de = NVIDIA, class 03 = display controller)
data "proxmox_hardware_pci" "gpus" {
  node_name           = "pve"
  pci_class_blacklist = []

  filters = {
    vendor_id = "10de"
    class     = "03"
  }
}
