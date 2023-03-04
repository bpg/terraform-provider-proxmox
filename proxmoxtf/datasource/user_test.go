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

// TestUserInstantiation tests whether the User instance can be instantiated.
func TestUserInstantiation(t *testing.T) {
	t.Parallel()
	s := User()

	if s == nil {
		t.Fatalf("Cannot instantiate User")
	}
}

// TestUserSchema tests the User schema.
func TestUserSchema(t *testing.T) {
	t.Parallel()
	s := User()

	test.AssertRequiredArguments(t, s, []string{
		mkDataSourceVirtualEnvironmentUserUserID,
	})

	test.AssertComputedAttributes(t, s, []string{
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

	test.AssertValueTypes(t, s, map[string]schema.ValueType{
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

	aclSchema := test.AssertNestedSchemaExistence(t, s, mkDataSourceVirtualEnvironmentGroupACL)

	test.AssertComputedAttributes(t, aclSchema, []string{
		mkDataSourceVirtualEnvironmentUserACLPath,
		mkDataSourceVirtualEnvironmentUserACLPropagate,
		mkDataSourceVirtualEnvironmentUserACLRoleID,
	})

	test.AssertValueTypes(t, aclSchema, map[string]schema.ValueType{
		mkDataSourceVirtualEnvironmentUserACLPath:      schema.TypeString,
		mkDataSourceVirtualEnvironmentUserACLPropagate: schema.TypeBool,
		mkDataSourceVirtualEnvironmentUserACLRoleID:    schema.TypeString,
	})
}
