/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package proxmoxtf

import (
	"testing"

	"github.com/hashicorp/terraform/helper/schema"
)

// TestDataSourceVirtualEnvironmentUserInstantiation tests whether the DataSourceVirtualEnvironmentUser instance can be instantiated.
func TestDataSourceVirtualEnvironmentUserInstantiation(t *testing.T) {
	s := dataSourceVirtualEnvironmentUser()

	if s == nil {
		t.Fatalf("Cannot instantiate dataSourceVirtualEnvironmentUser")
	}
}

// TestDataSourceVirtualEnvironmentUserSchema tests the dataSourceVirtualEnvironmentUser schema.
func TestDataSourceVirtualEnvironmentUserSchema(t *testing.T) {
	s := dataSourceVirtualEnvironmentUser()

	testRequiredArguments(t, s, []string{
		mkDataSourceVirtualEnvironmentUserUserID,
	})

	testComputedAttributes(t, s, []string{
		mkDataSourceVirtualEnvironmentUserACL,
		mkDataSourceVirtualEnvironmentUserComment,
		mkDataSourceVirtualEnvironmentUserEmail,
		mkDataSourceVirtualEnvironmentUserEnabled,
		mkDataSourceVirtualEnvironmentUserExpirationDate,
		mkDataSourceVirtualEnvironmentUserFirstName,
		mkDataSourceVirtualEnvironmentUserGroups,
		mkDataSourceVirtualEnvironmentUserKeys,
		mkDataSourceVirtualEnvironmentUserLastName,
	})

	testSchemaValueTypes(t, s, []string{
		mkDataSourceVirtualEnvironmentUserACL,
		mkDataSourceVirtualEnvironmentUserComment,
		mkDataSourceVirtualEnvironmentUserEmail,
		mkDataSourceVirtualEnvironmentUserEnabled,
		mkDataSourceVirtualEnvironmentUserExpirationDate,
		mkDataSourceVirtualEnvironmentUserFirstName,
		mkDataSourceVirtualEnvironmentUserGroups,
		mkDataSourceVirtualEnvironmentUserKeys,
		mkDataSourceVirtualEnvironmentUserLastName,
	}, []schema.ValueType{
		schema.TypeSet,
		schema.TypeString,
		schema.TypeString,
		schema.TypeBool,
		schema.TypeString,
		schema.TypeString,
		schema.TypeList,
		schema.TypeString,
		schema.TypeString,
	})

	aclSchema := testNestedSchemaExistence(t, s, mkDataSourceVirtualEnvironmentGroupACL)

	testComputedAttributes(t, aclSchema, []string{
		mkDataSourceVirtualEnvironmentUserACLPath,
		mkDataSourceVirtualEnvironmentUserACLPropagate,
		mkDataSourceVirtualEnvironmentUserACLRoleID,
	})

	testSchemaValueTypes(t, aclSchema, []string{
		mkDataSourceVirtualEnvironmentUserACLPath,
		mkDataSourceVirtualEnvironmentUserACLPropagate,
		mkDataSourceVirtualEnvironmentUserACLRoleID,
	}, []schema.ValueType{
		schema.TypeString,
		schema.TypeBool,
		schema.TypeString,
	})
}
