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

// TestGroupsInstantiation tests whether the Groups instance can be instantiated.
func TestGroupsInstantiation(t *testing.T) {
	t.Parallel()

	s := Groups()
	if s == nil {
		t.Fatalf("Cannot instantiate Groups")
	}
}

// TestGroupsSchema tests the Groups schema.
func TestGroupsSchema(t *testing.T) {
	t.Parallel()

	s := Groups()

	test.AssertComputedAttributes(t, s, []string{
		mkDataSourceVirtualEnvironmentGroupsComments,
		mkDataSourceVirtualEnvironmentGroupsGroupIDs,
	})

	test.AssertValueTypes(t, s, map[string]schema.ValueType{
		mkDataSourceVirtualEnvironmentGroupsComments: schema.TypeList,
		mkDataSourceVirtualEnvironmentGroupsGroupIDs: schema.TypeList,
	})
}
