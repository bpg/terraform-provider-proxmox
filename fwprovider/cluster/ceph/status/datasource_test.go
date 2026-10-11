//go:build acceptance || all

//testacc:tier=light
//testacc:resource=ceph

/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package status_test

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/bpg/terraform-provider-proxmox/fwprovider/test"
)

func TestAccDataSourceCephStatus(t *testing.T) {
	t.Parallel()

	te := test.InitEnvironment(t)
	te.RequireCeph()

	healthRe := regexp.MustCompile(`^HEALTH_(OK|WARN|ERR)$`)
	fsidRe := regexp.MustCompile(`^[0-9a-f-]{36}$`)

	tests := []struct {
		name string
		step []resource.TestStep
	}{
		{"cluster scope (no node_name)", []resource.TestStep{
			{
				Config: `data "proxmox_ceph_status" "cluster" {}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestMatchResourceAttr("data.proxmox_ceph_status.cluster", "fsid", fsidRe),
					resource.TestCheckResourceAttrPair(
						"data.proxmox_ceph_status.cluster", "id",
						"data.proxmox_ceph_status.cluster", "fsid",
					),
					resource.TestMatchResourceAttr("data.proxmox_ceph_status.cluster", "health_status", healthRe),
					resource.TestCheckResourceAttrSet("data.proxmox_ceph_status.cluster", "quorum_names.0"),
				),
			},
		}},
		{"node scope (node_name set)", []resource.TestStep{
			{
				Config: te.RenderConfig(`
					data "proxmox_ceph_status" "node" {
						node_name = "{{.NodeName}}"
					}`),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestMatchResourceAttr("data.proxmox_ceph_status.node", "fsid", fsidRe),
					resource.TestCheckResourceAttrPair(
						"data.proxmox_ceph_status.node", "id",
						"data.proxmox_ceph_status.node", "fsid",
					),
					resource.TestMatchResourceAttr("data.proxmox_ceph_status.node", "health_status", healthRe),
					resource.TestCheckResourceAttrSet("data.proxmox_ceph_status.node", "quorum_names.0"),
					resource.TestCheckResourceAttr("data.proxmox_ceph_status.node", "node_name", te.NodeName),
				),
			},
		}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resource.Test(t, resource.TestCase{
				ProtoV6ProviderFactories: te.AccProviders,
				Steps:                    tt.step,
			})
		})
	}
}
