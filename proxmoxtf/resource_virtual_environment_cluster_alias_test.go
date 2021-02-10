/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package proxmoxtf

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

// TestResourceVirtualEnvironmentAliasInstantiation tests whether the ResourceVirtualEnvironmentAlias instance can be instantiated.
func TestResourceVirtualEnvironmentAliasInstantiation(t *testing.T) {
	s := resourceVirtualEnvironmentClusterAlias()

	if s == nil {
		t.Fatalf("Cannot instantiate resourceVirtualEnvironmentAlias")
	}
}

// TestResourceVirtualEnvironmentAliasSchema tests the resourceVirtualEnvironmentAlias schema.
func TestResourceVirtualEnvironmentAliasSchema(t *testing.T) {
	s := resourceVirtualEnvironmentClusterAlias()

	testRequiredArguments(t, s, []string{
		mkResourceVirtualEnvironmentClusterAliasName,
		mkResourceVirtualEnvironmentClusterAliasCIDR,
	})

	testOptionalArguments(t, s, []string{
		mkResourceVirtualEnvironmentClusterAliasComment,
	})

	testValueTypes(t, s, map[string]schema.ValueType{
		mkResourceVirtualEnvironmentClusterAliasName:    schema.TypeString,
		mkResourceVirtualEnvironmentClusterAliasCIDR:    schema.TypeString,
		mkResourceVirtualEnvironmentClusterAliasComment: schema.TypeString,
	})
}
