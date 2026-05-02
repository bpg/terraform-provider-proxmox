//go:build acceptance || all

//testacc:tier=heavy
//testacc:resource=sdn

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

func TestAccDatasourceSDNControllerEVPN(t *testing.T) {
	// Cannot run in parallel due to SDN applier functionality affecting global state

	te := InitEnvironment(t)

	tests := []struct {
		name  string
		steps []resource.TestStep
	}{
		{"Read peers-only controller via datasource", []resource.TestStep{{
			Config: te.RenderConfig(`
				resource "proxmox_sdn_controller_evpn" "controller_evpn" {
				  id  = "ctrlE"
				  peers = ["10.0.0.1"]
				  asn = 65000
				  depends_on = [
				    proxmox_sdn_applier.finalizer
				  ]
				}

				resource "proxmox_sdn_applier" "main" {
				  depends_on = [
				    proxmox_sdn_controller_evpn.controller_evpn
				  ]
				}
				resource "proxmox_sdn_applier" "finalizer" {}

				data "proxmox_sdn_controller_evpn" "controller_evpn" {
				  id = proxmox_sdn_controller_evpn.controller_evpn.id
				}
			`),
			Check: resource.ComposeTestCheckFunc(
				ResourceAttributes("data.proxmox_sdn_controller_evpn.controller_evpn", map[string]string{
					"id":      "ctrlE",
					"peers.#": "1",
					"peers.0": "10.0.0.1",
					"asn":     "65000",
					"fabric":  "",
				}),
			),
		}}},
		{"Read fabric-only controller via datasource", []resource.TestStep{{
			Config: te.RenderConfig(`
				resource "proxmox_sdn_fabric_openfabric" "main" {
				  id = "main"
				  ip_prefix = "10.0.0.0/16"
				  depends_on = [
				    proxmox_sdn_applier.finalizer
				  ]
				}
				resource "proxmox_sdn_controller_evpn" "controller_evpn" {
				  id  = "ctrlE"
				  fabric = proxmox_sdn_fabric_openfabric.main.id
				  asn = 65000
				  depends_on = [
				    proxmox_sdn_applier.finalizer
				  ]
				}

				resource "proxmox_sdn_applier" "main" {
				  depends_on = [
				    proxmox_sdn_controller_evpn.controller_evpn
				  ]
				}
				resource "proxmox_sdn_applier" "finalizer" {}

				data "proxmox_sdn_controller_evpn" "controller_evpn" {
				  id = proxmox_sdn_controller_evpn.controller_evpn.id
				}
			`),
			Check: resource.ComposeTestCheckFunc(
				ResourceAttributes("data.proxmox_sdn_controller_evpn.controller_evpn", map[string]string{
					"id":      "ctrlE",
					"asn":     "65000",
					"fabric":  "main",
					"peers.#": "0",
				}),
			),
		}}},
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
