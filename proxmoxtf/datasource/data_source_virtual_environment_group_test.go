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

// TestDataSourceVirtualEnvironmentGroupInstantiation tests whether the DataSourceVirtualEnvironmentGroup instance can be instantiated.
func TestDataSourceVirtualEnvironmentGroupInstantiation(t *testing.T) {
	s := DataSourceVirtualEnvironmentGroup()

	if s == nil {
		t.Fatalf("Cannot instantiate DataSourceVirtualEnvironmentGroup")
	}
}

// TestDataSourceVirtualEnvironmentGroupSchema tests the DataSourceVirtualEnvironmentGroup schema.
func TestDataSourceVirtualEnvironmentGroupSchema(t *testing.T) {
	s := DataSourceVirtualEnvironmentGroup()

	test.AssertRequiredArguments(t, s, []string{
		mkDataSourceVirtualEnvironmentGroupID,
	})

	test.AssertComputedAttributes(t, s, []string{
		mkDataSourceVirtualEnvironmentGroupACL,
		mkDataSourceVirtualEnvironmentGroupComment,
		mkDataSourceVirtualEnvironmentGroupMembers,
	})

	test.AssertValueTypes(t, s, map[string]schema.ValueType{
		mkDataSourceVirtualEnvironmentGroupACL:     schema.TypeSet,
		mkDataSourceVirtualEnvironmentGroupID:      schema.TypeString,
		mkDataSourceVirtualEnvironmentGroupComment: schema.TypeString,
		mkDataSourceVirtualEnvironmentGroupMembers: schema.TypeSet,
	})

	aclSchema := test.AssertNestedSchemaExistence(t, s, mkDataSourceVirtualEnvironmentGroupACL)

	test.AssertComputedAttributes(t, aclSchema, []string{
		mkDataSourceVirtualEnvironmentGroupACLPath,
		mkDataSourceVirtualEnvironmentGroupACLPropagate,
		mkDataSourceVirtualEnvironmentGroupACLRoleID,
	})

	test.AssertValueTypes(t, aclSchema, map[string]schema.ValueType{
		mkDataSourceVirtualEnvironmentGroupACLPath:      schema.TypeString,
		mkDataSourceVirtualEnvironmentGroupACLPropagate: schema.TypeBool,
		mkDataSourceVirtualEnvironmentGroupACLRoleID:    schema.TypeString,
	})
}
