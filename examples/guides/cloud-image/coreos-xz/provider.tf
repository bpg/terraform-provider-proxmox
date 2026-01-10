terraform {
  required_providers {
    proxmox = {
      source  = "bpg/proxmox"
      version = "0.92.0" # x-release-please-version
    }
    http = {
      source  = "hashicorp/http"
      version = "~> 3.0"
    }
  }
}

provider "proxmox" {
  endpoint  = var.virtual_environment_endpoint
  api_token = var.virtual_environment_token
  # SSH configuration is required for file_id disk imports
  ssh {
    agent    = true
    username = "root"
  }
}

