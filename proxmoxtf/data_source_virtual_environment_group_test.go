/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package proxmoxtf

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

// TestDataSourceVirtualEnvironmentGroupInstantiation tests whether the DataSourceVirtualEnvironmentGroup instance can be instantiated.
func TestDataSourceVirtualEnvironmentGroupInstantiation(t *testing.T) {
	s := dataSourceVirtualEnvironmentGroup()

	if s == nil {
		t.Fatalf("Cannot instantiate dataSourceVirtualEnvironmentGroup")
	}
}

// TestDataSourceVirtualEnvironmentGroupSchema tests the dataSourceVirtualEnvironmentGroup schema.
func TestDataSourceVirtualEnvironmentGroupSchema(t *testing.T) {
	s := dataSourceVirtualEnvironmentGroup()

	testRequiredArguments(t, s, []string{
		mkDataSourceVirtualEnvironmentGroupID,
	})

	testComputedAttributes(t, s, []string{
		mkDataSourceVirtualEnvironmentGroupACL,
		mkDataSourceVirtualEnvironmentGroupComment,
		mkDataSourceVirtualEnvironmentGroupMembers,
	})

	testValueTypes(t, s, map[string]schema.ValueType{
		mkDataSourceVirtualEnvironmentGroupACL:     schema.TypeSet,
		mkDataSourceVirtualEnvironmentGroupID:      schema.TypeString,
		mkDataSourceVirtualEnvironmentGroupComment: schema.TypeString,
		mkDataSourceVirtualEnvironmentGroupMembers: schema.TypeSet,
	})

	aclSchema := testNestedSchemaExistence(t, s, mkDataSourceVirtualEnvironmentGroupACL)

	testComputedAttributes(t, aclSchema, []string{
		mkDataSourceVirtualEnvironmentGroupACLPath,
		mkDataSourceVirtualEnvironmentGroupACLPropagate,
		mkDataSourceVirtualEnvironmentGroupACLRoleID,
	})

	testValueTypes(t, aclSchema, map[string]schema.ValueType{
		mkDataSourceVirtualEnvironmentGroupACLPath:      schema.TypeString,
		mkDataSourceVirtualEnvironmentGroupACLPropagate: schema.TypeBool,
		mkDataSourceVirtualEnvironmentGroupACLRoleID:    schema.TypeString,
	})
}
