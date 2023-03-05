/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package firewall

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/bpg/terraform-provider-proxmox/proxmoxtf/test"
)

// TestIPSetsInstantiation tests whether the IPSets instance can be instantiated.
func TestIPSetsInstantiation(t *testing.T) {
	t.Parallel()
	s := IPSets()

	if s == nil {
		t.Fatalf("Cannot instantiate IPSets")
	}
}

// TestIPSetsSchema tests the IPSets schema.
func TestIPSetsSchema(t *testing.T) {
	t.Parallel()
	s := IPSets()

	test.AssertComputedAttributes(t, s, []string{
		mkIPSetsIPSetNames,
	})

	test.AssertValueTypes(t, s, map[string]schema.ValueType{
		mkIPSetsIPSetNames: schema.TypeList,
	})
}
