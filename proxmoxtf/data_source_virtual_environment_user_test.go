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

	testValueTypes(t, s, map[string]schema.ValueType{
		mkDataSourceVirtualEnvironmentUserACL:            schema.TypeSet,
		mkDataSourceVirtualEnvironmentUserComment:        schema.TypeString,
		mkDataSourceVirtualEnvironmentUserEmail:          schema.TypeString,
		mkDataSourceVirtualEnvironmentUserEnabled:        schema.TypeBool,
		mkDataSourceVirtualEnvironmentUserExpirationDate: schema.TypeString,
		mkDataSourceVirtualEnvironmentUserFirstName:      schema.TypeString,
		mkDataSourceVirtualEnvironmentUserGroups:         schema.TypeList,
		mkDataSourceVirtualEnvironmentUserKeys:           schema.TypeString,
		mkDataSourceVirtualEnvironmentUserLastName:       schema.TypeString,
	})

	aclSchema := testNestedSchemaExistence(t, s, mkDataSourceVirtualEnvironmentGroupACL)

	testComputedAttributes(t, aclSchema, []string{
		mkDataSourceVirtualEnvironmentUserACLPath,
		mkDataSourceVirtualEnvironmentUserACLPropagate,
		mkDataSourceVirtualEnvironmentUserACLRoleID,
	})

	testValueTypes(t, aclSchema, map[string]schema.ValueType{
		mkDataSourceVirtualEnvironmentUserACLPath:      schema.TypeString,
		mkDataSourceVirtualEnvironmentUserACLPropagate: schema.TypeBool,
		mkDataSourceVirtualEnvironmentUserACLRoleID:    schema.TypeString,
	})
}
