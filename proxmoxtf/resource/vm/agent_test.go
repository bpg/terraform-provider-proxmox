/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package resource

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/stretchr/testify/require"

	"github.com/bpg/terraform-provider-proxmox/proxmox/nodes/vms"
)

// TestVMAgentWaitForIPDisabledSchema verifies the agent.wait_for_ip block exposes a `disabled`
// master switch that defaults to false (waiting stays the default behavior). The switch is phrased
// so its zero value means "do not skip", which keeps VMs upgraded from older provider versions —
// whose stored state lacks the key — querying the agent (see #2928).
func TestVMAgentWaitForIPDisabledSchema(t *testing.T) {
	t.Parallel()

	agentBlock, ok := VM().Schema[mkAgent].Elem.(*schema.Resource)
	require.True(t, ok, "agent block must be a *schema.Resource")

	waitForIP, ok := agentBlock.Schema[mkAgentWaitForIP].Elem.(*schema.Resource)
	require.True(t, ok, "wait_for_ip block must be a *schema.Resource")

	disabled, ok := waitForIP.Schema[mkAgentWaitForIPDisabled]
	require.True(t, ok, "wait_for_ip must expose a `disabled` attribute")
	require.Equal(t, schema.TypeBool, disabled.Type)
	require.Equal(t, false, disabled.Default, "wait_for_ip.disabled must default to false")
}

// TestVMAgentWaitForIPConfigUpgradedStateDoesNotSkip reproduces #2928: a VM whose state was
// written by a provider version predating the wait_for_ip master switch has a wait_for_ip block
// with no explicit switch value. The SDK reads that absent bool back as its zero value, so the
// read path must not interpret it as "skip" — otherwise ipv4_addresses/ipv6_addresses become
// null on refresh for every existing VM.
func TestVMAgentWaitForIPConfigUpgradedStateDoesNotSkip(t *testing.T) {
	t.Parallel()

	// Simulate state stored by a provider version predating the wait_for_ip master switch: the
	// wait_for_ip block exists (ipv4 = true) but the flatmap has no switch key. Building
	// ResourceData straight from the InstanceState mirrors the refresh read path, where Terraform
	// core feeds stored state without applying schema defaults.
	is := &terraform.InstanceState{
		ID: "100",
		Attributes: map[string]string{
			"agent.#":                    "1",
			"agent.0.enabled":            "true",
			"agent.0.wait_for_ip.#":      "1",
			"agent.0.wait_for_ip.0.ipv4": "true",
		},
	}

	d := VM().Data(is)

	cfg := getAgentWaitForIPConfig(d)

	require.NotNil(t, cfg)
	require.False(t, cfg.Skip,
		"agent IP lookup must not be skipped for VMs upgraded from a provider without the wait_for_ip master switch")
	require.True(t, cfg.IPv4)
}

// TestAgentWaitForIPConfigFromBlock verifies the wait_for_ip block is translated into the
// right WaitForIPConfig, including the new disable path.
func TestAgentWaitForIPConfigFromBlock(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		block map[string]any
		want  *vms.WaitForIPConfig
	}{
		{
			name:  "disabled true -> skip",
			block: map[string]any{mkAgentWaitForIPDisabled: true},
			want:  &vms.WaitForIPConfig{Skip: true},
		},
		{
			name:  "disabled false, no family -> wait for any (nil)",
			block: map[string]any{mkAgentWaitForIPDisabled: false},
			want:  nil,
		},
		{
			name:  "disabled false, ipv4 -> wait for ipv4",
			block: map[string]any{mkAgentWaitForIPDisabled: false, mkAgentWaitForIPIPv4: true},
			want:  &vms.WaitForIPConfig{IPv4: true},
		},
		{
			name:  "disabled true wins over family",
			block: map[string]any{mkAgentWaitForIPDisabled: true, mkAgentWaitForIPIPv4: true},
			want:  &vms.WaitForIPConfig{Skip: true},
		},
		{
			name:  "disabled absent (zero value) defaults to wait -> #2928 upgraded state",
			block: map[string]any{mkAgentWaitForIPIPv4: true},
			want:  &vms.WaitForIPConfig{IPv4: true},
		},
		{
			name:  "empty block waits for any (nil)",
			block: map[string]any{},
			want:  nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			require.Equal(t, tt.want, agentWaitForIPConfigFromBlock(tt.block))
		})
	}
}
