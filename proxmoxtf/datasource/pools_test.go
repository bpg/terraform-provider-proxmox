/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package datasource

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/bpg/terraform-provider-proxmox/proxmoxtf/test"
)

// TestPoolsInstantiation tests whether the Pools instance can be instantiated.
func TestPoolsInstantiation(t *testing.T) {
	t.Parallel()

	s := Pools()
	if s == nil {
		t.Fatalf("Cannot instantiate Pools")
	}
}

// TestPoolsSchema tests the Pools schema.
func TestPoolsSchema(t *testing.T) {
	t.Parallel()

	s := Pools()

	test.AssertComputedAttributes(t, s, []string{
		mkDataSourceVirtualEnvironmentPoolsPoolIDs,
	})

	test.AssertValueTypes(t, s, map[string]schema.ValueType{
		mkDataSourceVirtualEnvironmentPoolsPoolIDs: schema.TypeList,
	})
}
