/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package proxmoxtf

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// TestDataSourceVirtualEnvironmentPoolsInstantiation tests whether the DataSourceVirtualEnvironmentPools instance can be instantiated.
func TestDataSourceVirtualEnvironmentPoolsInstantiation(t *testing.T) {
	s := dataSourceVirtualEnvironmentPools()

	if s == nil {
		t.Fatalf("Cannot instantiate dataSourceVirtualEnvironmentPools")
	}
}

// TestDataSourceVirtualEnvironmentPoolsSchema tests the dataSourceVirtualEnvironmentPools schema.
func TestDataSourceVirtualEnvironmentPoolsSchema(t *testing.T) {
	s := dataSourceVirtualEnvironmentPools()

	testComputedAttributes(t, s, []string{
		mkDataSourceVirtualEnvironmentPoolsPoolIDs,
	})

	testValueTypes(t, s, map[string]schema.ValueType{
		mkDataSourceVirtualEnvironmentPoolsPoolIDs: schema.TypeList,
	})
}
