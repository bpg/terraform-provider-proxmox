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

// TestAliasInstantiation tests whether the ResourceVirtualEnvironmentAlias instance can be instantiated.
func TestAliasInstantiation(t *testing.T) {
	t.Parallel()
	s := Alias()

	if s == nil {
		t.Fatalf("Cannot instantiate ResourceVirtualEnvironmentAlias")
	}
}

// TestAliasSchema tests the ResourceVirtualEnvironmentAlias schema.
func TestAliasSchema(t *testing.T) {
	t.Parallel()
	s := Alias()

	test.AssertRequiredArguments(t, s, []string{
		mkResourceVirtualEnvironmentClusterAliasName,
		mkResourceVirtualEnvironmentClusterAliasCIDR,
	})

	test.AssertOptionalArguments(t, s, []string{
		mkResourceVirtualEnvironmentClusterAliasComment,
	})

	test.AssertValueTypes(t, s, map[string]schema.ValueType{
		mkResourceVirtualEnvironmentClusterAliasName:    schema.TypeString,
		mkResourceVirtualEnvironmentClusterAliasCIDR:    schema.TypeString,
		mkResourceVirtualEnvironmentClusterAliasComment: schema.TypeString,
	})
}
