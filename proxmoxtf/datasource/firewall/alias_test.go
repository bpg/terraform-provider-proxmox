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

// TestAliasInstantiation tests whether the DataSourceVirtualEnvironmentAlias instance can be instantiated.
func TestAliasInstantiation(t *testing.T) {
	t.Parallel()
	s := DataSourceVirtualEnvironmentFirewallAlias()

	if s == nil {
		t.Fatalf("Cannot instantiate DataSourceVirtualEnvironmentAlias")
	}
}

// TestAliasSchema tests the DataSourceVirtualEnvironmentAlias schema.
func TestAliasSchema(t *testing.T) {
	t.Parallel()
	s := DataSourceVirtualEnvironmentFirewallAlias()

	test.AssertRequiredArguments(t, s, []string{
		mkDataSourceVirtualEnvironmentFirewallAliasName,
	})

	test.AssertComputedAttributes(t, s, []string{
		mkDataSourceVirtualEnvironmentFirewallAliasCIDR,
		mkDataSourceVirtualEnvironmentFirewallAliasComment,
	})

	test.AssertValueTypes(t, s, map[string]schema.ValueType{
		mkDataSourceVirtualEnvironmentFirewallAliasName:    schema.TypeString,
		mkDataSourceVirtualEnvironmentFirewallAliasCIDR:    schema.TypeString,
		mkDataSourceVirtualEnvironmentFirewallAliasComment: schema.TypeString,
	})
}
