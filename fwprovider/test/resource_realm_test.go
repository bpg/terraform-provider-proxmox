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

func TestAccResourceRealm(t *testing.T) {
	t.Parallel()

	te := InitEnvironment(t)

	tests := []struct {
		name string
		step []resource.TestStep
	}{
		{"multiline description", []resource.TestStep{{
			Config: te.RenderConfig(`
				 resource "proxmox_virtual_environment_realm" "test_realm1" {
					 realm = "{{.Realm}}"
					 type   = "{{.Type}}"
					 
					 description = <<-EOT
						 my
						 description
						 value
					 EOT
				 }`),
			Check: resource.ComposeTestCheckFunc(
				ResourceAttributes("proxmox_virtual_environment_realm.test_realm1", map[string]string{
					"description": "my\ndescription\nvalue",
				}),
			),
		}}},
		{"single line description", []resource.TestStep{{
			Config: te.RenderConfig(`
				 resource "proxmox_virtual_environment_realm" "test_realm2" {
					 nrealm = "{{.NodeName}}"
					 type   = "{{.Type}}"
					 
					 description = "my description value"
				 }`),
			Check: resource.ComposeTestCheckFunc(
				ResourceAttributes("proxmox_virtual_environment_realm.test_realm2", map[string]string{
					"description": "my description value",
				}),
			),
		}}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			resource.Test(t, resource.TestCase{
				ProtoV6ProviderFactories: te.AccProviders,
				Steps:                    tt.step,
			})
		})
	}
}
