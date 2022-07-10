/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package proxmoxtf

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// TestResourceVirtualEnvironmentPoolInstantiation tests whether the ResourceVirtualEnvironmentPool instance can be instantiated.
func TestResourceVirtualEnvironmentPoolInstantiation(t *testing.T) {
	s := resourceVirtualEnvironmentPool()

	if s == nil {
		t.Fatalf("Cannot instantiate resourceVirtualEnvironmentPool")
	}
}

// TestResourceVirtualEnvironmentPoolSchema tests the resourceVirtualEnvironmentPool schema.
func TestResourceVirtualEnvironmentPoolSchema(t *testing.T) {
	s := resourceVirtualEnvironmentPool()

	testRequiredArguments(t, s, []string{
		mkResourceVirtualEnvironmentPoolPoolID,
	})

	testOptionalArguments(t, s, []string{
		mkResourceVirtualEnvironmentPoolComment,
	})

	testComputedAttributes(t, s, []string{
		mkResourceVirtualEnvironmentPoolMembers,
	})

	membersSchema := testNestedSchemaExistence(t, s, mkResourceVirtualEnvironmentPoolMembers)

	testComputedAttributes(t, membersSchema, []string{
		mkResourceVirtualEnvironmentPoolMembersDatastoreID,
		mkResourceVirtualEnvironmentPoolMembersID,
		mkResourceVirtualEnvironmentPoolMembersNodeName,
		mkResourceVirtualEnvironmentPoolMembersType,
		mkResourceVirtualEnvironmentPoolMembersVMID,
	})

	testValueTypes(t, membersSchema, map[string]schema.ValueType{
		mkResourceVirtualEnvironmentPoolMembersDatastoreID: schema.TypeString,
		mkResourceVirtualEnvironmentPoolMembersID:          schema.TypeString,
		mkResourceVirtualEnvironmentPoolMembersNodeName:    schema.TypeString,
		mkResourceVirtualEnvironmentPoolMembersType:        schema.TypeString,
		mkResourceVirtualEnvironmentPoolMembersVMID:        schema.TypeInt,
	})
}
