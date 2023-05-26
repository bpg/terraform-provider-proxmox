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

	s := Pool()

	test.AssertRequiredArguments(t, s, []string{
		mkDataSourceVirtualEnvironmentPoolPoolID,
	})

	test.AssertComputedAttributes(t, s, []string{
		mkDataSourceVirtualEnvironmentPoolComment,
		mkDataSourceVirtualEnvironmentPoolMembers,
	})

	test.AssertValueTypes(t, s, map[string]schema.ValueType{
		mkDataSourceVirtualEnvironmentPoolComment: schema.TypeString,
		mkDataSourceVirtualEnvironmentPoolMembers: schema.TypeList,
		mkDataSourceVirtualEnvironmentPoolPoolID:  schema.TypeString,
	})

	membersSchema := test.AssertNestedSchemaExistence(t, s, mkDataSourceVirtualEnvironmentPoolMembers)

	test.AssertComputedAttributes(t, membersSchema, []string{
		mkDataSourceVirtualEnvironmentPoolMembersDatastoreID,
		mkDataSourceVirtualEnvironmentPoolMembersID,
		mkDataSourceVirtualEnvironmentPoolMembersNodeName,
		mkDataSourceVirtualEnvironmentPoolMembersType,
		mkDataSourceVirtualEnvironmentPoolMembersVMID,
	})

	test.AssertValueTypes(t, membersSchema, map[string]schema.ValueType{
		mkDataSourceVirtualEnvironmentPoolMembersDatastoreID: schema.TypeString,
		mkDataSourceVirtualEnvironmentPoolMembersID:          schema.TypeString,
		mkDataSourceVirtualEnvironmentPoolMembersNodeName:    schema.TypeString,
		mkDataSourceVirtualEnvironmentPoolMembersType:        schema.TypeString,
		mkDataSourceVirtualEnvironmentPoolMembersVMID:        schema.TypeInt,
	})
}
