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

// TestUsersInstantiation tests whether the Users instance can be instantiated.
func TestUsersInstantiation(t *testing.T) {
	t.Parallel()

	s := Users()
	if s == nil {
		t.Fatalf("Cannot instantiate Users")
	}
}

// TestUsersSchema tests the Users schema.
func TestUsersSchema(t *testing.T) {
	t.Parallel()

	s := Users()

	test.AssertComputedAttributes(t, s, []string{
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

	test.AssertValueTypes(t, s, map[string]schema.ValueType{
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
