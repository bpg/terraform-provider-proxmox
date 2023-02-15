/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package proxmoxtf

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// TestDataSourceVirtualEnvironmentAliasesInstantiation tests whether the DataSourceVirtualEnvironmentAliases instance can be instantiated.
func TestDataSourceVirtualEnvironmentAliasesInstantiation(t *testing.T) {
	s := dataSourceVirtualEnvironmentClusterAliases()

	if s == nil {
		t.Fatalf("Cannot instantiate dataSourceVirtualEnvironmentAliases")
	}
}

// TestDataSourceVirtualEnvironmentAliasesSchema tests the dataSourceVirtualEnvironmentAliases schema.
func TestDataSourceVirtualEnvironmentAliasesSchema(t *testing.T) {
	s := dataSourceVirtualEnvironmentClusterAliases()

	testComputedAttributes(t, s, []string{
		mkDataSourceVirtualEnvironmentClusterAliasesAliasIDs,
	})

	testValueTypes(t, s, map[string]schema.ValueType{
		mkDataSourceVirtualEnvironmentClusterAliasesAliasIDs: schema.TypeList,
	})
}
