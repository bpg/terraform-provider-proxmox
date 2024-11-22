/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package vm_test

import (
	"context"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/bpg/terraform-provider-proxmox/fwprovider/test"
	"github.com/bpg/terraform-provider-proxmox/proxmox/cluster"
	"github.com/bpg/terraform-provider-proxmox/proxmox/nodes/vms"
	"github.com/bpg/terraform-provider-proxmox/utils"
)

func TestBatchCreate(t *testing.T) {
	t.Parallel()

	const (
		numVMs = 30
	)

	if utils.GetAnyStringEnv("TF_ACC") == "" {
		t.Skip("Acceptance tests are disabled")
	}

	te := test.InitEnvironment(t)

	ctx := context.Background()

	gen := cluster.NewIDGenerator(te.ClusterClient(), cluster.IDGeneratorConfig{RandomIDs: false})

	sourceID, err := gen.NextID(ctx)
	require.NoError(t, err)

	err = te.NodeClient().VM(0).CreateVM(ctx, &vms.CreateRequestBody{VMID: sourceID})

	require.NoError(t, err, "failed to create VM %d", sourceID)

	ids := make([]int, numVMs)

	t.Cleanup(func() {
		_ = te.NodeClient().VM(sourceID).DeleteVM(ctx) //nolint:errcheck

		var wg sync.WaitGroup
		for _, id := range ids {
			wg.Add(1)

			go func() {
				defer wg.Done()

				if id > 100 {
					_ = te.NodeClient().VM(id).DeleteVM(ctx) //nolint:errcheck
				}
			}()
		}

		wg.Wait()
	})

	var wg sync.WaitGroup

	for i := range numVMs {
		wg.Add(1)

		go func() {
			defer wg.Done()

			id := 999900 + i
			if err == nil {
				err = te.NodeClient().VM(sourceID).CloneVM(ctx, 5, &vms.CloneRequestBody{VMIDNew: id})
				ids[i] = id
			}

			assert.NoError(t, err)
		}()
	}

	wg.Wait()
}
