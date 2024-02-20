/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package provider

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/bpg/terraform-provider-proxmox/proxmoxtf/resource"
	clusterfirewall "github.com/bpg/terraform-provider-proxmox/proxmoxtf/resource/cluster/firewall"
	container "github.com/bpg/terraform-provider-proxmox/proxmoxtf/resource/container"
	"github.com/bpg/terraform-provider-proxmox/proxmoxtf/resource/firewall"
	vm "github.com/bpg/terraform-provider-proxmox/proxmoxtf/resource/vm"
)

func createResourceMap() map[string]*schema.Resource {
	return map[string]*schema.Resource{
		"proxmox_virtual_environment_certificate":                     resource.Certificate(),
		"proxmox_virtual_environment_cluster_firewall":                clusterfirewall.Firewall(),
		"proxmox_virtual_environment_cluster_firewall_security_group": clusterfirewall.SecurityGroup(),
		"proxmox_virtual_environment_container":                       container.Container(),
		"proxmox_virtual_environment_dns":                             resource.DNS(),
		"proxmox_virtual_environment_file":                            resource.File(),
		"proxmox_virtual_environment_firewall_alias":                  firewall.Alias(),
		"proxmox_virtual_environment_firewall_ipset":                  firewall.IPSet(),
		"proxmox_virtual_environment_firewall_options":                firewall.Options(),
		"proxmox_virtual_environment_firewall_rules":                  firewall.Rules(),
		"proxmox_virtual_environment_group":                           resource.Group(),
		"proxmox_virtual_environment_hosts":                           resource.Hosts(),
		"proxmox_virtual_environment_pool":                            resource.Pool(),
		"proxmox_virtual_environment_role":                            resource.Role(),
		"proxmox_virtual_environment_time":                            resource.Time(),
		"proxmox_virtual_environment_user":                            resource.User(),
		"proxmox_virtual_environment_vm":                              vm.VM(),
	}
}
