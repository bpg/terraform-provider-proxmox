terraform {
  required_providers {
    local = {
      source = "hashicorp/local"
    }
    proxmox = {
      source  = "danitso/proxmox"
    }
    tls = {
      source = "hashicorp/tls"
    }
  }
  required_version = ">= 0.13"
}
