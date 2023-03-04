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

// TestResourceVirtualEnvironmentUserInstantiation tests whether the ResourceVirtualEnvironmentUser instance can be instantiated.
func TestResourceVirtualEnvironmentUserInstantiation(t *testing.T) {
	s := ResourceVirtualEnvironmentUser()

	if s == nil {
		t.Fatalf("Cannot instantiate ResourceVirtualEnvironmentUser")
	}
}

// TestResourceVirtualEnvironmentUserSchema tests the ResourceVirtualEnvironmentUser schema.
func TestResourceVirtualEnvironmentUserSchema(t *testing.T) {
	s := ResourceVirtualEnvironmentUser()

	test.AssertRequiredArguments(t, s, []string{
		mkResourceVirtualEnvironmentUserPassword,
		mkResourceVirtualEnvironmentUserUserID,
	})

	test.AssertOptionalArguments(t, s, []string{
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

	test.AssertValueTypes(t, s, map[string]schema.ValueType{
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

	aclSchema := test.AssertNestedSchemaExistence(t, s, mkResourceVirtualEnvironmentUserACL)

	test.AssertRequiredArguments(t, aclSchema, []string{
		mkResourceVirtualEnvironmentUserACLPath,
		mkResourceVirtualEnvironmentUserACLRoleID,
	})

	test.AssertOptionalArguments(t, aclSchema, []string{
		mkResourceVirtualEnvironmentUserACLPropagate,
	})

	test.AssertValueTypes(t, aclSchema, map[string]schema.ValueType{
		mkResourceVirtualEnvironmentUserACLPath:      schema.TypeString,
		mkResourceVirtualEnvironmentUserACLPropagate: schema.TypeBool,
		mkResourceVirtualEnvironmentUserACLRoleID:    schema.TypeString,
	})
}
