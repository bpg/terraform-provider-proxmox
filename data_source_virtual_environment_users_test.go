/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package main

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

	testSchemaValueTypes(t, s, []string{
		mkDataSourceVirtualEnvironmentUsersComments,
		mkDataSourceVirtualEnvironmentUsersEmails,
		mkDataSourceVirtualEnvironmentUsersEnabled,
		mkDataSourceVirtualEnvironmentUsersExpirationDates,
		mkDataSourceVirtualEnvironmentUsersFirstNames,
		mkDataSourceVirtualEnvironmentUsersGroups,
		mkDataSourceVirtualEnvironmentUsersKeys,
		mkDataSourceVirtualEnvironmentUsersLastNames,
		mkDataSourceVirtualEnvironmentUsersUserIDs,
	}, []schema.ValueType{
		schema.TypeList,
		schema.TypeList,
		schema.TypeList,
		schema.TypeList,
		schema.TypeList,
		schema.TypeList,
		schema.TypeList,
		schema.TypeList,
		schema.TypeList,
	})
}
