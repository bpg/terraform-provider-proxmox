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

// TestAliasesInstantiation tests whether the DataSourceVirtualEnvironmentAliases instance can be instantiated.
func TestAliasesInstantiation(t *testing.T) {
	t.Parallel()
	s := Aliases()

	if s == nil {
		t.Fatalf("Cannot instantiate DataSourceVirtualEnvironmentAliases")
	}
}

// TestAliasesSchema tests the DataSourceVirtualEnvironmentAliases schema.
func TestAliasesSchema(t *testing.T) {
	t.Parallel()
	s := Aliases()

	test.AssertComputedAttributes(t, s, []string{
		mkDataSourceVirtualEnvironmentClusterAliasesAliasNames,
	})

	test.AssertValueTypes(t, s, map[string]schema.ValueType{
		mkDataSourceVirtualEnvironmentClusterAliasesAliasNames: schema.TypeList,
	})
}
