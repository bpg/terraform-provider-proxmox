/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package proxmoxtf

import (
	"testing"

	"github.com/hashicorp/terraform/helper/schema"
)

// TestDataSourceVirtualEnvironmentUsersInstantiation tests whether the DataSourceVirtualEnvironmentUsers instance can be instantiated.
func TestDataSourceVirtualEnvironmentUsersInstantiation(t *testing.T) {
	s := dataSourceVirtualEnvironmentUsers()

	if s == nil {
		t.Fatalf("Cannot instantiate dataSourceVirtualEnvironmentUsers")
	}
}

// TestDataSourceVirtualEnvironmentUsersSchema tests the dataSourceVirtualEnvironmentUsers schema.
func TestDataSourceVirtualEnvironmentUsersSchema(t *testing.T) {
	s := dataSourceVirtualEnvironmentUsers()

	testComputedAttributes(t, s, []string{
		mkDataSourceVirtualEnvironmentUsersComments,
		mkDataSourceVirtualEnvironmentUsersEmails,
		mkDataSourceVirtualEnvironmentUsersEnabled,
		mkDataSourceVirtualEnvironmentUsersExpirationDates,
		mkDataSourceVirtualEnvironmentUsersFirstNames,
		mkDataSourceVirtualEnvironmentUsersGroups,
		mkDataSourceVirtualEnvironmentUsersKeys,
		mkDataSourceVirtualEnvironmentUsersLastNames,
		mkDataSourceVirtualEnvironmentUsersUserIDs,
	})

	testValueTypes(t, s, map[string]schema.ValueType{
		mkDataSourceVirtualEnvironmentUsersComments:        schema.TypeList,
		mkDataSourceVirtualEnvironmentUsersEmails:          schema.TypeList,
		mkDataSourceVirtualEnvironmentUsersEnabled:         schema.TypeList,
		mkDataSourceVirtualEnvironmentUsersExpirationDates: schema.TypeList,
		mkDataSourceVirtualEnvironmentUsersFirstNames:      schema.TypeList,
		mkDataSourceVirtualEnvironmentUsersGroups:          schema.TypeList,
		mkDataSourceVirtualEnvironmentUsersKeys:            schema.TypeList,
		mkDataSourceVirtualEnvironmentUsersLastNames:       schema.TypeList,
		mkDataSourceVirtualEnvironmentUsersUserIDs:         schema.TypeList,
	})
}
