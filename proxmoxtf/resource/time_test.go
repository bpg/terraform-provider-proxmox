/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package resource

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/bpg/terraform-provider-proxmox/proxmoxtf/test"
)

// TestTimeInstantiation tests whether the Time instance can be instantiated.
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

	s := Time()

	test.AssertRequiredArguments(t, s, []string{
		mkResourceVirtualEnvironmentTimeNodeName,
		mkResourceVirtualEnvironmentTimeTimeZone,
	})

	test.AssertComputedAttributes(t, s, []string{
		mkResourceVirtualEnvironmentTimeLocalTime,
		mkResourceVirtualEnvironmentTimeUTCTime,
	})

	test.AssertValueTypes(t, s, map[string]schema.ValueType{
		mkResourceVirtualEnvironmentTimeLocalTime: schema.TypeString,
		mkResourceVirtualEnvironmentTimeNodeName:  schema.TypeString,
		mkResourceVirtualEnvironmentTimeTimeZone:  schema.TypeString,
		mkResourceVirtualEnvironmentTimeUTCTime:   schema.TypeString,
	})
}
