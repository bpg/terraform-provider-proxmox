/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package firewall

import (
	"context"
	"fmt"
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

	test.AssertOptionalArguments(t, rules, []string{
		MkRule,
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

	rule := map[string]any{
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

	rule := map[string]any{
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

// deleteTestMockAPI is a minimal mock for testing RulesDelete error handling.
// RulesDelete's skip-missing logic (ErrNoRuleAtPosition → continue) is only reachable
// via a race condition (rule disappears between Read and Delete), so it cannot be
// tested with acceptance tests.
type deleteTestMockAPI struct {
	getRuleErrors map[int]error
	deletedPos    []int
}

func (m *deleteTestMockAPI) GetRulesID() string { return "test" }

func (m *deleteTestMockAPI) ListRules(context.Context) ([]*firewall.RuleListResponseData, error) {
	return nil, nil
}

func (m *deleteTestMockAPI) CreateRule(context.Context, *firewall.RuleCreateRequestBody) error {
	return nil
}

func (m *deleteTestMockAPI) UpdateRule(context.Context, int, *firewall.RuleUpdateRequestBody) error {
	return nil
}

func (m *deleteTestMockAPI) GetRule(_ context.Context, pos int) (*firewall.RuleGetResponseData, error) {
	if err, ok := m.getRuleErrors[pos]; ok {
		return nil, err
	}

	return &firewall.RuleGetResponseData{Action: "ACCEPT", Type: "in"}, nil
}

func (m *deleteTestMockAPI) DeleteRule(_ context.Context, pos int) error {
	m.deletedPos = append(m.deletedPos, pos)
	return nil
}

func newDeleteTestState(t *testing.T, rules []map[string]any) *schema.ResourceData {
	t.Helper()

	d := Rules().TestResourceData()

	ruleList := make([]any, len(rules))
	for i, r := range rules {
		ruleList[i] = r
	}

	err := d.Set(MkRule, ruleList)
	require.NoError(t, err)

	return d
}

func ruleState(pos int, comment string) map[string]any {
	return map[string]any{
		mkRulePos:       pos,
		mkRuleAction:    "ACCEPT",
		mkRuleType:      "in",
		mkRuleComment:   comment,
		mkRuleDPort:     "",
		mkRuleProto:     "",
		mkSecurityGroup: "",
		mkRuleDest:      "",
		mkRuleSource:    "",
		mkRuleSPort:     "",
		mkRuleEnabled:   true,
		mkRuleIFace:     "",
		mkRuleLog:       "",
		mkRuleMacro:     "",
	}
}

// TestRulesDeleteSkipsMissingRules verifies that RulesDelete skips rules
// that no longer exist (ErrNoRuleAtPosition) instead of failing.
func TestRulesDeleteSkipsMissingRules(t *testing.T) {
	t.Parallel()

	mock := &deleteTestMockAPI{
		getRuleErrors: map[int]error{
			1: fmt.Errorf("error retrieving firewall rule 1: %w", firewall.ErrNoRuleAtPosition),
		},
	}

	d := newDeleteTestState(t, []map[string]any{
		ruleState(0, "Allow HTTP"),
		ruleState(1, "Allow SSH"),
	})

	diags := RulesDelete(context.Background(), mock, d)
	require.False(t, diags.HasError(), "RulesDelete should not error for missing rules")
	require.Equal(t, []int{0}, mock.deletedPos, "Only existing rule should be deleted")
}

// TestRulesDeletePropagatesAPIErrors verifies that RulesDelete propagates
// real API errors (non-ErrNoRuleAtPosition) instead of silently ignoring them.
func TestRulesDeletePropagatesAPIErrors(t *testing.T) {
	t.Parallel()

	mock := &deleteTestMockAPI{
		getRuleErrors: map[int]error{
			1: fmt.Errorf("500 connection refused"),
		},
	}

	d := newDeleteTestState(t, []map[string]any{
		ruleState(0, "Allow HTTP"),
		ruleState(1, "Allow SSH"),
	})

	diags := RulesDelete(context.Background(), mock, d)
	require.True(t, diags.HasError(), "RulesDelete should propagate non-missing-rule API errors")
	require.Empty(t, mock.deletedPos, "No rules should be deleted when earlier position errors")
}

// TestComputeRuleSignature tests the signature computation for rules.
func TestComputeRuleSignature(t *testing.T) {
	t.Parallel()

	regularRule := map[string]any{
		mkRuleAction:    "ACCEPT",
		mkRuleType:      "in",
		mkRuleDest:      "192.168.1.0/24",
		mkRuleDPort:     "80",
		mkRuleSource:    "10.0.0.0/8",
		mkRuleSPort:     "1024",
		mkRuleProto:     "tcp",
		mkRuleMacro:     "",
		mkRuleIFace:     "net0",
		mkSecurityGroup: "",
		mkRuleComment:   "Test comment", // should be excluded
		mkRuleEnabled:   true,           // should be excluded
		mkRuleLog:       "info",         // should be excluded
	}

	sig1 := computeRuleSignature(regularRule)
	require.NotEmpty(t, sig1)
	require.Contains(t, sig1, "rule:")
	require.Contains(t, sig1, "ACCEPT")
	require.Contains(t, sig1, "in")

	regularRule[mkRuleComment] = "Different comment"
	sig2 := computeRuleSignature(regularRule)
	require.Equal(t, sig1, sig2, "Signature should not change when only comment changes")

	regularRule[mkRuleAction] = "DROP"
	sig3 := computeRuleSignature(regularRule)
	require.NotEqual(t, sig1, sig3, "Signature should change when action changes")

	groupRule := map[string]any{
		mkSecurityGroup: "foo",
		mkRuleIFace:     "net0",
		mkRuleComment:   "Group comment",
		mkRuleEnabled:   true,
		mkRuleAction:    "",
		mkRuleType:      "",
		mkRuleDest:      "",
		mkRuleDPort:     "",
		mkRuleSource:    "",
		mkRuleSPort:     "",
		mkRuleProto:     "",
		mkRuleMacro:     "",
		mkRuleLog:       "",
	}

	groupSig1 := computeRuleSignature(groupRule)
	require.NotEmpty(t, groupSig1)
	require.Contains(t, groupSig1, "group:")
	require.Contains(t, groupSig1, "foo")

	groupRule[mkRuleComment] = "Different group comment"
	groupSig2 := computeRuleSignature(groupRule)
	require.Equal(t, groupSig1, groupSig2, "Group signature should not change when only comment changes")

	groupRule[mkSecurityGroup] = "bar"
	groupSig3 := computeRuleSignature(groupRule)
	require.NotEqual(t, groupSig1, groupSig3, "Group signature should change when group name changes")
}
