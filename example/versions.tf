terraform {
  required_providers {
    local = {
      source = "hashicorp/local"
    }
    proxmox = {
      source  = "terraform.danitso.com/provider/proxmox"
    }
    tls = {
      source = "hashicorp/tls"
    }
  }
  required_version = ">= 0.13"
}
