terraform {
  required_providers {
    local = {
      source  = "hashicorp/local"
      version = "2.2.2"
    }
    tls = {
      source  = "hashicorp/tls"
      version = "3.1.0"
    }
    proxmox = {
      source  = "bpg/proxmox"
    }
  }
}
