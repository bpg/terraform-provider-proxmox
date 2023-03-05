/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package provider

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/bpg/terraform-provider-proxmox/proxmoxtf/datasource"
	"github.com/bpg/terraform-provider-proxmox/proxmoxtf/datasource/firewall"
)

func createDatasourceMap() map[string]*schema.Resource {
	return map[string]*schema.Resource{
		"proxmox_virtual_environment_firewall_alias":   firewall.Alias(),
		"proxmox_virtual_environment_firewall_aliases": firewall.Aliases(),
		"proxmox_virtual_environment_firewall_ipset":   firewall.IPSet(),
		"proxmox_virtual_environment_firewall_ipsets":  firewall.IPSets(),
		"proxmox_virtual_environment_datastores":       datasource.Datastores(),
		"proxmox_virtual_environment_dns":              datasource.DNS(),
		"proxmox_virtual_environment_group":            datasource.Group(),
		"proxmox_virtual_environment_groups":           datasource.Groups(),
		"proxmox_virtual_environment_hosts":            datasource.Hosts(),
		"proxmox_virtual_environment_nodes":            datasource.Nodes(),
		"proxmox_virtual_environment_pool":             datasource.Pool(),
		"proxmox_virtual_environment_pools":            datasource.Pools(),
		"proxmox_virtual_environment_role":             datasource.Role(),
		"proxmox_virtual_environment_roles":            datasource.Roles(),
		"proxmox_virtual_environment_time":             datasource.Time(),
		"proxmox_virtual_environment_user":             datasource.User(),
		"proxmox_virtual_environment_users":            datasource.Users(),
		"proxmox_virtual_environment_version":          datasource.Version(),
	}
}
