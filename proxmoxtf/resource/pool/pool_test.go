/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package pool

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/bpg/terraform-provider-proxmox/proxmoxtf/test"
)

// TestPoolInstantiation tests whether the Pool instance can be instantiated.
func TestPoolInstantiation(t *testing.T) {
	t.Parallel()

	s := Pool()
	if s == nil {
		t.Fatalf("Cannot instantiate Pool")
	}
}

// TestPoolSchema tests the Pool schema.
func TestPoolSchema(t *testing.T) {
	t.Parallel()

	s := Pool().Schema

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
