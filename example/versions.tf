terraform {
  required_providers {
    local = {
      source  = "hashicorp/local"
      version = "2.9.0"
    }
    tls = {
      source  = "hashicorp/tls"
      version = "4.3.0"
    }
    proxmox = {
      source = "bpg/proxmox"
    }
  }
}
