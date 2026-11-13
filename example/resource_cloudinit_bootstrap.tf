resource "proxmox_download_file" "ubuntu_24_04_cloud_img" {
  content_type = "import"
  datastore_id = "local"
  file_name    = "cloudinit-bootstrap-example-noble-amd64.qcow2"
  node_name    = data.proxmox_virtual_environment_nodes.example.names[0]
  url          = var.ubuntu_24_04_cloud_img_url
}

resource "proxmox_virtual_environment_vm" "cloudinit_bootstrap" {
  name      = "terraform-provider-proxmox-cloudinit-bootstrap"
  node_name = data.proxmox_virtual_environment_nodes.example.names[0]

  # proxmox_cloudinit_bootstrap starts the VM once the seed ISO is attached.
  started = false

  # The VM ends up running, so force-stop it on destroy; a graceful ACPI shutdown
  # may hang if the guest agent is not installed.
  stop_on_destroy = true

  disk {
    datastore_id = local.datastore_id
    import_from  = proxmox_download_file.ubuntu_24_04_cloud_img.id
    interface    = "virtio0"
    size         = 10
  }

  cdrom {
    interface = "ide0"
    file_id   = "none"
  }

  network_device {}

  agent {
    enabled = true
  }

  lifecycle {
    ignore_changes = [started]
  }
}

resource "proxmox_cloudinit_bootstrap" "example" {
  vm_id        = proxmox_virtual_environment_vm.cloudinit_bootstrap.vm_id
  datastore_id = "local"
  file_name    = "cloudinit-bootstrap-example-seed.iso"

  user_data = <<-EOT
    #cloud-config
    packages: [qemu-guest-agent]
    runcmd:
      - [systemctl, enable, --now, qemu-guest-agent]
  EOT

  network_config = yamlencode({
    version = 1
    config = [{
      type        = "physical"
      name        = "eth0"
      mac_address = lower(proxmox_virtual_environment_vm.cloudinit_bootstrap.network_device[0].mac_address)
      subnets     = [{ type = "dhcp" }]
    }]
  })

  lifecycle {
    replace_triggered_by = [proxmox_virtual_environment_vm.cloudinit_bootstrap.id]
  }
}

output "resource_proxmox_cloudinit_bootstrap_example_ipv4_addresses" {
  value = proxmox_cloudinit_bootstrap.example.ipv4_addresses
}
