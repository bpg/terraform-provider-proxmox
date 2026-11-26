//go:build acceptance || all

//testacc:tier=heavy
//testacc:resource=container

/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package test

import (
	"context"
	"fmt"
	"math/rand"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/stretchr/testify/require"

	"github.com/bpg/terraform-provider-proxmox/proxmox/nodes/containers"
	proxmoxtypes "github.com/bpg/terraform-provider-proxmox/proxmox/types"
)

func TestAccResourceContainerExistingID(t *testing.T) {
	t.Parallel()

	te := InitEnvironment(t)
	imageFileName := fmt.Sprintf("%d-alpine-3.22-default_20250617_amd64.tar.xz", time.Now().UnixMicro())
	testAccDownloadContainerTemplate(t, te, imageFileName)

	tests := []struct {
		name        string
		needsSource bool
		config      string
	}{
		{"container create", false, `
			resource "proxmox_virtual_environment_container" "test" {
				node_name    = "{{.NodeName}}"
				vm_id        = {{.ExistingID}}
				started      = false
				unprivileged = true

				disk {
					datastore_id = "{{.ContainerDatastoreID}}"
					size         = 4
				}

				operating_system {
					template_file_id = "local:vztmpl/{{.ImageFileName}}"
					type             = "alpine"
				}
			}`},
		{"container clone", true, `
			resource "proxmox_virtual_environment_container" "test" {
				node_name = "{{.NodeName}}"
				vm_id     = {{.ExistingID}}
				started   = false

				clone {
					vm_id = {{.SourceID}}
				}
			}`},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			te := InitEnvironment(t)
			ctx := context.Background()

			sourceID := 100000 + rand.Intn(99999)
			existingID := sourceID + 1

			ids := []int{existingID}
			if tt.needsSource {
				ids = append(ids, sourceID)
			}

			for _, id := range ids {
				body := &containers.CreateRequestBody{
					VMID:                 &id,
					Hostname:             new(fmt.Sprintf("existing-id-%d", id)),
					OSTemplateFileVolume: new("local:vztmpl/" + imageFileName),
					RootFS:               &containers.CustomRootFS{Volume: te.ContainerDatastoreID + ":1"},
					Unprivileged:         new(proxmoxtypes.CustomBool(true)),
				}
				require.NoError(t, te.NodeClient().Container(id).CreateContainer(ctx, body).Err())

				t.Cleanup(func() {
					if result := te.NodeClient().Container(id).DeleteContainer(context.Background(), true, true); result.Err() != nil {
						t.Logf("cleanup warning: failed to delete container %d: %v", id, result.Err())
					}
				})
			}

			existing, err := te.NodeClient().Container(existingID).GetContainer(ctx)
			require.NoError(t, err)

			// Runs after the post-test destroy, which deletes whatever container ended up in the state.
			t.Cleanup(func() {
				after, err := te.NodeClient().Container(existingID).GetContainer(context.Background())
				require.NoError(t, err, "container %d that already held the ID is gone", existingID)
				require.Equal(t, existing, after, "container %d that already held the ID was changed", existingID)
			})

			te.AddTemplateVars(map[string]any{
				"ImageFileName": imageFileName,
				"SourceID":      sourceID,
				"ExistingID":    existingID,
			})

			resource.ParallelTest(t, resource.TestCase{
				ProtoV6ProviderFactories: te.AccProviders,
				Steps: []resource.TestStep{{
					Config:      te.RenderConfig(tt.config),
					ExpectError: existingIDError,
				}},
			})
		})
	}
}
