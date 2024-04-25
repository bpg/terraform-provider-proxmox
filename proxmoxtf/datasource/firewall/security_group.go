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
	mkSecurityGroupName    = "name"
	mkSecurityGroupComment = "comment"
	mkRules                = "rules"
)

// SecurityGroupSchema defines the schema for the security group data source.
func SecurityGroupSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		mkSecurityGroupName: {
			Type:        schema.TypeString,
			Description: "Security group name",
			Required:    true,
		},
		mkSecurityGroupComment: {
			Type:        schema.TypeString,
			Description: "Security group comment",
			Computed:    true,
		},
		mkRules: {
			Type:        schema.TypeList,
			Description: "List of rules",
			Computed:    true,
			Elem:        &schema.Resource{Schema: RuleSchema()},
		},
	}
}
