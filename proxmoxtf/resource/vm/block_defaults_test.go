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

	"github.com/bpg/terraform-provider-proxmox/proxmoxtf/resource/vm/disk"
)

// Block-typed attributes must not carry a DefaultFunc. Terraform never applies
// one to a block, but Pulumi's Terraform bridge does, and a default the read
// never reports back becomes a permanent diff for every Pulumi user.
func TestBlockAttributesHaveNoDefaultFunc(t *testing.T) {
	t.Parallel()

	s := VM().Schema
	for _, key := range []string{mkAgent, mkCDROM, mkWatchdog} {
		require.Nil(t, s[key].DefaultFunc, key)
	}

	speed := disk.Schema()[disk.MkDisk].Elem.(*schema.Resource).Schema["speed"]
	require.Nil(t, speed.DefaultFunc, "disk.speed")
}

// The provider's own view of an absent block is unchanged: it still falls
// back to the same defaults the DefaultFuncs used to supply.
func TestGetBlockOrDefault(t *testing.T) {
	t.Parallel()

	r := VM()

	absent := schema.TestResourceDataRaw(t, r.Schema, map[string]any{"node_name": "pve"})

	agent, err := getBlockOrDefault(r, absent, mkAgent)
	require.NoError(t, err)
	require.Equal(t, agentBlockDefault()[0], agent)

	cdrom, err := getBlockOrDefault(r, absent, mkCDROM)
	require.NoError(t, err)
	require.Equal(t, cdromBlockDefault()[0], cdrom)

	watchdog, err := getBlockOrDefault(r, absent, mkWatchdog)
	require.NoError(t, err)
	require.Equal(t, watchdogBlockDefault()[0], watchdog)

	set := schema.TestResourceDataRaw(t, r.Schema, map[string]any{
		"node_name": "pve",
		mkAgent:     []any{map[string]any{mkAgentEnabled: true}},
	})

	agent, err = getBlockOrDefault(r, set, mkAgent)
	require.NoError(t, err)
	require.Equal(t, true, agent[mkAgentEnabled])
}
