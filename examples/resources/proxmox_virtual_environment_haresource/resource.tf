resource "proxmox_virtual_environment_haresource" "example" {
  depends_on = [
    proxmox_virtual_environment_hagroup.example
  ]
  resource_id = "vm:123"
  state       = "started"
  group       = "example"
  comment     = "Managed by Terraform"
}
