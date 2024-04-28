/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package tests

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccResourceTime(t *testing.T) {
	t.Parallel()

	te := initTestEnvironment(t)

	tests := []struct {
		name  string
		steps []resource.TestStep
	}{
		{"change timezone", []resource.TestStep{
			{
				Config: te.renderConfig(`resource "proxmox_virtual_environment_time" "node_time" {
				  node_name = "{{.NodeName}}"
				  time_zone = "America/New_York"
				}`),
				Check: testResourceAttributes("proxmox_virtual_environment_time.node_time", map[string]string{
					"time_zone": "America/New_York",
				}),
			},
			{
				Config: te.renderConfig(`resource "proxmox_virtual_environment_time" "node_time" {
				  node_name = "{{.NodeName}}"
				  time_zone = "UTC"
				}`),
				Check: testResourceAttributes("proxmox_virtual_environment_time.node_time", map[string]string{
					"time_zone": "UTC",
				}),
			},
			{
				Config: te.renderConfig(`resource "proxmox_virtual_environment_time" "node_time" {
				  node_name = "{{.NodeName}}"
				  time_zone = "UTC"
				}`),
				Check: testResourceAttributes("proxmox_virtual_environment_time.node_time", map[string]string{
					"time_zone": "UTC",
				}),
			},
		}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resource.Test(t, resource.TestCase{
				ProtoV6ProviderFactories: te.accProviders,
				Steps:                    tt.steps,
			})
		})
	}
}
