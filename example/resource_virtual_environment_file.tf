#===============================================================================
# Cloud Config (cloud-init)
#===============================================================================

resource "proxmox_virtual_environment_file" "user_config" {
  content_type = "snippets"
  datastore_id = "local"
  node_name    = data.proxmox_virtual_environment_nodes.example.names[0]

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
  datastore_id = "local"
  node_name    = data.proxmox_virtual_environment_nodes.example.names[0]

  source_raw {
    data = <<-EOF
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
  datastore_id = "local"
  node_name    = data.proxmox_virtual_environment_nodes.example.names[0]

  source_raw {
    data = <<-EOF
    local-hostname: myhost.internal
    EOF

    file_name = "meta-config.yaml"
  }
}
