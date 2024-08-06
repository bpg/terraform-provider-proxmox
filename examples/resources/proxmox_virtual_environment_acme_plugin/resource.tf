resource "proxmox_virtual_environment_acme_plugin" "example" {
  type   = "dns"
  plugin = "test"
  api    = "aws"
  data = {
    AWS_ACCESS_KEY_ID     = "EXAMPLE"
    AWS_SECRET_ACCESS_KEY = "EXAMPLE"
  }
}
