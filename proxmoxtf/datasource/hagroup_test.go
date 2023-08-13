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

// TestHAGroupInstantiation tests whether HAGroup can be instantiated.
func TestHAGroupInstantiation(t *testing.T) {
	t.Parallel()

	s := HAGroup()
	if s == nil {
		t.Fatalf("Cannot instantiate HAGroup")
	}
}

// TestHAGroupSchema tests the HAGroup schema.
func TestHAGroupSchema(t *testing.T) {
	t.Parallel()

	s := HAGroup()

	test.AssertRequiredArguments(t, s, []string{
		mkDataSourceVirtualEnvironmentHAGroupID,
	})

	test.AssertComputedAttributes(t, s, []string{
		mkDataSourceVirtualEnvironmentHAGroupComment,
		mkDataSourceVirtualEnvironmentHAGroupMembers,
		mkDataSourceVirtualEnvironmentHAGroupNoFailback,
		mkDataSourceVirtualEnvironmentHAGroupRestricted,
	})

	test.AssertValueTypes(t, s, map[string]schema.ValueType{
		mkDataSourceVirtualEnvironmentHAGroupID:         schema.TypeString,
		mkDataSourceVirtualEnvironmentHAGroupComment:    schema.TypeString,
		mkDataSourceVirtualEnvironmentHAGroupMembers:    schema.TypeList,
		mkDataSourceVirtualEnvironmentHAGroupNoFailback: schema.TypeBool,
		mkDataSourceVirtualEnvironmentHAGroupRestricted: schema.TypeBool,
	})

	membersSchema := test.AssertNestedSchemaExistence(t, s, mkDataSourceVirtualEnvironmentHAGroupMembers)

	test.AssertComputedAttributes(t, membersSchema, []string{
		mkDataSourceVirtualEnvironmentHAGroupMemberNodeName,
		mkDataSourceVirtualEnvironmentHAGroupMemberPriority,
	})

	test.AssertValueTypes(t, membersSchema, map[string]schema.ValueType{
		mkDataSourceVirtualEnvironmentHAGroupMemberNodeName: schema.TypeString,
		mkDataSourceVirtualEnvironmentHAGroupMemberPriority: schema.TypeInt,
	})
}
