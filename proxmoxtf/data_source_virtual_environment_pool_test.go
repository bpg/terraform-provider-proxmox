/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package proxmoxtf

import (
	"testing"

	"github.com/hashicorp/terraform/helper/schema"
)

// TestDataSourceVirtualEnvironmentPoolInstantiation tests whether the DataSourceVirtualEnvironmentPool instance can be instantiated.
func TestDataSourceVirtualEnvironmentPoolInstantiation(t *testing.T) {
	s := dataSourceVirtualEnvironmentPool()

	if s == nil {
		t.Fatalf("Cannot instantiate dataSourceVirtualEnvironmentPool")
	}
}

// TestDataSourceVirtualEnvironmentPoolSchema tests the dataSourceVirtualEnvironmentPool schema.
func TestDataSourceVirtualEnvironmentPoolSchema(t *testing.T) {
	s := dataSourceVirtualEnvironmentPool()

	testRequiredArguments(t, s, []string{
		mkDataSourceVirtualEnvironmentPoolPoolID,
	})

	testComputedAttributes(t, s, []string{
		mkDataSourceVirtualEnvironmentPoolComment,
		mkDataSourceVirtualEnvironmentPoolMembers,
	})

	testSchemaValueTypes(t, s, map[string]schema.ValueType{
		mkDataSourceVirtualEnvironmentPoolComment: schema.TypeString,
		mkDataSourceVirtualEnvironmentPoolMembers: schema.TypeList,
		mkDataSourceVirtualEnvironmentPoolPoolID:  schema.TypeString,
	})

	membersSchema := testNestedSchemaExistence(t, s, mkDataSourceVirtualEnvironmentPoolMembers)

	testComputedAttributes(t, membersSchema, []string{
		mkDataSourceVirtualEnvironmentPoolMembersDatastoreID,
		mkDataSourceVirtualEnvironmentPoolMembersID,
		mkDataSourceVirtualEnvironmentPoolMembersNodeName,
		mkDataSourceVirtualEnvironmentPoolMembersType,
		mkDataSourceVirtualEnvironmentPoolMembersVMID,
	})

	testSchemaValueTypes(t, membersSchema, map[string]schema.ValueType{
		mkDataSourceVirtualEnvironmentPoolMembersDatastoreID: schema.TypeString,
		mkDataSourceVirtualEnvironmentPoolMembersID:          schema.TypeString,
		mkDataSourceVirtualEnvironmentPoolMembersNodeName:    schema.TypeString,
		mkDataSourceVirtualEnvironmentPoolMembersType:        schema.TypeString,
		mkDataSourceVirtualEnvironmentPoolMembersVMID:        schema.TypeInt,
	})
}
