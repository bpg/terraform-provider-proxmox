//go:build acceptance || all

/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccResourceTime(t *testing.T) {
	t.Parallel()

	te := InitEnvironment(t)

	tests := []struct {
		name  string
		steps []resource.TestStep
	}{
		{"change timezone", []resource.TestStep{
			{
				Config: te.RenderConfig(`resource "proxmox_virtual_environment_time" "node_time" {
				  node_name = "{{.NodeName}}"
				  time_zone = "America/New_York"
				}`),
				Check: ResourceAttributes("proxmox_virtual_environment_time.node_time", map[string]string{
					"time_zone": "America/New_York",
				}),
			},
			{
				Config: te.RenderConfig(`resource "proxmox_virtual_environment_time" "node_time" {
				  node_name = "{{.NodeName}}"
				  time_zone = "UTC"
				}`),
				Check: ResourceAttributes("proxmox_virtual_environment_time.node_time", map[string]string{
					"time_zone": "UTC",
				}),
			},
			{
				Config: te.RenderConfig(`resource "proxmox_virtual_environment_time" "node_time" {
				  node_name = "{{.NodeName}}"
				  time_zone = "UTC"
				}`),
				Check: ResourceAttributes("proxmox_virtual_environment_time.node_time", map[string]string{
					"time_zone": "UTC",
				}),
			},
		}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resource.Test(t, resource.TestCase{
				ProtoV6ProviderFactories: te.AccProviders,
				Steps:                    tt.steps,
			})
		})
	}
}
