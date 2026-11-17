//go:build acceptance || all

//testacc:tier=heavy
//testacc:resource=sdn

/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package zone_test

import (
	"context"
	"fmt"
	"net/http"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"

	"github.com/bpg/terraform-provider-proxmox/fwprovider/test"
	"github.com/bpg/terraform-provider-proxmox/proxmox/cluster/sdn"
	proxmoxtypes "github.com/bpg/terraform-provider-proxmox/proxmox/types"
)

func TestAccResourceSDNZoneSimple(t *testing.T) {
	// Cannot run in parallel due to SDN applier functionality affecting global state

	te := test.InitEnvironment(t)

	tests := []struct {
		name  string
		steps []resource.TestStep
	}{
		{"create and update zones", []resource.TestStep{{
			Config: te.RenderConfig(`
				resource "proxmox_sdn_zone_simple" "zone_simple" {
				  id  = "zoneS"
				  nodes = ["{{.NodeName}}"]
				  ipam  = "pve"
				  mtu   = 1496
				  depends_on = [
					proxmox_sdn_applier.finalizer
				  ]
				}

				resource "proxmox_sdn_applier" "zone_simple_applier" {
				  depends_on = [
					proxmox_sdn_zone_simple.zone_simple
				  ]
				}

				resource "proxmox_sdn_applier" "finalizer" {}
			`),
			Check: resource.ComposeTestCheckFunc(
				test.ResourceAttributes("proxmox_sdn_zone_simple.zone_simple", map[string]string{
					"id":      "zoneS",
					"mtu":     "1496",
					"ipam":    "pve",
					"pending": "true",
					"state":   "new",
				}),
			),
		}, {
			Config: te.RenderConfig(`
				resource "proxmox_sdn_zone_simple" "zone_simple" {
				  id  = "zoneS"
				  nodes = ["{{.NodeName}}"]
				  mtu   = 1495
				  depends_on = [
					proxmox_sdn_applier.finalizer
				  ]
				}

				resource "proxmox_sdn_applier" "zone_simple_applier" {
				  depends_on = [
					proxmox_sdn_zone_simple.zone_simple
				  ]
				}

				resource "proxmox_sdn_applier" "finalizer" {}
			`),
			Check: resource.ComposeTestCheckFunc(
				test.ResourceAttributes("proxmox_sdn_zone_simple.zone_simple", map[string]string{
					"id":      "zoneS",
					"mtu":     "1495",
					"pending": "true",
					"state":   "changed",
				}),
				test.NoResourceAttributesSet("proxmox_sdn_zone_simple.zone_simple", []string{
					"ipam",
				}),
			),
			ResourceName: "proxmox_sdn_zone_simple.zone_simple",
			// ImportStateId:     "zoneS",
			// ImportState:       true,
			// ImportStateVerify: true,
		}}},
		{"create zone without mtu", []resource.TestStep{{
			Config: te.RenderConfig(`
				resource "proxmox_sdn_zone_simple" "zone_no_mtu" {
				  id  = "zoneNM"
				  nodes = ["{{.NodeName}}"]
				  depends_on = [
					proxmox_sdn_applier.finalizer
				  ]
				}

				resource "proxmox_sdn_applier" "zone_no_mtu_applier" {
				  depends_on = [
					proxmox_sdn_zone_simple.zone_no_mtu
				  ]
				}

				resource "proxmox_sdn_applier" "finalizer" {}
			`),
			Check: resource.ComposeTestCheckFunc(
				test.ResourceAttributes("proxmox_sdn_zone_simple.zone_no_mtu", map[string]string{
					"id":      "zoneNM",
					"pending": "true",
					"state":   "new",
				}),
				test.NoResourceAttributesSet("proxmox_sdn_zone_simple.zone_no_mtu", []string{
					"mtu",
				}),
			),
		}}},
		{"create zones with empty nodes", []resource.TestStep{{
			Config: te.RenderConfig(`
				resource "proxmox_sdn_zone_simple" "zone_simple2" {
				  id  = "zoneSE"
				  nodes = []
				  mtu   = 1496
				  depends_on = [
					proxmox_sdn_applier.finalizer
				  ]
				}

				resource "proxmox_sdn_applier" "zone_simple2_applier" {
				  depends_on = [
					proxmox_sdn_zone_simple.zone_simple2
				  ]
				}
				  
				resource "proxmox_sdn_applier" "finalizer" {}

			`),
			Check: resource.ComposeTestCheckFunc(
				test.ResourceAttributes("proxmox_sdn_zone_simple.zone_simple2", map[string]string{
					"id":      "zoneSE",
					"mtu":     "1496",
					"pending": "true",
					"nodes.#": "0",
					"state":   "new",
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

func TestAccResourceSDNZoneVLAN(t *testing.T) {
	// Cannot run in parallel due to SDN applier functionality affecting global state

	te := test.InitEnvironment(t)

	tests := []struct {
		name  string
		steps []resource.TestStep
	}{
		{"create and update VLAN zone", []resource.TestStep{{
			Config: te.RenderConfig(`
				resource "proxmox_sdn_zone_vlan" "zone_vlan" {
				  id    = "zoneV"
				  nodes = ["{{.NodeName}}"]
				  mtu   = 1496
				  bridge = "vmbr0"
				}
			`),
			Check: resource.ComposeTestCheckFunc(
				test.ResourceAttributes("proxmox_sdn_zone_vlan.zone_vlan", map[string]string{
					"id":      "zoneV",
					"mtu":     "1496",
					"bridge":  "vmbr0",
					"pending": "true",
					"state":   "new",
				}),
			),
		}, {
			Config: te.RenderConfig(`
				resource "proxmox_sdn_zone_vlan" "zone_vlan" {
				  id    = "zoneV"
				  nodes = ["{{.NodeName}}"]
				  mtu   = 1495
				  bridge = "vmbr0"
				}
			`),
			Check: resource.ComposeTestCheckFunc(
				test.ResourceAttributes("proxmox_sdn_zone_vlan.zone_vlan", map[string]string{
					"id":      "zoneV",
					"mtu":     "1495",
					"bridge":  "vmbr0",
					"pending": "true",
					"state":   "new",
				}),
			),
			ResourceName:      "proxmox_sdn_zone_vlan.zone_vlan",
			ImportStateId:     "zoneV",
			ImportState:       true,
			ImportStateVerify: true,
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

func TestAccResourceSDNZoneVLAN_NoNodes(t *testing.T) {
	// Cannot run in parallel due to SDN applier functionality affecting global state

	te := test.InitEnvironment(t)

	tests := []struct {
		name  string
		steps []resource.TestStep
	}{
		{
			name: "create VLAN zone without nodes",
			steps: []resource.TestStep{{
				Config: te.RenderConfig(`
					resource "proxmox_sdn_zone_vlan" "zone_vlan_no_nodes" {
					  id     = "zoneVNo"
					  bridge = "vmbr0"
					  mtu    = 1496
					}
				`),
				Check: resource.ComposeTestCheckFunc(
					test.ResourceAttributes("proxmox_sdn_zone_vlan.zone_vlan_no_nodes", map[string]string{
						"id":      "zoneVNo",
						"bridge":  "vmbr0",
						"mtu":     "1496",
						"pending": "true",
						"state":   "new",
					}),
				),
			}},
		},
		{
			name: "create VLAN zone with empty nodes list",
			steps: []resource.TestStep{{
				Config: te.RenderConfig(`
					resource "proxmox_sdn_zone_vlan" "zone_vlan_empty_nodes" {
					  id     = "zoneVEm"
					  nodes  = []
					  bridge = "vmbr0"
					  mtu    = 1496
					}
				`),
				Check: resource.ComposeTestCheckFunc(
					test.ResourceAttributes("proxmox_sdn_zone_vlan.zone_vlan_empty_nodes", map[string]string{
						"id":      "zoneVEm",
						"bridge":  "vmbr0",
						"mtu":     "1496",
						"pending": "true",
						"state":   "new",
					}),
				),
			}},
		},
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

func TestAccResourceSDNZoneQinQ(t *testing.T) {
	// Cannot run in parallel due to SDN applier functionality affecting global state

	te := test.InitEnvironment(t)

	tests := []struct {
		name  string
		steps []resource.TestStep
	}{
		{"create and update QinQ zone", []resource.TestStep{{
			Config: te.RenderConfig(`
				resource "proxmox_sdn_zone_qinq" "zone_qinq" {
				  id    = "zoneQ"
				  nodes = ["{{.NodeName}}"]
				  mtu   = 1496
				  bridge = "vmbr0"
				  service_vlan = 100
				  service_vlan_protocol = "802.1ad"
				}
			`),
			Check: resource.ComposeTestCheckFunc(
				test.ResourceAttributes("proxmox_sdn_zone_qinq.zone_qinq", map[string]string{
					"id":                    "zoneQ",
					"mtu":                   "1496",
					"bridge":                "vmbr0",
					"service_vlan":          "100",
					"service_vlan_protocol": "802.1ad",
					"pending":               "true",
					"state":                 "new",
				}),
			),
		}, {
			Config: te.RenderConfig(`
				resource "proxmox_sdn_zone_qinq" "zone_qinq" {
				  id    = "zoneQ"
				  nodes = ["{{.NodeName}}"]
				  mtu   = 1495
				  bridge = "vmbr0"
				  service_vlan = 200
				  service_vlan_protocol = "802.1q"
				}
			`),
			Check: resource.ComposeTestCheckFunc(
				test.ResourceAttributes("proxmox_sdn_zone_qinq.zone_qinq", map[string]string{
					"id":                    "zoneQ",
					"mtu":                   "1495",
					"bridge":                "vmbr0",
					"service_vlan":          "200",
					"service_vlan_protocol": "802.1q",
					"pending":               "true",
					"state":                 "new",
				}),
			),
			ResourceName:      "proxmox_sdn_zone_qinq.zone_qinq",
			ImportStateId:     "zoneQ",
			ImportState:       true,
			ImportStateVerify: true,
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

func TestAccResourceSDNZoneVXLAN(t *testing.T) {
	// Cannot run in parallel due to SDN applier functionality affecting global state

	te := test.InitEnvironment(t)

	tests := []struct {
		name  string
		steps []resource.TestStep
	}{
		{"create and update VXLAN zone", []resource.TestStep{{
			Config: te.RenderConfig(`
				resource "proxmox_sdn_zone_vxlan" "zone_vxlan" {
				  id    = "zoneX"
				  nodes = ["{{.NodeName}}"]
				  mtu   = 1450
				  peers = ["10.0.0.1", "10.0.0.2"]
				}
			`),
			Check: resource.ComposeTestCheckFunc(
				test.ResourceAttributes("proxmox_sdn_zone_vxlan.zone_vxlan", map[string]string{
					"id":      "zoneX",
					"mtu":     "1450",
					"pending": "true",
					"state":   "new",
				}),
			),
		}, {
			Config: te.RenderConfig(`
				resource "proxmox_sdn_zone_vxlan" "zone_vxlan" {
				  id    = "zoneX"
				  nodes = ["{{.NodeName}}"]
				  mtu   = 1440
				  peers = ["10.0.0.3", "10.0.0.4"]
				}
			`),
			Check: resource.ComposeTestCheckFunc(
				test.ResourceAttributes("proxmox_sdn_zone_vxlan.zone_vxlan", map[string]string{
					"id":      "zoneX",
					"mtu":     "1440",
					"pending": "true",
					"state":   "new",
				}),
			),
			ResourceName:      "proxmox_sdn_zone_vxlan.zone_vxlan",
			ImportStateId:     "zoneX",
			ImportState:       true,
			ImportStateVerify: true,
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

func TestAccResourceSDNZoneMTURemoval(t *testing.T) {
	// Cannot run in parallel due to SDN applier functionality affecting global state

	te := test.InitEnvironment(t)

	zoneWithMTU := te.RenderConfig(`
		resource "proxmox_sdn_zone_simple" "zone_mtu" {
		  id    = "zoneRM"
		  nodes = ["{{.NodeName}}"]
		  mtu   = 1496
		  depends_on = [proxmox_sdn_applier.finalizer]
		}

		resource "proxmox_sdn_applier" "apply_mtu" {
		  depends_on = [proxmox_sdn_zone_simple.zone_mtu]
		}

		resource "proxmox_sdn_applier" "finalizer" {}
	`)

	zoneWithoutMTU := `
		resource "proxmox_sdn_zone_simple" "zone_mtu" {
		  id    = "zoneRM"
		  nodes = ["{{.NodeName}}"]
		  depends_on = [proxmox_sdn_applier.finalizer]
		}

		resource "proxmox_sdn_applier" "apply_mtu" {
		  depends_on = [proxmox_sdn_zone_simple.zone_mtu]
		}

		resource "proxmox_sdn_applier" "finalizer" {}
	`

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: te.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: zoneWithMTU,
				Check: resource.ComposeTestCheckFunc(
					test.ResourceAttributes("proxmox_sdn_zone_simple.zone_mtu", map[string]string{
						"mtu": "1496",
					}),
					checkZoneRaw(te, "zoneRM", func(zone map[string]any) error {
						if fmt.Sprint(zone["mtu"]) != "1496" {
							return fmt.Errorf("expected applied mtu 1496, got %v", zone["mtu"])
						}

						return nil
					}),
				),
			},
			{
				// The applier is not re-run, so PVE reports the removal as a pending "deleted" marker.
				Config: te.RenderConfig(zoneWithoutMTU),
				Check: resource.ComposeTestCheckFunc(
					test.NoResourceAttributesSet("proxmox_sdn_zone_simple.zone_mtu", []string{"mtu"}),
					test.ResourceAttributes("proxmox_sdn_zone_simple.zone_mtu", map[string]string{
						"pending": "true",
					}),
					checkZoneRaw(te, "zoneRM", func(zone map[string]any) error {
						pending, _ := zone["pending"].(map[string]any)
						if fmt.Sprint(pending["mtu"]) != "deleted" {
							return fmt.Errorf("expected pending mtu \"deleted\", got %v", pending["mtu"])
						}

						return nil
					}),
				),
			},
			{
				Config: te.RenderConfig(zoneWithoutMTU + `
					resource "proxmox_sdn_applier" "apply_removal" {
					  depends_on = [proxmox_sdn_zone_simple.zone_mtu]
					}
				`),
				Check: checkZoneRaw(te, "zoneRM", func(zone map[string]any) error {
					if v, ok := zone["mtu"]; ok {
						return fmt.Errorf("expected no mtu after apply, got %v", v)
					}

					if v, ok := zone["pending"]; ok {
						return fmt.Errorf("expected no pending changes after apply, got %v", v)
					}

					return nil
				}),
			},
			{
				Config: te.RenderConfig(zoneWithoutMTU + `
					resource "proxmox_sdn_applier" "apply_removal" {
					  depends_on = [proxmox_sdn_zone_simple.zone_mtu]
					}
				`),
				Check: test.NoResourceAttributesSet("proxmox_sdn_zone_simple.zone_mtu", []string{"mtu", "pending"}),
			},
		},
	})
}

func checkZoneRaw(te *test.Environment, id string, check func(zone map[string]any) error) resource.TestCheckFunc {
	return func(_ *terraform.State) error {
		zones := te.ClusterClient().SDNZones()
		resBody := &struct {
			Data map[string]any `json:"data"`
		}{}

		params := &sdn.QueryParams{Pending: proxmoxtypes.CustomBool(true).Pointer()}

		err := zones.DoRequest(context.Background(), http.MethodGet, zones.ExpandPath(id), params, resBody)
		if err != nil {
			return fmt.Errorf("reading SDN zone %s: %w", id, err)
		}

		return check(resBody.Data)
	}
}

func TestAccResourceSDNZoneVXLANPeers(t *testing.T) {
	// Cannot run in parallel due to SDN applier functionality affecting global state

	te := test.InitEnvironment(t)

	config := func(peers string) string {
		return te.RenderConfig(fmt.Sprintf(`
			resource "proxmox_sdn_zone_vxlan" "zone_vxlan" {
			  id    = "zoneXP"
			  nodes = ["{{.NodeName}}"]
			  peers = %s
			}
		`, peers))
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: te.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: config(`["fd00::1", "fd00::2"]`),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("proxmox_sdn_zone_vxlan.zone_vxlan", "peers.#", "2"),
					resource.TestCheckTypeSetElemAttr("proxmox_sdn_zone_vxlan.zone_vxlan", "peers.*", "fd00::1"),
					resource.TestCheckTypeSetElemAttr("proxmox_sdn_zone_vxlan.zone_vxlan", "peers.*", "fd00::2"),
				),
			},
			{
				ResourceName:      "proxmox_sdn_zone_vxlan.zone_vxlan",
				ImportStateId:     "zoneXP",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config:      config(`[]`),
				PlanOnly:    true,
				ExpectError: regexp.MustCompile(`(?s)must\s+contain\s+at\s+least\s+1\s+elements`),
			},
			{
				Config:      config(`["10.0.0.0/24", "peer.example.com"]`),
				PlanOnly:    true,
				ExpectError: regexp.MustCompile(`(?s)must\s+be\s+a\s+valid\s+IPv4\s+or\s+IPv6\s+address`),
			},
			{
				// Post-test destroy uses the last config, which must be valid.
				Config:   config(`["fd00::1", "fd00::2"]`),
				PlanOnly: true,
			},
		},
	})
}
