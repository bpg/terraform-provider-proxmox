terraform {
  required_version = ">= 1.11"

  required_providers {
    vault = {
      source  = "hashicorp/vault"
      version = "~> 5.9"
    }
    proxmox = {
      source  = "bpg/proxmox"
      version = ">= 0.110"
    }
  }
}

# ---------------------------------------------------------------------------
# Vault provider
# Set VAULT_TOKEN (or VAULT_ROLE_ID + VAULT_SECRET_ID) in your environment.
# ---------------------------------------------------------------------------
provider "vault" {
  address = var.vault_address
}

# ---------------------------------------------------------------------------
# Proxmox API token — retrieved ephemerally from Vault (requires vault
# provider >= 5.9).
#
# "ephemeral" means Terraform reads the value during apply but never writes
# it to state or the plan file.  The Proxmox API token stored in Vault is
# therefore never at rest in your Terraform state backend.
#
# Expected Vault secret shape:
#   vault kv put secret/proxmox/api-token \
#     api_token="root@pam!my-token=xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx"
# ---------------------------------------------------------------------------
ephemeral "vault_kv_secret_v2" "proxmox_api_token" {
  mount = var.vault_kv_mount
  name  = var.vault_proxmox_token_path
}

# ---------------------------------------------------------------------------
# Proxmox provider — authenticated with the ephemeral API token.
#
# Provider configuration arguments are never stored in state by Terraform,
# and the token itself comes from an ephemeral resource, so the credential
# never appears in state or plan files.
# ---------------------------------------------------------------------------
provider "proxmox" {
  endpoint  = var.proxmox_endpoint
  api_token = ephemeral.vault_kv_secret_v2.proxmox_api_token.data["api_token"]
  insecure  = var.proxmox_insecure
}

# ---------------------------------------------------------------------------
# VM password — retrieved ephemerally from Vault.
#
# The password is fetched at apply time and flows into the write-only
# `initialization.user_account.password` attribute — Proxmox receives it
# but Terraform never records it in state.
#
# The username is NOT sensitive and is stored in state normally (see
# var.vm_username).  Proxmox echoes the username back from its API, so
# write-only semantics are not needed (and would prevent drift tracking).
#
# Expected Vault secret shape:
#   vault kv put secret/proxmox/vm-credentials \
#     password="s3cr3t!"
# ---------------------------------------------------------------------------
ephemeral "vault_kv_secret_v2" "vm_credentials" {
  mount = var.vault_kv_mount
  name  = var.vault_vm_credentials_path
}

# ---------------------------------------------------------------------------
# VM — cloned from a template and configured with cloud-init credentials.
#
# Uses proxmox_vm (Plugin Framework resource) which supports write-only
# attributes for ephemeral values.  The password is applied to Proxmox during
# apply and is never stored in Terraform state or plan files.
# ---------------------------------------------------------------------------
resource "proxmox_vm" "example" {
  id        = var.vm_id
  name      = var.vm_name
  node_name = var.proxmox_node

  clone = {
    vm_id = var.template_vm_id
    full  = true
  }

  cpu = {
    cores = var.vm_cores
    type  = "x86-64-v2-AES"
  }

  memory = {
    size = var.vm_memory_mb
  }

  agent = {
    enabled = true
  }

  network_device = [{
    model  = "virtio"
    bridge = var.vm_network_bridge
  }]

  initialization = {
    ip_config = [{
      ipv4_address = "dhcp"
    }]

    user_account = {
      # username is not sensitive; stored in state and tracked for drift
      username = var.vm_username
      # password is write-only: applied to Proxmox but never stored in state
      password = ephemeral.vault_kv_secret_v2.vm_credentials.data["password"]
    }
  }

  started         = true
  stop_on_destroy = true
}
