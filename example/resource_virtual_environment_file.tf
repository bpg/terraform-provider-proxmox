#===============================================================================
# Cloud Config (cloud-init)
#===============================================================================

resource "proxmox_virtual_environment_file" "user_config" {
  content_type = "snippets"
  datastore_id = element(data.proxmox_virtual_environment_datastores.example.datastore_ids, index(data.proxmox_virtual_environment_datastores.example.datastore_ids, "local"))
  node_name    = data.proxmox_virtual_environment_datastores.example.node_name

  source_raw {
    data = <<-EOF
    #cloud-config
    chpasswd:
      list: |
        ubuntu:example
      expire: false
    hostname: terraform-provider-proxmox-example
    users:
      - default
      - name: ubuntu
        groups: sudo
        shell: /bin/bash
        ssh-authorized-keys:
          - ${trimspace(tls_private_key.example.public_key_openssh)}
        sudo: ALL=(ALL) NOPASSWD:ALL
    EOF

    file_name = "terraform-provider-proxmox-example-user-config.yaml"
  }
}

resource "proxmox_virtual_environment_file" "vendor_config" {
  content_type = "snippets"
  datastore_id = element(data.proxmox_virtual_environment_datastores.example.datastore_ids, index(data.proxmox_virtual_environment_datastores.example.datastore_ids, "local"))
  node_name    = data.proxmox_virtual_environment_datastores.example.node_name

  source_raw {
    data = <<EOF
#cloud-config
runcmd:
    - apt update
    - apt install -y qemu-guest-agent
    - systemctl enable qemu-guest-agent
    - systemctl start qemu-guest-agent
    - echo "done" > /tmp/vendor-cloud-init-done
    EOF

    file_name = "terraform-provider-proxmox-example-vendor-config.yaml"
  }
}


resource "proxmox_virtual_environment_file" "meta_config" {
  content_type = "snippets"
  datastore_id = element(data.proxmox_virtual_environment_datastores.example.datastore_ids, index(data.proxmox_virtual_environment_datastores.example.datastore_ids, "local"))
  node_name    = data.proxmox_virtual_environment_datastores.example.node_name

  source_raw {
    data = <<EOF
local-hostname: myhost.internal
    EOF

    file_name = "meta-config.yaml"
  }
}

#===============================================================================
# Snippets
#===============================================================================

resource "proxmox_virtual_environment_file" "hook_script" {
  content_type = "snippets"
  datastore_id = element(data.proxmox_virtual_environment_datastores.example.datastore_ids, index(data.proxmox_virtual_environment_datastores.example.datastore_ids, "local"))
  node_name    = data.proxmox_virtual_environment_datastores.example.node_name
  # Hook scripts must be executable, otherwise the Proxmox VE API will reject the configuration for the VM/CT.
  file_mode = "0700"

  source_raw {
    data = <<-EOF
      #!/usr/bin/env bash

      echo "Running hook script"
      EOF
    file_name = "prepare-hook.sh"
  }
}
