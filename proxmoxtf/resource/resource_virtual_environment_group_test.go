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

// TestResourceVirtualEnvironmentGroupInstantiation tests whether the ResourceVirtualEnvironmentGroup instance can be instantiated.
func TestResourceVirtualEnvironmentGroupInstantiation(t *testing.T) {
	s := ResourceVirtualEnvironmentGroup()

	if s == nil {
		t.Fatalf("Cannot instantiate ResourceVirtualEnvironmentGroup")
	}
}

// TestResourceVirtualEnvironmentGroupSchema tests the ResourceVirtualEnvironmentGroup schema.
func TestResourceVirtualEnvironmentGroupSchema(t *testing.T) {
	s := ResourceVirtualEnvironmentGroup()

	test.AssertRequiredArguments(t, s, []string{
		mkResourceVirtualEnvironmentGroupID,
	})

	test.AssertOptionalArguments(t, s, []string{
		mkResourceVirtualEnvironmentGroupACL,
		mkResourceVirtualEnvironmentGroupComment,
	})

	test.AssertComputedAttributes(t, s, []string{
		mkResourceVirtualEnvironmentGroupMembers,
	})

	test.AssertValueTypes(t, s, map[string]schema.ValueType{
		mkResourceVirtualEnvironmentGroupACL:     schema.TypeSet,
		mkResourceVirtualEnvironmentGroupComment: schema.TypeString,
		mkResourceVirtualEnvironmentGroupID:      schema.TypeString,
		mkResourceVirtualEnvironmentGroupMembers: schema.TypeSet,
	})

	aclSchema := test.AssertNestedSchemaExistence(t, s, mkResourceVirtualEnvironmentGroupACL)

	test.AssertRequiredArguments(t, aclSchema, []string{
		mkResourceVirtualEnvironmentGroupACLPath,
		mkResourceVirtualEnvironmentGroupACLRoleID,
	})

	test.AssertOptionalArguments(t, aclSchema, []string{
		mkResourceVirtualEnvironmentGroupACLPropagate,
	})

	test.AssertValueTypes(t, aclSchema, map[string]schema.ValueType{
		mkResourceVirtualEnvironmentGroupACLPath:      schema.TypeString,
		mkResourceVirtualEnvironmentGroupACLPropagate: schema.TypeBool,
		mkResourceVirtualEnvironmentGroupACLRoleID:    schema.TypeString,
	})
}
