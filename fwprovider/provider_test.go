//go:build acceptance || all

/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package fwprovider_test

import (
	"context"
	"regexp"
	"sync"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/bpg/terraform-provider-proxmox/fwprovider/test"
	"github.com/bpg/terraform-provider-proxmox/proxmox/cluster"
	"github.com/bpg/terraform-provider-proxmox/proxmox/helpers/ptr"
	"github.com/bpg/terraform-provider-proxmox/proxmox/nodes/vms"
	"github.com/bpg/terraform-provider-proxmox/utils"
)

func TestIDGenerator_Sequence(t *testing.T) {
	t.Parallel()

	const (
		numIDs     = 10
		numBusyIDs = 30
	)

	if utils.GetAnyStringEnv("TF_ACC") == "" {
		t.Skip("Acceptance tests are disabled")
	}

	te := test.InitEnvironment(t)

	ctx := context.Background()

	gen := cluster.NewIDGenerator(te.ClusterClient(), cluster.IDGeneratorConfig{RandomIDs: false})
	firstID, err := gen.NextID(ctx)
	require.NoError(t, err)

	firstBusyID := firstID + 5

	_, err = te.ClusterClient().GetNextID(ctx, ptr.Ptr(firstBusyID))
	require.NoError(t, err, "the VM ID %d should be available", firstBusyID)

	for i := range numBusyIDs {
		busyID := firstBusyID + i
		err = te.NodeClient().VM(0).CreateVM(ctx, &vms.CreateRequestBody{VMID: busyID})
		require.NoError(t, err, "failed to create VM %d", busyID)
	}

	t.Cleanup(func() {
		var wg sync.WaitGroup

		for i := range numBusyIDs {
			wg.Add(1)

			go func() {
				defer wg.Done()

				busyID := firstBusyID + i
				err = te.NodeClient().VM(busyID).DeleteVM(ctx)
				assert.NoError(t, err, "failed to delete VM %d", busyID)
			}()
		}

		wg.Wait()
	})

	ids := make([]int, numIDs)

	t.Cleanup(func() {
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

	for i := range numIDs {
		wg.Add(1)

		go func() {
			defer wg.Done()

			id, err := gen.NextID(ctx)
			if err == nil {
				err = te.NodeClient().VM(0).CreateVM(ctx, &vms.CreateRequestBody{VMID: id})
				ids[i] = id
			}

			assert.NoError(t, err)
		}()
	}

	wg.Wait()
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

func TestProviderAuth(t *testing.T) {
	if utils.GetAnyStringEnv("TF_ACC") == "" {
		t.Skip("Acceptance tests are disabled")
	}

	te := test.InitEnvironment(t)

	tests := []struct {
		name  string
		steps []resource.TestStep
	}{
		{"no credentials", []resource.TestStep{{
			Config: `
					provider "proxmox" {
						api_token              = ""
						username               = ""
					}
					data "proxmox_virtual_environment_version" "test" {}
					`,
			ExpectError: regexp.MustCompile(`must provide either username and password, an API token, or a ticket`),
		}}},
		{"invalid api token", []resource.TestStep{{
			Config: `
					provider "proxmox" {
						api_token = "invalid-token"
					}
					data "proxmox_virtual_environment_version" "test" {}
					`,
			ExpectError: regexp.MustCompile(`the API token must be in the format 'USER@REALM!TOKENID=UUID'`),
		}}},
		{"invalid username", []resource.TestStep{{
			Config: `
					provider "proxmox" {
						api_token              = ""
						username               = "root"
					}
					data "proxmox_virtual_environment_version" "test" {}
					`,
			ExpectError: regexp.MustCompile(`username must end with '@pve' or '@pam'`),
		}}},
		{"missing password", []resource.TestStep{{
			Config: `
					provider "proxmox" {
						api_token              = ""
						username               = "root@pam"
						password			   = ""
					}
					data "proxmox_virtual_environment_version" "test" {}
					`,
			ExpectError: regexp.MustCompile(`must provide either username and password, an API token, or a ticket`),
		}}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resource.ParallelTest(t, resource.TestCase{
				ProtoV6ProviderFactories: te.AccProviders,
				Steps:                    tt.steps,
			})
		})
	}
}
