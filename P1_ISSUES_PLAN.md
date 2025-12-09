# P1 Priority Issues - Fix Plan

## Most Critical Issues (Ranked by Impact)

### 1. #2368 - VM Creation Timeout (BLOCKING) ✅
- **Title**: virtio-scsi-single with agent.enabled = true causes indefinite timeout on VM creation
- **URL**: https://github.com/bpg/terraform-provider-proxmox/issues/2368
- **Impact**: Blocks VM creation entirely
- **Created**: Nov 22, 2025 (recent)
- **Status**: Completed - Fixed HTTP 500 error retry logic for agent detection

### 2. #2218 - Disk Removal Not Working ✅
- **Title**: Disk removal does not remove disks, and causes downstream issues
- **URL**: https://github.com/bpg/terraform-provider-proxmox/issues/2218
- **Impact**: Removed disks aren't cleaned up, causing conflicts
- **Created**: Oct 2, 2025
- **Status**: Completed - Fixed disk size mismatch when re-adding removed disk

### 3. #2259 - Network Device Removal Not Working
- **Title**: Removing a network device from a VM configuration does not remove it from Proxmox
- **URL**: https://github.com/bpg/terraform-provider-proxmox/issues/2259
- **Impact**: Network devices aren't cleaned up
- **Created**: Oct 20, 2025

### 4. #1435 - Container Disk Resize Forces Replacement (High User Impact)
- **Title**: Changing disk size of 'proxmox_virtual_environment_container' forces replacement
- **URL**: https://github.com/bpg/terraform-provider-proxmox/issues/1435
- **Impact**: 9 reactions, 8 comments - unnecessary container replacement
- **Created**: Jul 10, 2024

### 5. #2195 - IP Config Not Removed
- **Title**: `ipconfig` not removed when deleting a network device
- **URL**: https://github.com/bpg/terraform-provider-proxmox/issues/2195
- **Impact**: Related to #2259 - IP config cleanup issue
- **Created**: Sep 22, 2025

### 6. #1515 - EFI Disk Changes Recreate VM ✅
- **Title**: Changing EFI disk's parameters should not recreate the whole VM
- **URL**: https://github.com/bpg/terraform-provider-proxmox/issues/1515
- **Impact**: 3 reactions - unnecessary VM recreation
- **Created**: Sep 5, 2024
- **Status**: Completed - Removed ForceNew from updatable EFI disk parameters

### 7. #538 - Unwanted VM Reboots (Old, High Engagement)
- **Title**: VM Should not be rebooted on hotplug added resources
- **URL**: https://github.com/bpg/terraform-provider-proxmox/issues/538
- **Impact**: 6 reactions - unwanted reboots during hotplug operations
- **Created**: Sep 4, 2023
- **Status**: In progress

### 8. #1998 - Import Causes Recreation
- **Title**: Field `file_id` not populated when importing vm, causing recreation
- **URL**: https://github.com/bpg/terraform-provider-proxmox/issues/1998
- **Impact**: Import functionality broken - 5 comments
- **Created**: Jun 16, 2025

## Summary

**Top Priorities:**
1. **#2368** - Blocks VM creation (critical blocker)
2. **#2218, #2259, #2195** - Resource cleanup failures (data consistency)
3. **#1435** - High user impact (9 reactions)
4. **#1515, #538** - Unnecessary recreations/reboots

