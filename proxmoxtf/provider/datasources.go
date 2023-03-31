/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package provider

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/bpg/terraform-provider-proxmox/proxmoxtf/datasource"
)

func createDatasourceMap() map[string]*schema.Resource {
	return map[string]*schema.Resource{
		// "proxmox_virtual_environment_cluster_firewall_alias":           cluster.FirewallAlias(),
		// "proxmox_virtual_environment_cluster_firewall_aliases":         cluster.FirewallAliases(),
		// "proxmox_virtual_environment_cluster_firewall_ipset":           cluster.FirewallIPSet(),
		// "proxmox_virtual_environment_cluster_firewall_ipsets":          cluster.FirewallIPSets(),
		// "proxmox_virtual_environment_cluster_firewall_security_group":  cluster.FirewallSecurityGroup(),
		// "proxmox_virtual_environment_cluster_firewall_security_groups": cluster.FirewallSecurityGroups(),
		"proxmox_virtual_environment_datastores": datasource.Datastores(),
		"proxmox_virtual_environment_dns":        datasource.DNS(),
		"proxmox_virtual_environment_group":      datasource.Group(),
		"proxmox_virtual_environment_groups":     datasource.Groups(),
		"proxmox_virtual_environment_hosts":      datasource.Hosts(),
		"proxmox_virtual_environment_nodes":      datasource.Nodes(),
		"proxmox_virtual_environment_pool":       datasource.Pool(),
		"proxmox_virtual_environment_pools":      datasource.Pools(),
		"proxmox_virtual_environment_role":       datasource.Role(),
		"proxmox_virtual_environment_roles":      datasource.Roles(),
		"proxmox_virtual_environment_time":       datasource.Time(),
		"proxmox_virtual_environment_user":       datasource.User(),
		"proxmox_virtual_environment_users":      datasource.Users(),
		"proxmox_virtual_environment_version":    datasource.Version(),
		"proxmox_virtual_environment_vm":         datasource.VM(),
		"proxmox_virtual_environment_vms":        datasource.VMs(),
	}
}
