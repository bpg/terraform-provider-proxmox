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

// TestUserInstantiation tests whether the User instance can be instantiated.
func TestUserInstantiation(t *testing.T) {
	t.Parallel()

	s := User()
	if s == nil {
		t.Fatalf("Cannot instantiate User")
	}
}

// TestUserSchema tests the User schema.
func TestUserSchema(t *testing.T) {
	t.Parallel()

	s := User().Schema

	test.AssertRequiredArguments(t, s, []string{
		mkResourceVirtualEnvironmentUserUserID,
	})

	test.AssertOptionalArguments(t, s, []string{
		mkResourceVirtualEnvironmentUserComment,
		mkResourceVirtualEnvironmentUserEmail,
		mkResourceVirtualEnvironmentUserEnabled,
		mkResourceVirtualEnvironmentUserExpirationDate,
		mkResourceVirtualEnvironmentUserFirstName,
		mkResourceVirtualEnvironmentUserGroups,
		mkResourceVirtualEnvironmentUserKeys,
		mkResourceVirtualEnvironmentUserLastName,
		mkResourceVirtualEnvironmentUserPassword,
	})

	test.AssertValueTypes(t, s, map[string]schema.ValueType{
		mkResourceVirtualEnvironmentUserComment:        schema.TypeString,
		mkResourceVirtualEnvironmentUserEmail:          schema.TypeString,
		mkResourceVirtualEnvironmentUserEnabled:        schema.TypeBool,
		mkResourceVirtualEnvironmentUserExpirationDate: schema.TypeString,
		mkResourceVirtualEnvironmentUserFirstName:      schema.TypeString,
		mkResourceVirtualEnvironmentUserGroups:         schema.TypeSet,
		mkResourceVirtualEnvironmentUserKeys:           schema.TypeString,
		mkResourceVirtualEnvironmentUserLastName:       schema.TypeString,
		mkResourceVirtualEnvironmentUserPassword:       schema.TypeString,
		mkResourceVirtualEnvironmentUserUserID:         schema.TypeString,
	})
}
