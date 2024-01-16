terraform {
  required_providers {
    local = {
      source  = "hashicorp/local"
      version = "2.4.1"
    }
    tls = {
      source  = "hashicorp/tls"
      version = "3.4.0"
    }
    proxmox = {
      source = "bpg/proxmox"
    }
  }
}
