/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package firewall

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/bpg/terraform-provider-proxmox/proxmoxtf/test"
)

// TestDataSourceVirtualEnvironmentAliasesInstantiation tests whether the DataSourceVirtualEnvironmentAliases instance can be instantiated.
func TestDataSourceVirtualEnvironmentAliasesInstantiation(t *testing.T) {
	s := DataSourceVirtualEnvironmentClusterAliases()

	if s == nil {
		t.Fatalf("Cannot instantiate DataSourceVirtualEnvironmentAliases")
	}
}

// TestDataSourceVirtualEnvironmentAliasesSchema tests the DataSourceVirtualEnvironmentAliases schema.
func TestDataSourceVirtualEnvironmentAliasesSchema(t *testing.T) {
	s := DataSourceVirtualEnvironmentClusterAliases()

	test.AssertComputedAttributes(t, s, []string{
		mkDataSourceVirtualEnvironmentClusterAliasesAliasIDs,
	})

	test.AssertValueTypes(t, s, map[string]schema.ValueType{
		mkDataSourceVirtualEnvironmentClusterAliasesAliasIDs: schema.TypeList,
	})
}
