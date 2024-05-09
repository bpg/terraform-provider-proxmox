#!/usr/bin/env sh
# ACL can be imported using its unique identifier, e.g.: {path}?entity_id={group|user@realm|user@realm!token}?role_id={role}
terraform import proxmox_virtual_environment_acl.operations_automation_monitoring /?entity_id=monitor@pve&role_id=operations-monitoring
