/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package controller

import (
	"context"
	"encoding/json"
	"sort"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/bpg/terraform-provider-proxmox/fwprovider/types/stringset"
	"github.com/bpg/terraform-provider-proxmox/proxmox/cluster/sdn/controllers"
	proxmoxtypes "github.com/bpg/terraform-provider-proxmox/proxmox/types"
)

func sortedPeers(t *testing.T, v stringset.Value) []string {
	t.Helper()

	diags := &diag.Diagnostics{}
	peers := v.ValueList(context.Background(), diags)
	require.False(t, diags.HasError())

	sort.Strings(peers)

	return peers
}

func TestFromAPI_BGP(t *testing.T) {
	t.Parallel()

	peers := proxmoxtypes.CustomCommaSeparatedList([]string{"10.0.0.1"})
	data := &controllers.ControllerData{
		Controller: controllers.Controller{
			ASNumber:            new(int64(65000)),
			Nodes:               new(proxmoxtypes.CustomCommaSeparatedList([]string{"pve-01"})),
			Peers:               &peers,
			EBGP:                proxmoxtypes.CustomBoolPtr(new(true)),
			EBPGMultiHop:        proxmoxtypes.CustomInt64OptionPtr(new(int64(2))),
			Loopback:            new("lo"),
			BgpMultiPathAsRelax: proxmoxtypes.CustomBoolPtr(new(true)),
		},
		Digest: new("digest-1"),
	}

	m := &bgpModel{}
	diags := &diag.Diagnostics{}
	m.fromAPI("ctrlB", data, diags)
	require.False(t, diags.HasError())

	assert.Equal(t, "ctrlB", m.ID.ValueString())
	assert.Equal(t, types.StringValue("digest-1"), m.Digest)
	assert.Equal(t, types.Int64Value(65000), m.ASNumber)
	assert.Equal(t, types.StringValue("pve-01"), m.NodeName)
	assert.Equal(t, []string{"10.0.0.1"}, sortedPeers(t, m.Peers))
	assert.Equal(t, types.BoolValue(true), m.EBGP)
	assert.Equal(t, types.Int64Value(2), m.EBGPMultiHop)
	assert.Equal(t, types.StringValue("lo"), m.Loopback)
	assert.Equal(t, types.BoolValue(true), m.MultiPathRelax)
}

func TestFromAPI_BGP_PrefersPending(t *testing.T) {
	t.Parallel()

	peers := proxmoxtypes.CustomCommaSeparatedList([]string{"10.0.0.1"})
	pendingPeers := proxmoxtypes.CustomCommaSeparatedList([]string{"10.0.0.2"})
	pendingDeletedNodes := proxmoxtypes.CustomCommaSeparatedList([]string{"deleted"})

	data := &controllers.ControllerData{
		Controller: controllers.Controller{
			ASNumber: new(int64(65000)),
			Nodes:    new(proxmoxtypes.CustomCommaSeparatedList([]string{"pve-01"})),
			Peers:    &peers,
		},
		Digest: new("digest-1"),
		Pending: &controllers.Controller{
			ASNumber:     new(int64(65100)),
			Nodes:        &pendingDeletedNodes,
			Peers:        &pendingPeers,
			EBGP:         proxmoxtypes.CustomBoolPtr(new(true)),
			EBPGMultiHop: proxmoxtypes.CustomInt64OptionPtr(new(int64(3))),
			Loopback:     new("deleted"),
		},
	}

	m := &bgpModel{}
	diags := &diag.Diagnostics{}
	m.fromAPI("ctrlB", data, diags)
	require.False(t, diags.HasError())

	assert.Equal(t, types.Int64Value(65100), m.ASNumber)
	assert.Equal(t, types.StringNull(), m.NodeName)
	assert.Equal(t, []string{"10.0.0.2"}, sortedPeers(t, m.Peers))
	assert.Equal(t, types.BoolValue(true), m.EBGP)
	assert.Equal(t, types.Int64Value(3), m.EBGPMultiHop)
	assert.Equal(t, types.StringNull(), m.Loopback)
}

// TestFromAPI_BGP_PendingDeletedMarkers reproduces the acceptance-test
// scenario of removing optional BGP attributes: PVE keeps the old values in
// the running configuration and writes the literal "deleted" string for the
// removed keys in the pending block. Removed attributes must refresh to null
// (bools to their false default), and the old running values must not leak through.
func TestFromAPI_BGP_PendingDeletedMarkers(t *testing.T) {
	t.Parallel()

	payload := `{
		"controller": "ctrlB",
		"digest": "digest-1",
		"asn": 65000,
		"node": "pve-01",
		"peers": "10.0.0.1",
		"ebgp": 1,
		"ebgp-multihop": 3,
		"loopback": "lo",
		"bgp-multipath-as-path-relax": 1,
		"pending": {
			"ebgp": "deleted",
			"ebgp-multihop": "deleted",
			"loopback": "deleted",
			"bgp-multipath-as-path-relax": "deleted",
			"node": "deleted",
			"peers": "deleted"
		}
	}`

	var data controllers.ControllerData
	require.NoError(t, json.Unmarshal([]byte(payload), &data))

	m := &bgpModel{}
	diags := &diag.Diagnostics{}
	m.fromAPI("ctrlB", &data, diags)
	require.False(t, diags.HasError())

	assert.Equal(t, types.Int64Value(65000), m.ASNumber)
	assert.Equal(t, types.StringNull(), m.NodeName)
	assert.Equal(t, stringset.NullValue(), m.Peers)
	assert.Equal(t, types.BoolValue(false), m.EBGP)
	assert.Equal(t, types.Int64Null(), m.EBGPMultiHop)
	assert.Equal(t, types.StringNull(), m.Loopback)
	assert.Equal(t, types.BoolValue(false), m.MultiPathRelax)
}

