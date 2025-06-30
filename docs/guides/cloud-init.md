---
layout: page
page_title: "Configure a VM with Cloud-Init"
subcategory: Guides
description: |-
    This guide explains how to use the Proxmox provider to create and manage virtual machines using cloud-init.
---

# Configure a VM with Cloud-Init

## Native Proxmox Cloud-Init Support

Proxmox supports cloud-init natively, so you can use the `initialization` block to configure your VM:

```terraform
data "local_file" "ssh_public_key" {
  filename = "./id_rsa.pub"
}

resource "proxmox_virtual_environment_vm" "ubuntu_vm" {
  name      = "test-ubuntu"
  node_name = "pve"

  # should be true if qemu agent is not installed / enabled on the VM
  stop_on_destroy = true

  initialization {
    ip_config {
      ipv4 {
        address = "192.168.3.233/24"
        gateway = "192.168.3.1"
      }
    }

    user_account {
      username = "ubuntu"
      keys     = [trimspace(data.local_file.ssh_public_key.content)]
    }
  }

  disk {
    datastore_id = "local-lvm"
    file_id      = proxmox_virtual_environment_download_file.ubuntu_cloud_image.id
    interface    = "virtio0"
    iothread     = true
    discard      = "on"
    size         = 20
  }

  network_device {
    bridge = "vmbr0"
  }
}

resource "proxmox_virtual_environment_download_file" "ubuntu_cloud_image" {
  content_type = "iso"
  datastore_id = "local"
  node_name    = "pve"

  url = "https://cloud-images.ubuntu.com/jammy/current/jammy-server-cloudimg-amd64.img"
}
```

Note that many cloud images do not have `qemu-guest-agent` installed by default, so you won't be able to retrieve the dynamic IP address of the VM from Proxmox, as this is agent's responsibility. You can use the `ip_config` block to configure a static IP address instead.

## Custom Cloud-Init Configuration

Due to several limitations of the native Proxmox cloud-init support, you may want to use a custom Cloud-Init configuration instead. This allows you to adjust the VM configuration to your needs and install the `qemu-guest-agent` and additional packages.

To use a custom cloud-init configuration, create a cloud-config snippet file and pass it to the VM as a `user_data_file_id` parameter. Use the `proxmox_virtual_environment_file` resource to create the file. Ensure the "Snippets" content type is enabled on the target datastore in Proxmox before applying the configuration below.

Note that you need to explicitly set the `hostname` in the provided cloud-init configuration, as the custom user data cloud-init config overwrites the config set by Proxmox, as shown in the example below.

```terraform
data "local_file" "ssh_public_key" {
  filename = "./id_rsa.pub"
}

resource "proxmox_virtual_environment_file" "user_data_cloud_config" {
  content_type = "snippets"
  datastore_id = "local"
  node_name    = "pve"

  source_raw {
    data = <<-EOF
    #cloud-config
    hostname: test-ubuntu
    timezone: America/Toronto
    users:
      - default
      - name: ubuntu
        groups:
          - sudo
        shell: /bin/bash
        ssh_authorized_keys:
          - ${trimspace(data.local_file.ssh_public_key.content)}
        sudo: ALL=(ALL) NOPASSWD:ALL
    package_update: true
    packages:
      - qemu-guest-agent
      - net-tools
      - curl
    runcmd:
      - systemctl enable qemu-guest-agent
      - systemctl start qemu-guest-agent
      - echo "done" > /tmp/cloud-config.done
    EOF

    file_name = "user-data-cloud-config.yaml"
  }
}
```

```terraform
resource "proxmox_virtual_environment_vm" "ubuntu_vm" {
  name      = "test-ubuntu"
  node_name = "pve"

  agent {
    enabled = true
  }

  cpu {
    cores = 2
  }

  memory {
    dedicated = 2048
  }

  disk {
    datastore_id = "local-lvm"
    file_id      = proxmox_virtual_environment_download_file.ubuntu_cloud_image.id
    interface    = "virtio0"
    iothread     = true
    discard      = "on"
    size         = 20
  }

  initialization {
    ip_config {
      ipv4 {
        address = "dhcp"
      }
    }

    user_data_file_id = proxmox_virtual_environment_file.user_data_cloud_config.id
  }

  network_device {
    bridge = "vmbr0"
  }

}

resource "proxmox_virtual_environment_download_file" "ubuntu_cloud_image" {
  content_type = "iso"
  datastore_id = "local"
  node_name    = "pve"

  url = "https://cloud-images.ubuntu.com/jammy/current/jammy-server-cloudimg-amd64.img"
}

output "vm_ipv4_address" {
  value = proxmox_virtual_environment_vm.ubuntu_vm.ipv4_addresses[1][0]
}
```

If you wish to keep the user data cloud-init config generic, for example, when deploying multiple VMs for a Kubernetes cluster, you can use a separate snippet with the metadata cloud-init config to set the hostname. Note that it uses the `local-hostname` configuration parameter. See more details in the [cloud-init documentation](https://docs.cloud-init.io/en/latest/reference/yaml_examples/update_hostname.html).

```terraform
resource "proxmox_virtual_environment_file" "meta_data_cloud_config" {
  content_type = "snippets"
  datastore_id = "local"
  node_name    = "pve"

  source_raw {
    data = <<-EOF
    #cloud-config
    local-hostname: test-ubuntu
    EOF

    file_name = "meta-data-cloud-config.yaml"
  }
}
```

```terraform
resource "proxmox_virtual_environment_vm" "ubuntu_vm" {
  # ...

  initialization {
    # ...
    user_data_file_id = proxmox_virtual_environment_file.user_data_cloud_config.id
    meta_data_file_id = proxmox_virtual_environment_file.meta_data_cloud_config.id
  }

  # ...
}
```
