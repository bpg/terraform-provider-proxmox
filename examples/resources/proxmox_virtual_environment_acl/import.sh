#!/usr/bin/env sh
# ACL can be imported using its unique identifier, e.g.: {path}?{group|user@realm|user@realm!token}?{role}
terraform import proxmox_virtual_environment_acl.operations_automation_monitoring /?monitor@pve?operations-monitoring
