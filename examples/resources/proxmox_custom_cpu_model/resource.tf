resource "proxmox_custom_cpu_model" "vermeer" {
  id             = "RYZEN-Vermeer-v3"
  reported_model = "EPYC-Milan-v3"
  flags          = ["-pcid"]
  phys_bits      = "48"
}

# Reference the model from a VM through its prefixed CPU type:
#
# cpu {
#   type = proxmox_custom_cpu_model.vermeer.cputype
# }
