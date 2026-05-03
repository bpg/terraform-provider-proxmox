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

func TestAccResourceSDNControllerEVPN(t *testing.T) {
	// Cannot run in parallel due to SDN applier functionality affecting global state

	te := InitEnvironment(t)

	tests := []struct {
		name  string
		steps []resource.TestStep
	}{
		{"Update single peer", []resource.TestStep{{
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
			`),
			Check: resource.ComposeTestCheckFunc(
				ResourceAttributes("proxmox_sdn_controller_evpn.controller_evpn", map[string]string{
					"id":      "ctrlE",
					"peers.#": "1",
					"peers.0": "10.0.0.1",
					"asn":     "65000",
				}),
				NoResourceAttributesSet("proxmox_sdn_controller_evpn.controller_evpn", []string{
					"fabric",
				}),
			),
		}, {
			Config: te.RenderConfig(`
				resource "proxmox_sdn_controller_evpn" "controller_evpn" {
				  id  = "ctrlE"
				  peers = ["10.0.0.2"]
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
			`),
			Check: resource.ComposeTestCheckFunc(
				ResourceAttributes("proxmox_sdn_controller_evpn.controller_evpn", map[string]string{
					"id":      "ctrlE",
					"peers.#": "1",
					"peers.0": "10.0.0.2",
					"asn":     "65000",
				}),
				NoResourceAttributesSet("proxmox_sdn_controller_evpn.controller_evpn", []string{
					"fabric",
				}),
			),
			ResourceName:      "proxmox_sdn_controller_evpn.controller_evpn",
			ImportStateId:     "ctrlE",
			ImportState:       true,
			ImportStateVerify: true,
		}}},
		{"Update multiple peers", []resource.TestStep{{
			Config: te.RenderConfig(`
				resource "proxmox_sdn_controller_evpn" "controller_evpn" {
				  id  = "ctrlE"
				  peers = ["10.0.0.1", "10.0.0.2", "10.0.0.3"]
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
			`),

			Check: resource.ComposeTestCheckFunc(
				ResourceAttributes("proxmox_sdn_controller_evpn.controller_evpn", map[string]string{
					"id":      "ctrlE",
					"peers.#": "3",
					"peers.0": "10.0.0.1",
					"peers.1": "10.0.0.2",
					"peers.2": "10.0.0.3",
					"asn":     "65000",
				}),
				NoResourceAttributesSet("proxmox_sdn_controller_evpn.controller_evpn", []string{
					"fabric",
				}),
			),
		}, {
			Config: te.RenderConfig(`
				resource "proxmox_sdn_controller_evpn" "controller_evpn" {
				  id  = "ctrlE"
				  peers = ["10.0.0.1", "10.0.0.4"]
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
			`),
			Check: resource.ComposeTestCheckFunc(
				ResourceAttributes("proxmox_sdn_controller_evpn.controller_evpn", map[string]string{
					"id":      "ctrlE",
					"peers.#": "2",
					"peers.0": "10.0.0.1",
					"peers.1": "10.0.0.4",
					"asn":     "65000",
				}),
				NoResourceAttributesSet("proxmox_sdn_controller_evpn.controller_evpn", []string{
					"fabric",
				}),
			),
		}}},
		{"Update fabric", []resource.TestStep{{
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
			`),
			Check: resource.ComposeTestCheckFunc(
				ResourceAttributes("proxmox_sdn_controller_evpn.controller_evpn", map[string]string{
					"id":     "ctrlE",
					"fabric": "main",
					"asn":    "65000",
				}),
				NoResourceAttributesSet("proxmox_sdn_controller_evpn.controller_evpn", []string{
					"peers",
				}),
			),
		}, {
			Config: te.RenderConfig(`
				resource "proxmox_sdn_fabric_openfabric" "main" {
				  id = "main"
				  ip_prefix = "10.0.0.0/16"
				  depends_on = [
				    proxmox_sdn_applier.finalizer
				  ]
				}
				resource "proxmox_sdn_fabric_openfabric" "main2" {
				  id = "main2"
				  ip_prefix = "10.1.0.0/16"
				  depends_on = [
				    proxmox_sdn_applier.finalizer
				  ]
				}
				resource "proxmox_sdn_controller_evpn" "controller_evpn" {
				  id  = "ctrlE"
				  asn = 65000
				  fabric = proxmox_sdn_fabric_openfabric.main2.id
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
			`),
			Check: resource.ComposeTestCheckFunc(
				ResourceAttributes("proxmox_sdn_controller_evpn.controller_evpn", map[string]string{
					"id":     "ctrlE",
					"asn":    "65000",
					"fabric": "main2",
				}),
				NoResourceAttributesSet("proxmox_sdn_controller_evpn.controller_evpn", []string{
					"peers",
				}),
			),
		}}},
		{"Change peers to fabric", []resource.TestStep{{
			Config: te.RenderConfig(`
				resource "proxmox_sdn_controller_evpn" "controller_evpn" {
				  id  = "ctrlE"
				  asn = 65000
				  peers = ["10.0.0.1"]
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
			`),
			Check: resource.ComposeTestCheckFunc(
				ResourceAttributes("proxmox_sdn_controller_evpn.controller_evpn", map[string]string{
					"id":      "ctrlE",
					"peers.#": "1",
					"peers.0": "10.0.0.1",
					"asn":     "65000",
				}),
				NoResourceAttributesSet("proxmox_sdn_controller_evpn.controller_evpn", []string{
					"fabric",
				}),
			),
		}, {
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
			`),
			Check: resource.ComposeTestCheckFunc(
				ResourceAttributes("proxmox_sdn_controller_evpn.controller_evpn", map[string]string{
					"id":     "ctrlE",
					"fabric": "main",
					"asn":    "65000",
				}),
				NoResourceAttributesSet("proxmox_sdn_controller_evpn.controller_evpn", []string{
					"peers",
				}),
			),
		}}},
		{"Change fabric to peers", []resource.TestStep{{
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
			`),
			Check: resource.ComposeTestCheckFunc(
				ResourceAttributes("proxmox_sdn_controller_evpn.controller_evpn", map[string]string{
					"id":     "ctrlE",
					"fabric": "main",
					"asn":    "65000",
				}),
				NoResourceAttributesSet("proxmox_sdn_controller_evpn.controller_evpn", []string{
					"peers",
				}),
			),
		}, {
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
				  asn = 65000
				  peers = ["10.0.0.1"]
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
			`),
			Check: resource.ComposeTestCheckFunc(
				ResourceAttributes("proxmox_sdn_controller_evpn.controller_evpn", map[string]string{
					"id":      "ctrlE",
					"peers.#": "1",
					"peers.0": "10.0.0.1",
					"asn":     "65000",
				}),
				NoResourceAttributesSet("proxmox_sdn_controller_evpn.controller_evpn", []string{
					"fabric",
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
