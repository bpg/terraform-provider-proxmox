resource "proxmox_sdn_zone_simple" "application" {
  id = "appzone"

  depends_on = [
    proxmox_sdn_applier.finalizer,
  ]
}

resource "proxmox_sdn_vnet" "application" {
  id   = "appnet"
  zone = proxmox_sdn_zone_simple.application.id
}

resource "proxmox_sdn_applier" "application" {
  depends_on = [
    proxmox_sdn_vnet.application,
  ]
}

resource "proxmox_sdn_applier" "finalizer" {}

resource "proxmox_sdn_vnet_firewall_rules" "application" {
  vnet = proxmox_sdn_vnet.application.id

  depends_on = [
    proxmox_sdn_applier.application,
  ]

  rules = [
    {
      action  = "ACCEPT"
      comment = "Allow HTTPS from the service tier"
      source  = "192.0.2.0/24"
      dest    = "198.51.100.0/24"
      proto   = "tcp"
      dport   = "443"
    },
    {
      action  = "DROP"
      comment = "Deny all other service-tier traffic"
      source  = "192.0.2.0/24"
      dest    = "198.51.100.0/24"
      log     = "info"
    },
  ]
}
