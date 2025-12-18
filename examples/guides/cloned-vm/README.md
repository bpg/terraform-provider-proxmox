# Cloned VM Example

This example demonstrates the `proxmox_virtual_environment_cloned_vm` resource, which provides fine-grained control over VM cloning with explicit device management.

## Features Demonstrated

1. **Partial Management**: Clone a VM and manage only specific network interfaces while preserving others
2. **Selective Deletion**: Explicitly delete inherited devices using the `delete` block
3. **Full Management**: Manage all inherited devices with custom configurations
4. **Disk Management**: Resize inherited disks and add new disks using slot-based addressing

## Key Concepts

### Map-Based Device Addressing

Unlike the legacy clone approach, `cloned_vm` uses map-based addressing for devices:

- Network devices: `net0`, `net1`, `net2`, etc.
- Disk devices: `scsi0`, `virtio0`, `sata0`, `ide0`, etc.

This provides stable, addressable references to specific devices.

### Opt-In Management

Only devices and configuration explicitly declared in your Terraform code are managed. Inherited settings from the template are preserved unless:

1. You explicitly configure them (override)
2. You list them in the `delete` block (remove)

### Explicit Deletion Semantics

Removing a device from your Terraform configuration **does not delete it from the VM**. To delete inherited devices, you must explicitly list them in the `delete` block:

```terraform
delete = {
  network = ["net1", "net2"]
  disk    = ["scsi1"]
}
```

## Prerequisites

1. A Proxmox VE cluster (version 9.x recommended)
2. An API token with appropriate permissions
3. SSH access configured for the Terraform user
4. An SSH public key file at `./id_rsa.pub`

## Configuration

Create a `terraform.tfvars` file:

```hcl
virtual_environment_endpoint   = "https://your-proxmox-host:8006/"
virtual_environment_token      = "root@pam!terraform=your-api-token"
virtual_environment_node_name  = "pve"
datastore_id                   = "local-lvm"
```

## Usage

```bash
# Initialize Terraform
terraform init

# Review the plan
terraform plan

# Apply the configuration
terraform apply

# Clean up
terraform destroy
```

## What Gets Created

1. **Template VM**: A multi-NIC Ubuntu template with cloud-init configuration
2. **Partial Management Clone**: Manages only `net0`, preserves `net1` and `net2`
3. **Selective Deletion Clone**: Manages `net0`, deletes `net1` and `net2`
4. **Full Management Clone**: Explicitly manages all three network interfaces
5. **Disk Management Clone**: Resizes boot disk and adds a data disk

## Examples Explained

### Example 1: Partial Management

```terraform
resource "proxmox_virtual_environment_cloned_vm" "partial_management" {
  # Only manage net0, net1 and net2 remain on the VM but aren't tracked
  network = {
    net0 = { ... }
  }
}
```

**Result**: VM has all 3 NICs, but Terraform only tracks/manages `net0`.

### Example 2: Selective Deletion

```terraform
resource "proxmox_virtual_environment_cloned_vm" "selective_deletion" {
  network = {
    net0 = { ... }
  }

  delete = {
    network = ["net1", "net2"]
  }
}
```

**Result**: VM has only `net0`, the other NICs are deleted during creation.

### Example 3: Full Management

```terraform
resource "proxmox_virtual_environment_cloned_vm" "full_management" {
  network = {
    net0 = { tag = 100 }
    net1 = { tag = 200 }
    net2 = { tag = 300 }
  }
}
```

**Result**: All 3 NICs are managed with different VLAN tags.

### Example 4: Disk Management

```terraform
resource "proxmox_virtual_environment_cloned_vm" "disk_management" {
  disk = {
    virtio0 = { size_gb = 50 }  # Resize from 20GB
    virtio1 = { size_gb = 100 } # Add new disk
  }
}
```

**Result**: Boot disk expanded to 50GB, new 100GB data disk added.

## When to Use This Resource

Use `proxmox_virtual_environment_cloned_vm` when:

- You need fine-grained control over which devices are managed
- Your templates have complex device configurations
- You want to preserve some inherited configuration without tracking it in Terraform
- You need explicit control over device deletion
- You want stable, slot-based device addressing

For simpler clone scenarios, consider using `proxmox_virtual_environment_vm` with the `clone` block.

## See Also

- [Resource Documentation](https://registry.terraform.io/providers/bpg/proxmox/latest/docs/resources/virtual_environment_cloned_vm)
- [Clone VM Guide](https://registry.terraform.io/providers/bpg/proxmox/latest/docs/guides/clone-vm)
- [VM2 Clone Migration Guide](https://registry.terraform.io/providers/bpg/proxmox/latest/docs/guides/migration-vm2-clone)
