//go:build acceptance || all

//testacc:tier=light
//testacc:resource=misc

/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package config_test

import (
	"context"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/stretchr/testify/require"

	"github.com/bpg/terraform-provider-proxmox/fwprovider/test"
	"github.com/bpg/terraform-provider-proxmox/proxmox/nodes"
)

func TestAccResourceNodeConfig(t *testing.T) {
	t.Parallel()

	te := test.InitEnvironment(t)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: te.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: te.RenderConfig(`
					resource "proxmox_node_config" "test" {
						node_name   = "{{.NodeName}}"
						description = "test notes"
					}

					data "proxmox_node_config" "test" {
						node_name = proxmox_node_config.test.node_name
					}
				`),
				Check: resource.ComposeTestCheckFunc(
					test.ResourceAttributes("proxmox_node_config.test", map[string]string{
						"description": "test notes",
					}),
					test.ResourceAttributes("data.proxmox_node_config.test", map[string]string{
						"description": "test notes",
					}),
				),
			},
			// Raw heredoc ends in \n
			// validator must reject it and guide the user to trimspace()
			{
				Config: te.RenderConfig(`
							resource "proxmox_node_config" "test" {
								node_name   = "{{.NodeName}}"
								description = <<-EOT
								Multi-line notes
								EOT
							}`),
				ExpectError: regexp.MustCompile(`must not end with a newline`),
			},
			// trimspace() strips the trailing newline before validation
			// correct heredoc idiom
			{
				Config: te.RenderConfig(`
							resource "proxmox_node_config" "test" {
								node_name   = "{{.NodeName}}"
								description = trimspace(<<-EOT
									Multi-line notes
								EOT
								)
							}`),
				Check: resource.ComposeTestCheckFunc(
					test.ResourceAttributes("proxmox_node_config.test", map[string]string{
						"description": "Multi-line notes",
					}),
				),
			},
			{
				Config: te.RenderConfig(`
					resource "proxmox_node_config" "test" {
						node_name   = "{{.NodeName}}"
						description = "updated notes"
					}
				`),
				Check: resource.ComposeTestCheckFunc(
					test.ResourceAttributes("proxmox_node_config.test", map[string]string{
						"description": "updated notes",
					}),
				),
			},
			{
				Config: te.RenderConfig(`
					resource "proxmox_node_config" "test" {
						node_name = "{{.NodeName}}"
					}
				`),
				Check: resource.ComposeTestCheckFunc(
					test.NoResourceAttributesSet("proxmox_node_config.test", []string{"description"}),
				),
			},
			{
				ResourceName:      "proxmox_node_config.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

// TestAccResourceNodeConfigStartAllOnbootDelay verifies the resource can be read on a node that
// has startall-onboot-delay explicitly set. PVE serializes that field as a string, which crashed
// the GET response decode before the fix.
func TestAccResourceNodeConfigStartAllOnbootDelay(t *testing.T) {
	t.Parallel()

	te := test.InitEnvironment(t)

	ctx := context.Background()

	// Set startall-onboot-delay out-of-band so PVE returns it as a string on the next GET.
	err := te.NodeClient().UpdateConfig(ctx, &nodes.ConfigUpdateRequestBody{
		StartAllOnbootDelay: new(15),
	})
	require.NoError(t, err)

	t.Cleanup(func() {
		err := te.NodeClient().UpdateConfig(ctx, &nodes.ConfigUpdateRequestBody{
			Delete: []string{"startall-onboot-delay"},
		})
		if err != nil {
			t.Logf("cleanup: failed to delete startall-onboot-delay: %v", err)
		}
	})

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: te.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: te.RenderConfig(`
					resource "proxmox_node_config" "test" {
						node_name   = "{{.NodeName}}"
						description = "test notes"
					}

					data "proxmox_node_config" "test" {
						node_name = proxmox_node_config.test.node_name
					}
				`),
				Check: resource.ComposeTestCheckFunc(
					test.ResourceAttributes("proxmox_node_config.test", map[string]string{
						"description": "test notes",
					}),
					test.ResourceAttributes("data.proxmox_node_config.test", map[string]string{
						"description": "test notes",
					}),
				),
			},
		},
	})
}

func TestAccResourceNodeConfigEmptyDescription(t *testing.T) {
	t.Parallel()

	te := test.InitEnvironment(t)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: te.AccProviders,
		Steps: []resource.TestStep{
			// empty string must be rejected by the validator to not cause a plan/state mismatch.
			{
				Config: te.RenderConfig(`
					resource "proxmox_node_config" "test_empty" {
						node_name   = "{{.NodeName}}"
						description = ""
					}
				`),
				ExpectError: regexp.MustCompile(`string length must be at least 1`),
			},
		},
	})
}
