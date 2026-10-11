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

// TestGroupInstantiation tests whether the Group instance can be instantiated.
func TestGroupInstantiation(t *testing.T) {
	t.Parallel()

	s := Group()
	if s == nil {
		t.Fatalf("Cannot instantiate Group")
	}
}

// TestGroupSchema tests the Group schema.
func TestGroupSchema(t *testing.T) {
	t.Parallel()

	s := Group().Schema

	test.AssertRequiredArguments(t, s, []string{
		mkResourceVirtualEnvironmentGroupID,
	})

	test.AssertOptionalArguments(t, s, []string{
		mkResourceVirtualEnvironmentGroupACL,
		mkResourceVirtualEnvironmentGroupComment,
	})

	test.AssertComputedAttributes(t, s, []string{
		mkResourceVirtualEnvironmentGroupMembers,
	})

	test.AssertValueTypes(t, s, map[string]schema.ValueType{
		mkResourceVirtualEnvironmentGroupACL:     schema.TypeSet,
		mkResourceVirtualEnvironmentGroupComment: schema.TypeString,
		mkResourceVirtualEnvironmentGroupID:      schema.TypeString,
		mkResourceVirtualEnvironmentGroupMembers: schema.TypeSet,
	})

	aclSchema := test.AssertNestedSchemaExistence(t, s, mkResourceVirtualEnvironmentGroupACL)

	test.AssertRequiredArguments(t, aclSchema, []string{
		mkResourceVirtualEnvironmentGroupACLPath,
		mkResourceVirtualEnvironmentGroupACLRoleID,
	})

	test.AssertOptionalArguments(t, aclSchema, []string{
		mkResourceVirtualEnvironmentGroupACLPropagate,
	})

	test.AssertValueTypes(t, aclSchema, map[string]schema.ValueType{
		mkResourceVirtualEnvironmentGroupACLPath:      schema.TypeString,
		mkResourceVirtualEnvironmentGroupACLPropagate: schema.TypeBool,
		mkResourceVirtualEnvironmentGroupACLRoleID:    schema.TypeString,
	})
}

// TestGroupACLDeprecationWarning verifies the deprecated `acl` block warns only when it is
// actually set in config, not on every group that omits it.
func TestGroupACLDeprecationWarning(t *testing.T) {
	t.Parallel()

	s := Group().Schema

	if hasACLDeprecationWarning(s, map[string]any{
		mkResourceVirtualEnvironmentGroupID: "test-group",
	}) {
		t.Error("unexpected acl deprecation warning when no acl block is configured")
	}

	if !hasACLDeprecationWarning(s, map[string]any{
		mkResourceVirtualEnvironmentGroupID: "test-group",
		mkResourceVirtualEnvironmentGroupACL: []any{
			map[string]any{
				mkResourceVirtualEnvironmentGroupACLPath:   "/",
				mkResourceVirtualEnvironmentGroupACLRoleID: "Administrator",
			},
		},
	}) {
		t.Error("expected acl deprecation warning when an acl block is configured")
	}
}
