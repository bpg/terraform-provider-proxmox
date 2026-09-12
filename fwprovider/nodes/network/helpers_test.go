//go:build acceptance || all

//testacc:tier=heavy
//testacc:resource=network

/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package network_test

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"

	"github.com/bpg/terraform-provider-proxmox/fwprovider/test"
)

// checkInterfaceActive asserts whether PVE reports an interface as active on the node.
//
// A staged change is listed but not up in the kernel, so `active` stays false until a reload.
// Terraform state alone would pass whether or not `reload` was honoured.
func checkInterfaceActive(te *test.Environment, iface string, want bool) resource.TestCheckFunc {
	return func(_ *terraform.State) error {
		ifaces, err := te.NodeClient().ListNetworkInterfaces(context.Background())
		if err != nil {
			return fmt.Errorf("listing network interfaces: %w", err)
		}

		for _, i := range ifaces {
			if i.Iface != iface {
				continue
			}

			got := i.Active != nil && bool(*i.Active)
			if got != want {
				return fmt.Errorf("interface %q: active = %t, want %t", iface, got, want)
			}

			return nil
		}

		return fmt.Errorf("interface %q not found on node", iface)
	}
}

// reloadNode applies whatever is staged, so a test that leaves a pending change on purpose
// does not hand the next test a dirty node.
func reloadNode(te *test.Environment) resource.TestCheckFunc {
	return func(_ *terraform.State) error {
		return te.NodeClient().ReloadNetworkConfiguration(context.Background())
	}
}

// checkNodePendingChanges asserts whether the node has staged network changes.
//
// A staged delete drops out of the interface list immediately, so `active` cannot tell it apart
// from an applied one. PVE only sends `changes` while something is pending.
func checkNodePendingChanges(te *test.Environment, want bool) resource.TestCheckFunc {
	return func(_ *terraform.State) error {
		var body struct {
			Changes *string `json:"changes"`
		}

		path := fmt.Sprintf("nodes/%s/network", te.NodeName)
		if err := te.Client().DoRequest(context.Background(), http.MethodGet, path, nil, &body); err != nil {
			return fmt.Errorf("reading node network: %w", err)
		}

		got := body.Changes != nil && *body.Changes != ""
		if got != want {
			return fmt.Errorf("node %q: pending changes = %t, want %t", te.NodeName, got, want)
		}

		return nil
	}
}

// checkStagedDestroy asserts that a destroy with `reload = false` left its delete staged rather
// than applied, then clears the pending change off the shared node.
func checkStagedDestroy(te *test.Environment) resource.TestCheckFunc {
	return resource.ComposeTestCheckFunc(
		checkNodePendingChanges(te, true),
		reloadNode(te),
	)
}
