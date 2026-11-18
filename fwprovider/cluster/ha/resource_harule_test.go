//go:build acceptance || all

//testacc:tier=medium
//testacc:resource=ha

/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package ha_test

import (
	"fmt"
	"math/rand"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/stretchr/testify/require"

	"github.com/bpg/terraform-provider-proxmox/fwprovider/test"
)

func TestAccResourceHARuleBoolToggle(t *testing.T) {
	te := test.InitEnvironment(t)

	skipIfNoHA(t, te)

	vmID := 100000 + rand.Intn(99999)
	ruleID := fmt.Sprintf("test-bool-toggle-%d", vmID)

	te.AddTemplateVars(map[string]any{
		"TestVMID": vmID,
		"RuleID":   ruleID,
	})

	config := func(bools string) string {
		return te.RenderConfig(`
			resource "proxmox_virtual_environment_vm" "test_vm" {
				node_name = "{{.NodeName}}"
				vm_id     = {{.TestVMID}}
				started   = false
				name      = "test-harule-bool-toggle"
			}

			resource "proxmox_haresource" "test_res" {
				depends_on  = [proxmox_virtual_environment_vm.test_vm]
				resource_id = "vm:{{.TestVMID}}"
				state       = "stopped"
			}

			resource "proxmox_harule" "test_rule" {
				depends_on = [proxmox_haresource.test_res]
				rule       = "{{.RuleID}}"
				type       = "node-affinity"
				resources  = ["vm:{{.TestVMID}}"]
				nodes      = { "{{.NodeName}}" = 1 }
				` + bools + `
			}
		`)
	}

	check := func(disable, strict bool) resource.TestCheckFunc {
		return resource.ComposeTestCheckFunc(
			resource.TestCheckResourceAttr("proxmox_harule.test_rule", "disable", fmt.Sprint(disable)),
			resource.TestCheckResourceAttr("proxmox_harule.test_rule", "strict", fmt.Sprint(strict)),
			func(*terraform.State) error {
				rule, err := te.ClusterClient().HA().Rules().Get(t.Context(), ruleID)
				require.NoError(t, err)
				require.Equal(t, disable, rule.Disable != nil && bool(*rule.Disable), "disable in PVE")
				require.Equal(t, strict, rule.Strict != nil && bool(*rule.Strict), "strict in PVE")

				return nil
			},
		)
	}

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: te.AccProviders,
		Steps: []resource.TestStep{
			{Config: config(""), Check: check(false, false)},
			{Config: config("disable = true\nstrict = true"), Check: check(true, true)},
			{Config: config(""), Check: check(false, false)},
			{Config: config("disable = true"), Check: check(true, false)},
			{
				ResourceName:      "proxmox_harule.test_rule",
				ImportState:       true,
				ImportStateId:     ruleID,
				ImportStateVerify: true,
			},
			{Config: config("disable = false"), Check: check(false, false)},
		},
	})
}
