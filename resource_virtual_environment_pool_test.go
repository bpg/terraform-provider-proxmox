/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package main

import (
	"testing"

	"github.com/hashicorp/terraform/helper/schema"
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

	testSchemaValueTypes(t, membersSchema, []string{
		mkResourceVirtualEnvironmentPoolMembersDatastoreID,
		mkResourceVirtualEnvironmentPoolMembersID,
		mkResourceVirtualEnvironmentPoolMembersNodeName,
		mkResourceVirtualEnvironmentPoolMembersType,
		mkResourceVirtualEnvironmentPoolMembersVMID,
	}, []schema.ValueType{
		schema.TypeString,
		schema.TypeString,
		schema.TypeString,
		schema.TypeString,
		schema.TypeInt,
	})
}
