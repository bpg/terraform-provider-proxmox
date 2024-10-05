/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package fwprovider_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/bpg/terraform-provider-proxmox/fwprovider/test"
	"github.com/bpg/terraform-provider-proxmox/proxmox/cluster"
	"github.com/bpg/terraform-provider-proxmox/proxmox/helpers/ptr"
	"github.com/bpg/terraform-provider-proxmox/proxmox/nodes/vms"
	"github.com/bpg/terraform-provider-proxmox/utils"
)

func TestIDGenerator_Sequence(t *testing.T) {
	t.Parallel()

	const numIDs = 10

	if utils.GetAnyStringEnv("TF_ACC") == "" {
		t.Skip("Acceptance tests are disabled")
	}

	te := test.InitEnvironment(t)

	ctx := context.Background()

	gen := cluster.NewIDGenerator(te.ClusterClient(), cluster.IDGeneratorConfig{RandomIDs: false})
	firstID, err := gen.NextID(ctx)
	require.NoError(t, err)

	busyID := firstID + 5

	_, err = te.ClusterClient().GetNextID(ctx, ptr.Ptr(busyID))
	require.NoError(t, err, "the VM ID %d should be available", busyID)

	err = te.NodeClient().VM(0).CreateVM(ctx, &vms.CreateRequestBody{VMID: busyID})
	require.NoError(t, err, "failed to create VM %d", busyID)

	t.Cleanup(func() {
		err = te.NodeClient().VM(busyID).DeleteVM(ctx)
		require.NoError(t, err, "failed to delete VM %d", busyID)
	})

	ids := make([]int, numIDs)

	t.Cleanup(func() {
		for _, id := range ids {
			if id > 100 {
				_ = te.NodeClient().VM(id).DeleteVM(ctx) //nolint:errcheck
			}
		}
	})

	prevID := firstID

	for i := range numIDs {
		id, err := gen.NextID(ctx)
		require.NoError(t, err)
		err = te.NodeClient().VM(0).CreateVM(ctx, &vms.CreateRequestBody{VMID: id})
		ids[i] = id

		require.NoError(t, err)
		require.Greater(t, id, prevID, "the generated ID should be greater than the previous one")

		prevID = id
	}
}

func TestIDGenerator_Random(t *testing.T) {
	t.Parallel()

	const (
		numIDs       = 7
		randomIDStat = 1000
		randomIDEnd  = 1010
	)

	if utils.GetAnyStringEnv("TF_ACC") == "" {
		t.Skip("Acceptance tests are disabled")
	}

	te := test.InitEnvironment(t)

	ctx := context.Background()

	gen := cluster.NewIDGenerator(te.ClusterClient(), cluster.IDGeneratorConfig{RandomIDs: true, RandomIDStat: randomIDStat, RandomIDEnd: randomIDEnd})

	ids := make([]int, numIDs)

	t.Cleanup(func() {
		for _, id := range ids {
			if id > 100 {
				_ = te.NodeClient().VM(id).DeleteVM(ctx) //nolint:errcheck
			}
		}
	})

	for i := range numIDs {
		id, err := gen.NextID(ctx)
		require.NoError(t, err)
		err = te.NodeClient().VM(0).CreateVM(ctx, &vms.CreateRequestBody{VMID: id})
		ids[i] = id

		require.NoError(t, err)
	}
}
