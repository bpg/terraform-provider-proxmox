terraform {
  required_providers {
    local = {
      source  = "hashicorp/local"
      version = "2.5.2"
    }
    tls = {
      source  = "hashicorp/tls"
      version = "4.1.0"
    }
    proxmox = {
      source = "bpg/proxmox"
    }
  }
}
