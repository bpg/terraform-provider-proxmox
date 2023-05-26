/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package datasource

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/bpg/terraform-provider-proxmox/proxmoxtf/test"
)

// TestVersionInstantiation tests whether the Version instance can be instantiated.
func TestVersionInstantiation(t *testing.T) {
	t.Parallel()

	s := Version()
	if s == nil {
		t.Fatalf("Cannot instantiate Version")
	}
}

// TestVersionSchema tests the Version schema.
func TestVersionSchema(t *testing.T) {
	t.Parallel()

	s := Version()

	test.AssertComputedAttributes(t, s, []string{
		mkDataSourceVirtualEnvironmentVersionKeyboardLayout,
		mkDataSourceVirtualEnvironmentVersionRelease,
		mkDataSourceVirtualEnvironmentVersionRepositoryID,
		mkDataSourceVirtualEnvironmentVersionVersion,
	})

	test.AssertValueTypes(t, s, map[string]schema.ValueType{
		mkDataSourceVirtualEnvironmentVersionKeyboardLayout: schema.TypeString,
		mkDataSourceVirtualEnvironmentVersionRelease:        schema.TypeString,
		mkDataSourceVirtualEnvironmentVersionRepositoryID:   schema.TypeString,
		mkDataSourceVirtualEnvironmentVersionVersion:        schema.TypeString,
	})
}
