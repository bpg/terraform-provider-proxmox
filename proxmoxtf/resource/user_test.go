/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package resource

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/bpg/terraform-provider-proxmox/proxmoxtf/test"
)

// hasACLDeprecationWarning reports whether validating the given config against the
// resource schema produces an "Argument is deprecated" warning for the `acl` attribute.
func hasACLDeprecationWarning(s map[string]*schema.Schema, config map[string]any) bool {
	diags := schema.InternalMap(s).Validate(terraform.NewResourceConfigRaw(config))

	for _, d := range diags {
		if d.Severity == diag.Warning && d.Summary == "Argument is deprecated" {
			return true
		}
	}

	return false
}

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
		mkResourceVirtualEnvironmentUserACL,
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
		mkResourceVirtualEnvironmentUserACL:            schema.TypeSet,
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

	aclSchema := test.AssertNestedSchemaExistence(t, s, mkResourceVirtualEnvironmentUserACL)

	test.AssertRequiredArguments(t, aclSchema, []string{
		mkResourceVirtualEnvironmentUserACLPath,
		mkResourceVirtualEnvironmentUserACLRoleID,
	})

	test.AssertOptionalArguments(t, aclSchema, []string{
		mkResourceVirtualEnvironmentUserACLPropagate,
	})

	test.AssertValueTypes(t, aclSchema, map[string]schema.ValueType{
		mkResourceVirtualEnvironmentUserACLPath:      schema.TypeString,
		mkResourceVirtualEnvironmentUserACLPropagate: schema.TypeBool,
		mkResourceVirtualEnvironmentUserACLRoleID:    schema.TypeString,
	})
}

// TestUserACLDeprecationWarning verifies the deprecated `acl` block warns only when it is
// actually set in config, not on every user that omits it.
func TestUserACLDeprecationWarning(t *testing.T) {
	t.Parallel()

	s := User().Schema

	if hasACLDeprecationWarning(s, map[string]any{
		mkResourceVirtualEnvironmentUserUserID: "test@pve",
	}) {
		t.Error("unexpected acl deprecation warning when no acl block is configured")
	}

	if !hasACLDeprecationWarning(s, map[string]any{
		mkResourceVirtualEnvironmentUserUserID: "test@pve",
		mkResourceVirtualEnvironmentUserACL: []any{
			map[string]any{
				mkResourceVirtualEnvironmentUserACLPath:   "/",
				mkResourceVirtualEnvironmentUserACLRoleID: "Administrator",
			},
		},
	}) {
		t.Error("expected acl deprecation warning when an acl block is configured")
	}
}
