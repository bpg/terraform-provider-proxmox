/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package main

import (
	"testing"

	"github.com/hashicorp/terraform/helper/schema"
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

	testSchemaValueTypes(t, s, []string{
		mkDataSourceVirtualEnvironmentGroupACL,
		mkDataSourceVirtualEnvironmentGroupID,
		mkDataSourceVirtualEnvironmentGroupComment,
		mkDataSourceVirtualEnvironmentGroupMembers,
	}, []schema.ValueType{
		schema.TypeSet,
		schema.TypeString,
		schema.TypeString,
		schema.TypeSet,
	})

	aclSchema := testNestedSchemaExistence(t, s, mkDataSourceVirtualEnvironmentGroupACL)

	testComputedAttributes(t, aclSchema, []string{
		mkDataSourceVirtualEnvironmentGroupACLPath,
		mkDataSourceVirtualEnvironmentGroupACLPropagate,
		mkDataSourceVirtualEnvironmentGroupACLRoleID,
	})

	testSchemaValueTypes(t, aclSchema, []string{
		mkDataSourceVirtualEnvironmentGroupACLPath,
		mkDataSourceVirtualEnvironmentGroupACLPropagate,
		mkDataSourceVirtualEnvironmentGroupACLRoleID,
	}, []schema.ValueType{
		schema.TypeString,
		schema.TypeBool,
		schema.TypeString,
	})
}
