/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package firewall

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const (
	mkSecurityGroupsSecurityGroupNames = "security_group_names"
)

// SecurityGroupsSchema defines the schema for the security groups.
func SecurityGroupsSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		mkSecurityGroupsSecurityGroupNames: {
			Type:        schema.TypeList,
			Description: "Security Group Names",
			Computed:    true,
			Elem:        &schema.Schema{Type: schema.TypeString},
		},
	}
}
