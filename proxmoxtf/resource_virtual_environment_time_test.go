/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package proxmoxtf

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// TestResourceVirtualEnvironmentTimeInstantiation tests whether the ResourceVirtualEnvironmentTime instance can be instantiated.
func TestResourceVirtualEnvironmentTimeInstantiation(t *testing.T) {
	s := resourceVirtualEnvironmentTime()

	if s == nil {
		t.Fatalf("Cannot instantiate resourceVirtualEnvironmentTime")
	}
}

// TestResourceVirtualEnvironmentTimeSchema tests the resourceVirtualEnvironmentTime schema.
func TestResourceVirtualEnvironmentTimeSchema(t *testing.T) {
	s := resourceVirtualEnvironmentTime()

	testRequiredArguments(t, s, []string{
		mkResourceVirtualEnvironmentTimeNodeName,
		mkResourceVirtualEnvironmentTimeTimeZone,
	})

	testComputedAttributes(t, s, []string{
		mkResourceVirtualEnvironmentTimeLocalTime,
		mkResourceVirtualEnvironmentTimeUTCTime,
	})

	testValueTypes(t, s, map[string]schema.ValueType{
		mkResourceVirtualEnvironmentTimeLocalTime: schema.TypeString,
		mkResourceVirtualEnvironmentTimeNodeName:  schema.TypeString,
		mkResourceVirtualEnvironmentTimeTimeZone:  schema.TypeString,
		mkResourceVirtualEnvironmentTimeUTCTime:   schema.TypeString,
	})
}
