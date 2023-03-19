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

	"github.com/bpg/terraform-provider-proxmox/proxmoxtf/test"
)

// TestSecurityGroupInstantiation tests whether the SecurityGroup instance can be instantiated.
func TestSecurityGroupInstantiation(t *testing.T) {
	t.Parallel()
	require.NotNilf(t, SecurityGroup(), "Cannot instantiate SecurityGroup")
}

// TestSecurityGroupSchema tests the SecurityGroup schema.
func TestSecurityGroupSchema(t *testing.T) {
	t.Parallel()
	s := SecurityGroup()

	test.AssertRequiredArguments(t, s, []string{
		mkGroupName,
	})

	test.AssertOptionalArguments(t, s, []string{
		mkGroupComment,
	})

	test.AssertValueTypes(t, s, map[string]schema.ValueType{
		mkGroupName:    schema.TypeString,
		mkGroupComment: schema.TypeString,
	})

	ruleSchema := test.AssertNestedSchemaExistence(t, s, mkRule)

	test.AssertRequiredArguments(t, ruleSchema, []string{
		mkRuleAction,
		mkRuleType,
	})

	test.AssertOptionalArguments(t, ruleSchema, []string{
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

	test.AssertValueTypes(t, ruleSchema, map[string]schema.ValueType{
		mkRulePos:     schema.TypeInt,
		mkRuleAction:  schema.TypeString,
		mkRuleType:    schema.TypeString,
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
