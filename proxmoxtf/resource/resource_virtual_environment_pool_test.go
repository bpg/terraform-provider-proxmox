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

// TestResourceVirtualEnvironmentPoolInstantiation tests whether the ResourceVirtualEnvironmentPool instance can be instantiated.
func TestResourceVirtualEnvironmentPoolInstantiation(t *testing.T) {
	s := ResourceVirtualEnvironmentPool()

	if s == nil {
		t.Fatalf("Cannot instantiate ResourceVirtualEnvironmentPool")
	}
}

// TestResourceVirtualEnvironmentPoolSchema tests the ResourceVirtualEnvironmentPool schema.
func TestResourceVirtualEnvironmentPoolSchema(t *testing.T) {
	s := ResourceVirtualEnvironmentPool()

	test.AssertRequiredArguments(t, s, []string{
		mkResourceVirtualEnvironmentPoolPoolID,
	})

	test.AssertOptionalArguments(t, s, []string{
		mkResourceVirtualEnvironmentPoolComment,
	})

	test.AssertComputedAttributes(t, s, []string{
		mkResourceVirtualEnvironmentPoolMembers,
	})

	membersSchema := test.AssertNestedSchemaExistence(t, s, mkResourceVirtualEnvironmentPoolMembers)

	test.AssertComputedAttributes(t, membersSchema, []string{
		mkResourceVirtualEnvironmentPoolMembersDatastoreID,
		mkResourceVirtualEnvironmentPoolMembersID,
		mkResourceVirtualEnvironmentPoolMembersNodeName,
		mkResourceVirtualEnvironmentPoolMembersType,
		mkResourceVirtualEnvironmentPoolMembersVMID,
	})

	test.AssertValueTypes(t, membersSchema, map[string]schema.ValueType{
		mkResourceVirtualEnvironmentPoolMembersDatastoreID: schema.TypeString,
		mkResourceVirtualEnvironmentPoolMembersID:          schema.TypeString,
		mkResourceVirtualEnvironmentPoolMembersNodeName:    schema.TypeString,
		mkResourceVirtualEnvironmentPoolMembersType:        schema.TypeString,
		mkResourceVirtualEnvironmentPoolMembersVMID:        schema.TypeInt,
	})
}
