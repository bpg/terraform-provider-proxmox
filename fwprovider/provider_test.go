//go:build acceptance || all

//testacc:tier=heavy
//testacc:resource=vm

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
	"github.com/bpg/terraform-provider-proxmox/proxmox/access"
	"github.com/bpg/terraform-provider-proxmox/proxmox/api"
	"github.com/bpg/terraform-provider-proxmox/proxmox/cluster"
	"github.com/bpg/terraform-provider-proxmox/proxmox/nodes/vms"
	proxmoxtypes "github.com/bpg/terraform-provider-proxmox/proxmox/types"
	"github.com/bpg/terraform-provider-proxmox/utils"
)

func TestAccIDGenerator_Sequence(t *testing.T) {
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

	_, err = te.ClusterClient().GetNextID(ctx, new(firstBusyID))
	require.NoError(t, err, "the VM ID %d should be available", firstBusyID)

	for i := range numBusyIDs {
		busyID := firstBusyID + i
		result := te.NodeClient().VM(0).CreateVM(ctx, &vms.CreateRequestBody{VMID: busyID})
		require.NoError(t, result.Err(), "failed to create VM %d", busyID)
	}

	t.Cleanup(func() {
		var wg sync.WaitGroup

		for i := range numBusyIDs {
			wg.Add(1)

			go func() {
				defer wg.Done()

				busyID := firstBusyID + i
				result := te.NodeClient().VM(busyID).DeleteVM(ctx, true, true)
				assert.NoError(t, result.Err(), "failed to delete VM %d", busyID)
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
					if result := te.NodeClient().VM(id).DeleteVM(ctx, true, true); result.Err() != nil {
						t.Logf("cleanup warning: failed to delete VM %d: %v", id, result.Err())
					}
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
				result := te.NodeClient().VM(0).CreateVM(ctx, &vms.CreateRequestBody{VMID: id})
				ids[i] = id
				err = result.Err()
			}

			assert.NoError(t, err)
		}()
	}

	wg.Wait()
}

func TestAccIDGenerator_SequenceNoReuseBeforeCreate(t *testing.T) {
	if utils.GetAnyStringEnv("TF_ACC") == "" {
		t.Skip("Acceptance tests are disabled")
	}

	t.Setenv("TMPDIR", t.TempDir())

	te := test.InitEnvironment(t)
	ctx := context.Background()

	gen := cluster.NewIDGenerator(te.ClusterClient(), cluster.IDGeneratorConfig{RandomIDs: false})

	firstID, err := gen.NextID(ctx)
	require.NoError(t, err)

	busyIDs := []int{firstID + 1, firstID + 2}
	createBusyVMs(ctx, t, te, busyIDs)

	ids := []int{firstID}

	for range 3 {
		id, err := gen.NextID(ctx)
		require.NoError(t, err)

		ids = append(ids, id)
	}

	seen := map[int]bool{}

	for _, id := range ids {
		assert.False(t, seen[id], "ID %d handed out twice: %v", id, ids)
		assert.NotContains(t, busyIDs, id)

		seen[id] = true
	}
}

func TestAccIDGenerator_SequenceSkipsGuestsHiddenFromToken(t *testing.T) {
	const numBusyIDs = 40

	if utils.GetAnyStringEnv("TF_ACC") == "" {
		t.Skip("Acceptance tests are disabled")
	}

	t.Setenv("TMPDIR", t.TempDir())

	te := test.InitEnvironment(t)
	ctx := context.Background()

	userID := test.SafeResourceName("idgen") + "@pve"
	tokenName := "noacl"

	require.NoError(t, te.AccessClient().CreateUser(ctx, &access.UserCreateRequestBody{ID: userID}))

	t.Cleanup(func() {
		_ = te.AccessClient().DeleteUserToken(context.Background(), userID, tokenName)
		_ = te.AccessClient().DeleteUser(context.Background(), userID)
	})

	token, err := te.AccessClient().CreateUserToken(ctx, userID, tokenName, &access.UserTokenCreateRequestBody{
		PrivSeparate: proxmoxtypes.CustomBool(true).Pointer(),
	})
	require.NoError(t, err)

	creds, err := api.NewCredentials("", "", "", token, "", "")
	require.NoError(t, err)

	conn, err := api.NewConnection(utils.GetAnyStringEnv("PROXMOX_VE_ENDPOINT"), true, "", nil)
	require.NoError(t, err)

	limitedClient, err := api.NewClient(creds, conn)
	require.NoError(t, err)

	gen := cluster.NewIDGenerator(&cluster.Client{Client: limitedClient}, cluster.IDGeneratorConfig{RandomIDs: false})

	firstID, err := gen.NextID(ctx)
	require.NoError(t, err)

	busyIDs := make([]int, numBusyIDs)
	for i := range busyIDs {
		busyIDs[i] = firstID + 1 + i
	}

	createBusyVMs(ctx, t, te, busyIDs)

	id, err := gen.NextID(ctx)
	require.NoError(t, err)
	assert.Greater(t, id, firstID+numBusyIDs)
}

func createBusyVMs(ctx context.Context, t *testing.T, te *test.Environment, ids []int) {
	t.Helper()

	t.Cleanup(func() {
		var wg sync.WaitGroup

		for _, id := range ids {
			wg.Go(func() {
				if result := te.NodeClient().VM(id).DeleteVM(ctx, true, true); result.Err() != nil {
					t.Logf("cleanup warning: failed to delete VM %d: %v", id, result.Err())
				}
			})
		}

		wg.Wait()
	})

	for _, id := range ids {
		result := te.NodeClient().VM(0).CreateVM(ctx, &vms.CreateRequestBody{VMID: id})
		require.NoError(t, result.Err(), "failed to create VM %d", id)
	}
}

func TestAccIDGenerator_Random(t *testing.T) {
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
				if result := te.NodeClient().VM(id).DeleteVM(ctx, true, true); result.Err() != nil {
					t.Logf("cleanup warning: failed to delete VM %d: %v", id, result.Err())
				}
			}
		}
	})

	for i := range numIDs {
		id, err := gen.NextID(ctx)
		require.NoError(t, err)
		result := te.NodeClient().VM(0).CreateVM(ctx, &vms.CreateRequestBody{VMID: id})
		ids[i] = id

		require.NoError(t, result.Err())
	}
}

func TestAccProviderAuth(t *testing.T) {
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
