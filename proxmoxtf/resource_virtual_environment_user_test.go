/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package proxmoxtf

import (
	"testing"

	"github.com/hashicorp/terraform/helper/schema"
)

// TestResourceVirtualEnvironmentUserInstantiation tests whether the ResourceVirtualEnvironmentUser instance can be instantiated.
func TestResourceVirtualEnvironmentUserInstantiation(t *testing.T) {
	s := resourceVirtualEnvironmentUser()

	if s == nil {
		t.Fatalf("Cannot instantiate resourceVirtualEnvironmentUser")
	}
}

// TestResourceVirtualEnvironmentUserSchema tests the resourceVirtualEnvironmentUser schema.
func TestResourceVirtualEnvironmentUserSchema(t *testing.T) {
	s := resourceVirtualEnvironmentUser()

	testRequiredArguments(t, s, []string{
		mkResourceVirtualEnvironmentUserPassword,
		mkResourceVirtualEnvironmentUserUserID,
	})

	testOptionalArguments(t, s, []string{
		mkResourceVirtualEnvironmentUserACL,
		mkResourceVirtualEnvironmentUserComment,
		mkResourceVirtualEnvironmentUserEmail,
		mkResourceVirtualEnvironmentUserEnabled,
		mkResourceVirtualEnvironmentUserExpirationDate,
		mkResourceVirtualEnvironmentUserFirstName,
		mkResourceVirtualEnvironmentUserGroups,
		mkResourceVirtualEnvironmentUserKeys,
		mkResourceVirtualEnvironmentUserLastName,
	})

	testSchemaValueTypes(t, s, map[string]schema.ValueType{
		mkResourceVirtualEnvironmentUserACL:            schema.TypeSet,
		mkResourceVirtualEnvironmentUserComment:        schema.TypeString,
		mkResourceVirtualEnvironmentUserEmail:          schema.TypeString,
		mkResourceVirtualEnvironmentUserEnabled:        schema.TypeBool,
		mkResourceVirtualEnvironmentUserExpirationDate: schema.TypeString,
		mkResourceVirtualEnvironmentUserFirstName:      schema.TypeString,
		mkResourceVirtualEnvironmentUserGroups:         schema.TypeSet,
		mkResourceVirtualEnvironmentUserKeys:           schema.TypeString,
		mkResourceVirtualEnvironmentUserLastName:       schema.TypeString,
		mkResourceVirtualEnvironmentUserPassword:       schema.TypeString,
		mkResourceVirtualEnvironmentUserUserID:         schema.TypeString,
	})

	aclSchema := testNestedSchemaExistence(t, s, mkResourceVirtualEnvironmentUserACL)

	testRequiredArguments(t, aclSchema, []string{
		mkResourceVirtualEnvironmentUserACLPath,
		mkResourceVirtualEnvironmentUserACLRoleID,
	})

	testOptionalArguments(t, aclSchema, []string{
		mkResourceVirtualEnvironmentUserACLPropagate,
	})

	testSchemaValueTypes(t, aclSchema, map[string]schema.ValueType{
		mkResourceVirtualEnvironmentUserACLPath:      schema.TypeString,
		mkResourceVirtualEnvironmentUserACLPropagate: schema.TypeBool,
		mkResourceVirtualEnvironmentUserACLRoleID:    schema.TypeString,
	})
}
