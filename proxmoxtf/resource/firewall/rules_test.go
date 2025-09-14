/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package firewall

import (
	"context"
	"fmt"
	"sort"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/require"

	"github.com/bpg/terraform-provider-proxmox/proxmox/firewall"
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

// TestRulesReadDriftDetection tests drift detection when rules are manually deleted.
func TestRulesReadDriftDetection(t *testing.T) {
	t.Parallel()

	mockAPI := &mockFirewallRuleAPI{
		rules: map[int]*firewall.RuleGetResponseData{
			0: {
				Action: "ACCEPT",
				Type:   "in",
				BaseRule: firewall.BaseRule{
					Comment: stringPtr("Allow HTTP"),
					DPort:   stringPtr("80"),
					Proto:   stringPtr("tcp"),
				},
			},
			2: {
				Action: "ACCEPT",
				Type:   "in",
				BaseRule: firewall.BaseRule{
					Comment: stringPtr("Allow HTTPS"),
					DPort:   stringPtr("443"),
					Proto:   stringPtr("tcp"),
				},
			},
		},
		rulesID: "test-rules-id",
	}

	// Create ResourceData with the rules schema
	rulesResource := Rules()
	d := rulesResource.TestResourceData()
	err := d.Set(MkRule, []interface{}{
		map[string]interface{}{
			mkRulePos:     0,
			mkRuleAction:  "ACCEPT",
			mkRuleType:    "in",
			mkRuleComment: "Allow HTTP",
			mkRuleDPort:   "80",
			mkRuleProto:   "tcp",
		},
		map[string]interface{}{
			mkRulePos:     1,
			mkRuleAction:  "ACCEPT",
			mkRuleType:    "in",
			mkRuleComment: "Allow SSH",
			mkRuleDPort:   "22",
			mkRuleProto:   "tcp",
		},
		map[string]interface{}{
			mkRulePos:     2,
			mkRuleAction:  "ACCEPT",
			mkRuleType:    "in",
			mkRuleComment: "Allow HTTPS",
			mkRuleDPort:   "443",
			mkRuleProto:   "tcp",
		},
	})
	require.NoError(t, err)

	diags := RulesRead(context.Background(), mockAPI, d)
	require.False(t, diags.HasError(), "RulesRead should not return errors")

	rules := d.Get(MkRule).([]interface{})
	require.Len(t, rules, 2, "Should have 2 rules after drift detection")

	rule0 := rules[0].(map[string]interface{})
	require.Equal(t, "Allow HTTP", rule0[mkRuleComment])
	require.Equal(t, "80", rule0[mkRuleDPort])

	rule1 := rules[1].(map[string]interface{})
	require.Equal(t, "Allow HTTPS", rule1[mkRuleComment])
	require.Equal(t, "443", rule1[mkRuleDPort])
}

// mockFirewallRuleAPI is a mock implementation for testing.
type mockFirewallRuleAPI struct {
	rules   map[int]*firewall.RuleGetResponseData
	rulesID string
}

func (m *mockFirewallRuleAPI) GetRulesID() string {
	return m.rulesID
}

func (m *mockFirewallRuleAPI) CreateRule(ctx context.Context, d *firewall.RuleCreateRequestBody) error {
	return nil
}

func (m *mockFirewallRuleAPI) GetRule(ctx context.Context, pos int) (*firewall.RuleGetResponseData, error) {
	rule, exists := m.rules[pos]
	if !exists {
		return nil, fmt.Errorf("500 no rule at position %d", pos)
	}

	return rule, nil
}

func (m *mockFirewallRuleAPI) ListRules(ctx context.Context) ([]*firewall.RuleListResponseData, error) {
	keys := make([]int, 0, len(m.rules))
	for k := range m.rules {
		keys = append(keys, k)
	}

	sort.Ints(keys)

	result := make([]*firewall.RuleListResponseData, 0, len(m.rules))
	for _, pos := range keys {
		result = append(result, &firewall.RuleListResponseData{Pos: pos})
	}

	return result, nil
}

func (m *mockFirewallRuleAPI) UpdateRule(ctx context.Context, pos int, d *firewall.RuleUpdateRequestBody) error {
	return nil
}

func (m *mockFirewallRuleAPI) DeleteRule(ctx context.Context, pos int) error {
	return nil
}

func stringPtr(s string) *string {
	return &s
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
