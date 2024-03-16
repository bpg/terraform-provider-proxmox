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

// TestTimeInstantiation tests whether the Roles instance can be instantiated.
func TestTimeInstantiation(t *testing.T) {
	t.Parallel()

	s := Time()
	if s == nil {
		t.Fatalf("Cannot instantiate Time")
	}
}

// TestTimeSchema tests the Time schema.
func TestTimeSchema(t *testing.T) {
	t.Parallel()

	s := Time().Schema

	test.AssertRequiredArguments(t, s, []string{
		mkDataSourceVirtualEnvironmentTimeNodeName,
	})

	test.AssertComputedAttributes(t, s, []string{
		mkDataSourceVirtualEnvironmentTimeLocalTime,
		mkDataSourceVirtualEnvironmentTimeTimeZone,
		mkDataSourceVirtualEnvironmentTimeUTCTime,
	})

	test.AssertValueTypes(t, s, map[string]schema.ValueType{
		mkDataSourceVirtualEnvironmentTimeLocalTime: schema.TypeString,
		mkDataSourceVirtualEnvironmentTimeNodeName:  schema.TypeString,
		mkDataSourceVirtualEnvironmentTimeTimeZone:  schema.TypeString,
		mkDataSourceVirtualEnvironmentTimeUTCTime:   schema.TypeString,
	})
}
