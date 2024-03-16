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

// TestRuleInstantiation tests whether the Rules instance can be instantiated.
func TestRuleInstantiation(t *testing.T) {
	t.Parallel()
	require.NotNilf(t, Rules(), "Cannot instantiate Rules")
}

// TestRuleSchema tests the Rules Schema.
func TestRuleSchema(t *testing.T) {
	t.Parallel()

	rules := Rules().Schema

	test.AssertRequiredArguments(t, rules, []string{
		MkRule,
	})

	test.AssertOptionalArguments(t, rules, []string{
		mkSelectorVMID,
		mkSelectorNodeName,
	})

	nested := test.AssertNestedSchemaExistence(t, rules, MkRule)

	test.AssertOptionalArguments(t, nested, []string{
		mkSecurityGroup,
		mkRuleAction,
		mkRuleType,
		mkRuleComment,
		mkRuleDest,
		mkRuleDPort,
		mkRuleEnabled,
		mkRuleIFace,
		mkRuleLog,
		mkRuleMacro,
		mkRuleProto,
		mkRuleSource,
		mkRuleSPort,
	})

	test.AssertValueTypes(t, nested, map[string]schema.ValueType{
		mkRulePos:     schema.TypeInt,
		mkRuleAction:  schema.TypeString,
		mkRuleType:    schema.TypeString,
		mkRuleComment: schema.TypeString,
		mkRuleDest:    schema.TypeString,
		mkRuleDPort:   schema.TypeString,
		mkRuleEnabled: schema.TypeBool,
		mkRuleIFace:   schema.TypeString,
		mkRuleLog:     schema.TypeString,
		mkRuleMacro:   schema.TypeString,
		mkRuleProto:   schema.TypeString,
		mkRuleSource:  schema.TypeString,
		mkRuleSPort:   schema.TypeString,
	})
}
