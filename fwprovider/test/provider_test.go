//go:build acceptance || all

//testacc:tier=light
//testacc:resource=misc

/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package test

import (
	"fmt"
	"net/url"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/stretchr/testify/require"

	"github.com/bpg/terraform-provider-proxmox/utils"
)

func TestAccProviderSSHNodeAddressSource(t *testing.T) {
	te := InitEnvironment(t)

	nodeName := utils.GetAnyStringEnv("PROXMOX_VE_ACC_NODE_NAME")
	if nodeName == "" {
		nodeName = "pve"
	}

	nodeAddress := utils.GetAnyStringEnv("PROXMOX_VE_ACC_NODE_SSH_ADDRESS")
	if nodeAddress == "" {
		endpoint := utils.GetAnyStringEnv("PROXMOX_VE_ENDPOINT")

		u, err := url.Parse(endpoint)
		require.NoError(t, err)

		nodeAddress = u.Hostname()
	}

	nodePort := utils.GetAnyStringEnv("PROXMOX_VE_ACC_NODE_SSH_PORT")
	if nodePort == "" {
		nodePort = "22"
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: te.AccProviders,
		Steps: []resource.TestStep{{
			Config: fmt.Sprintf(`
				provider "proxmox" {
					ssh {
						node_address_source = "dns"
						node {
							name    = %q
							address = %q
							port    = %s
						}
					}
				}

				data "proxmox_virtual_environment_version" "test" {}
			`, nodeName, nodeAddress, nodePort),
			Check: resource.TestCheckResourceAttrSet("data.proxmox_virtual_environment_version.test", "version"),
		}},
	})
}
