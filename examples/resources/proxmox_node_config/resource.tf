resource "proxmox_node_config" "example" {
  node_name   = "pve"
  description = "Managed by Terraform"
}

resource "proxmox_node_config" "multiline" {
  node_name = "pve2"
  description = trimspace(<<-EOT
    Managed by Terraform.

    Owner: ops@example.com
    Role: storage
  EOT
  )
}
