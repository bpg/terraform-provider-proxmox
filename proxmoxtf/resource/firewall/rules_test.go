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

// TestMapToBaseRuleWithEmptyValues tests empty value handling for issue #1504.
func TestMapToBaseRuleWithEmptyValues(t *testing.T) {
	t.Parallel()

	rule := map[string]interface{}{
		mkRuleComment: "",
		mkRuleDest:    "",
		mkRuleDPort:   "",
		mkRuleEnabled: true,
		mkRuleIFace:   "",
		mkRuleLog:     "",
		mkRuleMacro:   "",
		mkRuleProto:   "",
		mkRuleSource:  "",
		mkRuleSPort:   "",
	}

	baseRule := mapToBaseRule(rule)

	require.NotNil(t, baseRule.Comment)
	require.NotNil(t, baseRule.Dest)
	require.NotNil(t, baseRule.DPort)
	require.NotNil(t, baseRule.Enable)
	require.NotNil(t, baseRule.Macro)
	require.NotNil(t, baseRule.Proto)
	require.NotNil(t, baseRule.Source)
	require.NotNil(t, baseRule.SPort)

	require.Empty(t, *baseRule.Comment)
	require.Empty(t, *baseRule.Dest)
	require.Empty(t, *baseRule.DPort)
	require.True(t, bool(*baseRule.Enable))
	require.Empty(t, *baseRule.Macro)
	require.Empty(t, *baseRule.Proto)
	require.Empty(t, *baseRule.Source)
	require.Empty(t, *baseRule.SPort)

	require.Nil(t, baseRule.IFace)
	require.Nil(t, baseRule.Log)
}

// TestMapToBaseRuleWithNonEmptyValues tests non-empty value handling.
func TestMapToBaseRuleWithNonEmptyValues(t *testing.T) {
	t.Parallel()

	rule := map[string]interface{}{
		mkRuleComment: "Test comment",
		mkRuleDest:    "192.168.1.5",
		mkRuleDPort:   "80",
		mkRuleEnabled: false,
		mkRuleIFace:   "net0",
		mkRuleLog:     "info",
		mkRuleMacro:   "HTTP",
		mkRuleProto:   "tcp",
		mkRuleSource:  "192.168.1.0/24",
		mkRuleSPort:   "8080",
	}

	baseRule := mapToBaseRule(rule)

	require.Equal(t, "Test comment", *baseRule.Comment)
	require.Equal(t, "192.168.1.5", *baseRule.Dest)
	require.Equal(t, "80", *baseRule.DPort)
	require.False(t, bool(*baseRule.Enable))
	require.Equal(t, "net0", *baseRule.IFace)
	require.Equal(t, "info", *baseRule.Log)
	require.Equal(t, "HTTP", *baseRule.Macro)
	require.Equal(t, "tcp", *baseRule.Proto)
	require.Equal(t, "192.168.1.0/24", *baseRule.Source)
	require.Equal(t, "8080", *baseRule.SPort)
}
