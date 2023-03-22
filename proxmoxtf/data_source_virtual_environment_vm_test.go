/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package proxmoxtf

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// TestDataSourceVirtualEnvironmentVMInstantiation tests whether the dataSourceVirtualEnvironmentVM instance can be instantiated.
func TestDataSourceVirtualEnvironmentVMInstantiation(t *testing.T) {
	t.Parallel()

	s := dataSourceVirtualEnvironmentVM()

	if s == nil {
		t.Fatalf("Cannot instantiate dataSourceVirtualEnvironmentVM")
	}
}

// TestDataSourceVirtualEnvironmentVMSchema tests the dataSourceVirtualEnvironmentVM schema.
func TestDataSourceVirtualEnvironmentVMSchema(t *testing.T) {
	t.Parallel()

	s := dataSourceVirtualEnvironmentVM()

	testComputedAttributes(t, s, []string{
		mkDataSourceVirtualEnvironmentVMName,
		mkDataSourceVirtualEnvironmentVMTags,
	})

	testValueTypes(t, s, map[string]schema.ValueType{
		mkDataSourceVirtualEnvironmentVMName:     schema.TypeString,
		mkDataSourceVirtualEnvironmentVMNodeName: schema.TypeString,
		mkDataSourceVirtualEnvironmentVMTags:     schema.TypeList,
		mkDataSourceVirtualEnvironmentVMVMID:     schema.TypeInt,
	})
}
