#!/usr/bin/env sh
# Resource pool membership can be imported using its unique identifier, e.g.: {pool_id}/{type}/{member_id}
terraform import proxmox_virtual_environment_pool_membership.example_membership test-pool/vm/102
