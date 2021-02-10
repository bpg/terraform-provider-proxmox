/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package proxmoxtf

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

// TestDataSourceVirtualEnvironmentAliasInstantiation tests whether the DataSourceVirtualEnvironmentAlias instance can be instantiated.
func TestDataSourceVirtualEnvironmentAliasInstantiation(t *testing.T) {
	s := dataSourceVirtualEnvironmentClusterAlias()

	if s == nil {
		t.Fatalf("Cannot instantiate dataSourceVirtualEnvironmentAlias")
	}
}

// TestDataSourceVirtualEnvironmentAliasSchema tests the dataSourceVirtualEnvironmentAlias schema.
func TestDataSourceVirtualEnvironmentAliasSchema(t *testing.T) {
	s := dataSourceVirtualEnvironmentClusterAlias()

	testRequiredArguments(t, s, []string{
		mkDataSourceVirtualEnvironmentClusterAliasName,
	})

	testComputedAttributes(t, s, []string{
		mkDataSourceVirtualEnvironmentClusterAliasCIDR,
		mkDataSourceVirtualEnvironmentClusterAliasComment,
	})

	testValueTypes(t, s, map[string]schema.ValueType{
		mkDataSourceVirtualEnvironmentClusterAliasName:    schema.TypeString,
		mkDataSourceVirtualEnvironmentClusterAliasCIDR:    schema.TypeString,
		mkDataSourceVirtualEnvironmentClusterAliasComment: schema.TypeString,
	})
}
