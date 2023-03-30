/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package firewall

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/require"

	"github.com/bpg/terraform-provider-proxmox/proxmoxtf/structure"
)

// TestSecurityGroupSchemaInstantiation tests whether the SecurityGroupSchema instance can be instantiated.
func TestSecurityGroupSchemaInstantiation(t *testing.T) {
	t.Parallel()
	require.NotNilf(t, SecurityGroupSchema(), "Cannot instantiate SecurityGroupSchema")
}

// TestSecurityGroupSchema tests the SecurityGroupSchema.
func TestSecurityGroupSchema(t *testing.T) {
	t.Parallel()
	s := SecurityGroupSchema()

	structure.AssertRequiredArguments(t, s, []string{
		mkSecurityGroupName,
	})

	structure.AssertOptionalArguments(t, s, []string{
		mkSecurityGroupComment,
	})

	structure.AssertValueTypes(t, s, map[string]schema.ValueType{
		mkSecurityGroupName:    schema.TypeString,
		mkSecurityGroupComment: schema.TypeString,
	})

	ruleSchema := structure.AssertNestedSchemaExistence(t, s, MkRule).Schema

	structure.AssertRequiredArguments(t, ruleSchema, []string{
		MkRuleAction,
		MkRuleType,
	})

	structure.AssertOptionalArguments(t, ruleSchema, []string{
		mkRuleComment,
		mkRuleDest,
		mkRuleDPort,
		mkRuleEnable,
		mkRuleIFace,
		mkRuleLog,
		mkRuleMacro,
		mkRuleProto,
		mkRuleSource,
		mkRuleSPort,
	})

	structure.AssertValueTypes(t, ruleSchema, map[string]schema.ValueType{
		MkRulePos:     schema.TypeInt,
		MkRuleAction:  schema.TypeString,
		MkRuleType:    schema.TypeString,
		mkRuleComment: schema.TypeString,
		mkRuleDest:    schema.TypeString,
		mkRuleDPort:   schema.TypeString,
		mkRuleEnable:  schema.TypeBool,
		mkRuleIFace:   schema.TypeString,
		mkRuleLog:     schema.TypeString,
		mkRuleMacro:   schema.TypeString,
		mkRuleProto:   schema.TypeString,
		mkRuleSource:  schema.TypeString,
		mkRuleSPort:   schema.TypeString,
	})
}
