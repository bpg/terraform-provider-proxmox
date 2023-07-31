/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package resource

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/bpg/terraform-provider-proxmox/proxmoxtf/test"
)

// TestGroupInstantiation tests whether the Group instance can be instantiated.
func TestGroupInstantiation(t *testing.T) {
	t.Parallel()

	s := Group()
	if s == nil {
		t.Fatalf("Cannot instantiate Group")
	}
}

// TestGroupSchema tests the Group schema.
func TestGroupSchema(t *testing.T) {
	t.Parallel()

	s := Group().Schema

	test.AssertRequiredArguments(t, s, []string{
		mkResourceVirtualEnvironmentGroupID,
	})

	test.AssertOptionalArguments(t, s, []string{
		mkResourceVirtualEnvironmentGroupComment,
	})

	test.AssertComputedAttributes(t, s, []string{
		mkResourceVirtualEnvironmentGroupMembers,
	})

	test.AssertValueTypes(t, s, map[string]schema.ValueType{
		mkResourceVirtualEnvironmentGroupComment: schema.TypeString,
		mkResourceVirtualEnvironmentGroupID:      schema.TypeString,
		mkResourceVirtualEnvironmentGroupMembers: schema.TypeSet,
	})
}
