/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package proxmoxtf

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// TestDataSourceVirtualEnvironmentVMsInstantiation tests whether the dataSourceVirtualEnvironmentVMs instance can be instantiated.
func TestDataSourceVirtualEnvironmentVMsInstantiation(t *testing.T) {
	t.Parallel()

	s := dataSourceVirtualEnvironmentVMs()

	if s == nil {
		t.Fatalf("Cannot instantiate dataSourceVirtualEnvironmentVMs")
	}
}

// TestDataSourceVirtualEnvironmentVMsSchema tests the dataSourceVirtualEnvironmentVMs schema.
func TestDataSourceVirtualEnvironmentVMsSchema(t *testing.T) {
	t.Parallel()

	s := dataSourceVirtualEnvironmentVMs()

	testComputedAttributes(t, s, []string{
		mkDataSourceVirtualEnvironmentVMs,
	})

	testValueTypes(t, s, map[string]schema.ValueType{
		mkDataSourceVirtualEnvironmentVMNodeName: schema.TypeString,
		mkDataSourceVirtualEnvironmentVMTags:     schema.TypeList,
		mkDataSourceVirtualEnvironmentVMs:        schema.TypeList,
	})

	vmsSchema := testNestedSchemaExistence(t, s, mkDataSourceVirtualEnvironmentVMs)

	testComputedAttributes(t, vmsSchema, []string{
		mkDataSourceVirtualEnvironmentVMName,
		mkDataSourceVirtualEnvironmentVMTags,
	})

	testValueTypes(t, vmsSchema, map[string]schema.ValueType{
		mkDataSourceVirtualEnvironmentVMName:     schema.TypeString,
		mkDataSourceVirtualEnvironmentVMNodeName: schema.TypeString,
		mkDataSourceVirtualEnvironmentVMTags:     schema.TypeList,
		mkDataSourceVirtualEnvironmentVMVMID:     schema.TypeInt,
	})
}
