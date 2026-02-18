terraform {
  required_providers {
    local = {
      source  = "hashicorp/local"
      version = "2.7.0"
    }
    tls = {
      source  = "hashicorp/tls"
      version = "4.2.1"
    }
    proxmox = {
      source = "bpg/proxmox"
    }
  }
}
