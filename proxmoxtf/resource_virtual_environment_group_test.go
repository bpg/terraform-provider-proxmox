/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package proxmoxtf

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

// TestResourceVirtualEnvironmentGroupInstantiation tests whether the ResourceVirtualEnvironmentGroup instance can be instantiated.
func TestResourceVirtualEnvironmentGroupInstantiation(t *testing.T) {
	s := resourceVirtualEnvironmentGroup()

	if s == nil {
		t.Fatalf("Cannot instantiate resourceVirtualEnvironmentGroup")
	}
}

// TestResourceVirtualEnvironmentGroupSchema tests the resourceVirtualEnvironmentGroup schema.
func TestResourceVirtualEnvironmentGroupSchema(t *testing.T) {
	s := resourceVirtualEnvironmentGroup()

	testRequiredArguments(t, s, []string{
		mkResourceVirtualEnvironmentGroupID,
	})

	testOptionalArguments(t, s, []string{
		mkResourceVirtualEnvironmentGroupACL,
		mkResourceVirtualEnvironmentGroupComment,
	})

	testComputedAttributes(t, s, []string{
		mkResourceVirtualEnvironmentGroupMembers,
	})

	testValueTypes(t, s, map[string]schema.ValueType{
		mkResourceVirtualEnvironmentGroupACL:     schema.TypeSet,
		mkResourceVirtualEnvironmentGroupComment: schema.TypeString,
		mkResourceVirtualEnvironmentGroupID:      schema.TypeString,
		mkResourceVirtualEnvironmentGroupMembers: schema.TypeSet,
	})

	aclSchema := testNestedSchemaExistence(t, s, mkResourceVirtualEnvironmentGroupACL)

	testRequiredArguments(t, aclSchema, []string{
		mkResourceVirtualEnvironmentGroupACLPath,
		mkResourceVirtualEnvironmentGroupACLRoleID,
	})

	testOptionalArguments(t, aclSchema, []string{
		mkResourceVirtualEnvironmentGroupACLPropagate,
	})

	testValueTypes(t, aclSchema, map[string]schema.ValueType{
		mkResourceVirtualEnvironmentGroupACLPath:      schema.TypeString,
		mkResourceVirtualEnvironmentGroupACLPropagate: schema.TypeBool,
		mkResourceVirtualEnvironmentGroupACLRoleID:    schema.TypeString,
	})
}
