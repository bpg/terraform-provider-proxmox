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

// TestResourceVirtualEnvironmentTimeInstantiation tests whether the ResourceVirtualEnvironmentTime instance can be instantiated.
func TestResourceVirtualEnvironmentTimeInstantiation(t *testing.T) {
	s := ResourceVirtualEnvironmentTime()

	if s == nil {
		t.Fatalf("Cannot instantiate ResourceVirtualEnvironmentTime")
	}
}

// TestResourceVirtualEnvironmentTimeSchema tests the ResourceVirtualEnvironmentTime schema.
func TestResourceVirtualEnvironmentTimeSchema(t *testing.T) {
	s := ResourceVirtualEnvironmentTime()

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