func TestFromAPI_BGP_PendingZeroMultiHopStaysAbsent(t *testing.T) {
	t.Parallel()

	data := &controllers.ControllerData{
		Controller: controllers.Controller{
			ASNumber: new(int64(65000)),
		},
		Pending: &controllers.Controller{
			EBPGMultiHop: proxmoxtypes.CustomInt64OptionPtr(new(int64(0))),
		},
	}

	m := &bgpModel{}
	diags := &diag.Diagnostics{}
	m.fromAPI("ctrlB", data, diags)
	require.False(t, diags.HasError())

	assert.Equal(t, types.Int64Null(), m.EBGPMultiHop)
}

func TestFromAPIForDatasource_BGP(t *testing.T) {
	t.Parallel()

	data := &controllers.ControllerData{
		Controller: controllers.Controller{
			ASNumber: new(int64(65000)),
			Nodes:    new(proxmoxtypes.CustomCommaSeparatedList([]string{"pve-01"})),
			EBGP:     proxmoxtypes.CustomBoolPtr(new(true)),
		},
		Digest: new("digest-1"),
	}

	m := &bgpModel{}
	diags := &diag.Diagnostics{}
	m.fromAPIForDatasource("ctrlB", data, diags)
	require.False(t, diags.HasError())

	assert.Equal(t, "ctrlB", m.ID.ValueString())
	assert.Equal(t, types.StringValue("digest-1"), m.Digest)
	assert.Equal(t, types.Int64Value(65000), m.ASNumber)
	assert.Equal(t, types.StringValue("pve-01"), m.NodeName)
	assert.Empty(t, sortedPeers(t, m.Peers))
	assert.Equal(t, types.BoolValue(true), m.EBGP)
	assert.Equal(t, types.Int64Value(0), m.EBGPMultiHop)
	assert.Equal(t, types.StringValue(""), m.Loopback)
	assert.Equal(t, types.BoolValue(false), m.MultiPathRelax)
}

func TestToAPI_BGP(t *testing.T) {
	t.Parallel()

	m := &bgpModel{
		genericModel:   genericModel{ID: types.StringValue("ctrlB")},
		ASNumber:       types.Int64Value(65000),
		NodeName:       types.StringValue("pve-01"),
		EBGP:           types.BoolValue(true),
		EBGPMultiHop:   types.Int64Value(2),
		Loopback:       types.StringValue("lo"),
		MultiPathRelax: types.BoolValue(true),
	}

	diags := &diag.Diagnostics{}
	m.Peers = stringset.NewValueList([]string{"10.0.0.2", "10.0.0.1"}, diags)
	require.False(t, diags.HasError())

	data := m.toAPI(context.Background(), diags)
	require.False(t, diags.HasError())

	assert.Equal(t, "ctrlB", data.ID)
	assert.Equal(t, new(int64(65000)), data.ASNumber)
	assert.Equal(t, proxmoxtypes.CustomCommaSeparatedList([]string{"pve-01"}), *data.Nodes)
	assert.Nil(t, data.Type)

	require.NotNil(t, data.Peers)
	assert.Equal(t, []string{"10.0.0.1", "10.0.0.2"}, sorted(*data.Peers))

	assert.Equal(t, proxmoxtypes.CustomBool(true), *data.EBGP)
	assert.Equal(t, proxmoxtypes.CustomInt64Option(2), *data.EBPGMultiHop)
	assert.Equal(t, "lo", *data.Loopback)
	assert.Equal(t, proxmoxtypes.CustomBool(true), *data.BgpMultiPathAsRelax)
}

func TestCheckDeletedBGPFields(t *testing.T) {
	t.Parallel()

	fullState := &bgpModel{
		EBGP:           types.BoolValue(true),
		EBGPMultiHop:   types.Int64Value(2),
		Loopback:       types.StringValue("lo"),
		MultiPathRelax: types.BoolValue(true),
	}

	emptyPlan := &bgpModel{
		EBGP:           types.BoolNull(),
		EBGPMultiHop:   types.Int64Null(),
		Loopback:       types.StringNull(),
		MultiPathRelax: types.BoolNull(),
	}

	assert.Equal(t,
		[]string{"ebgp", "ebgp-multihop", "loopback", "bgp-multipath-as-path-relax"},
		checkDeletedBGPFields(emptyPlan, fullState),
	)

	assert.Empty(t, checkDeletedBGPFields(fullState, fullState))
	assert.Empty(t, checkDeletedBGPFields(emptyPlan, emptyPlan))
}

func TestEBGPMultiHopRequiresEBGP(t *testing.T) {
	t.Parallel()

	assert.False(t, ebgpMultiHopRequiresEBGP(types.BoolValue(true), types.Int64Null()))
	assert.True(t, ebgpMultiHopRequiresEBGP(types.BoolNull(), types.Int64Value(2)))
	assert.True(t, ebgpMultiHopRequiresEBGP(types.BoolValue(false), types.Int64Value(2)))
	assert.False(t, ebgpMultiHopRequiresEBGP(types.BoolValue(true), types.Int64Value(2)))
	assert.False(t, ebgpMultiHopRequiresEBGP(types.BoolUnknown(), types.Int64Value(2)))
}

func sorted(values []string) []string {
	result := make([]string, len(values))
	copy(result, values)
	sort.Strings(result)

	return result
}
