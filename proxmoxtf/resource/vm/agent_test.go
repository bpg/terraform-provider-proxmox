/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package resource

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/require"

	"github.com/bpg/terraform-provider-proxmox/proxmox/nodes/vms"
)

// TestVMAgentWaitForIPEnabledSchema verifies the agent.wait_for_ip block exposes an
// `enabled` master switch that defaults to true (waiting stays the default behavior).
func TestVMAgentWaitForIPEnabledSchema(t *testing.T) {
	t.Parallel()

	agentBlock, ok := VM().Schema[mkAgent].Elem.(*schema.Resource)
	require.True(t, ok, "agent block must be a *schema.Resource")

	waitForIP, ok := agentBlock.Schema[mkAgentWaitForIP].Elem.(*schema.Resource)
	require.True(t, ok, "wait_for_ip block must be a *schema.Resource")

	enabled, ok := waitForIP.Schema[mkAgentWaitForIPEnabled]
	require.True(t, ok, "wait_for_ip must expose an `enabled` attribute")
	require.Equal(t, schema.TypeBool, enabled.Type)
	require.Equal(t, true, enabled.Default, "wait_for_ip.enabled must default to true")
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
			name:  "enabled false -> skip",
			block: map[string]any{mkAgentWaitForIPEnabled: false},
			want:  &vms.WaitForIPConfig{Skip: true},
		},
		{
			name:  "enabled true, no family -> wait for any (nil)",
			block: map[string]any{mkAgentWaitForIPEnabled: true},
			want:  nil,
		},
		{
			name:  "enabled true, ipv4 -> wait for ipv4",
			block: map[string]any{mkAgentWaitForIPEnabled: true, mkAgentWaitForIPIPv4: true},
			want:  &vms.WaitForIPConfig{IPv4: true},
		},
		{
			name:  "enabled false wins over family",
			block: map[string]any{mkAgentWaitForIPEnabled: false, mkAgentWaitForIPIPv4: true},
			want:  &vms.WaitForIPConfig{Skip: true},
		},
		{
			name:  "enabled absent defaults to wait (backward compat)",
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
