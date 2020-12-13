/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package proxmoxtf

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

// TestDataSourceVirtualEnvironmentTimeInstantiation tests whether the DataSourceVirtualEnvironmentRoles instance can be instantiated.
func TestDataSourceVirtualEnvironmentTimeInstantiation(t *testing.T) {
	s := dataSourceVirtualEnvironmentTime()

	if s == nil {
		t.Fatalf("Cannot instantiate dataSourceVirtualEnvironmentTime")
	}
}

// TestDataSourceVirtualEnvironmentTimeSchema tests the dataSourceVirtualEnvironmentTime schema.
func TestDataSourceVirtualEnvironmentTimeSchema(t *testing.T) {
	s := dataSourceVirtualEnvironmentTime()

	testRequiredArguments(t, s, []string{
		mkDataSourceVirtualEnvironmentTimeNodeName,
	})

	testComputedAttributes(t, s, []string{
		mkDataSourceVirtualEnvironmentTimeLocalTime,
		mkDataSourceVirtualEnvironmentTimeTimeZone,
		mkDataSourceVirtualEnvironmentTimeUTCTime,
	})

	testValueTypes(t, s, map[string]schema.ValueType{
		mkDataSourceVirtualEnvironmentTimeLocalTime: schema.TypeString,
		mkDataSourceVirtualEnvironmentTimeNodeName:  schema.TypeString,
		mkDataSourceVirtualEnvironmentTimeTimeZone:  schema.TypeString,
		mkDataSourceVirtualEnvironmentTimeUTCTime:   schema.TypeString,
	})
}
