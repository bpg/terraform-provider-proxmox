resource "proxmox_virtual_environment_cluster_options" "options" {
  language                  = "en"
  keyboard                  = "pl"
  email_from                = "ged@gont.earthsea"
  bandwidth_limit_migration = 555555
  bandwidth_limit_default   = 666666
  max_workers               = 5
  migration_cidr            = "10.0.0.0/8"
  migration_type            = "secure"
}
