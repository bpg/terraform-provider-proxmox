/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package firewall

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/require"

	"github.com/bpg/terraform-provider-proxmox/proxmoxtf/resource/firewall"
	"github.com/bpg/terraform-provider-proxmox/proxmoxtf/test"
)

// TestSecurityGroupInstantiation tests whether the SecurityGroup instance can be instantiated.
func TestSecurityGroupInstantiation(t *testing.T) {
	t.Parallel()
	require.NotNilf(t, SecurityGroup(), "Cannot instantiate SecurityGroup")
}

// TestSecurityGroupSchema tests the SecurityGroup Schema.
func TestSecurityGroupSchema(t *testing.T) {
	t.Parallel()

	s := SecurityGroup().Schema

	test.AssertRequiredArguments(t, s, []string{
		mkSecurityGroupName,
	})

	test.AssertOptionalArguments(t, s, []string{
		mkSecurityGroupComment,
	})

	test.AssertValueTypes(t, s, map[string]schema.ValueType{
		mkSecurityGroupName:    schema.TypeString,
		mkSecurityGroupComment: schema.TypeString,
	})

	test.AssertNestedSchemaExistence(t, s, firewall.MkRule)
}

// TestSecurityGroupNameValidation tests that security group names with uppercase letters are rejected.
func TestSecurityGroupNameValidation(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		value   any
		wantErr string
	}{
		{"lowercase", "mytestfilter", ""},
		{"lowercase with digits and dash", "web-servers-01", ""},
		{"capital letters", "MyTestFilter", `use "mytestfilter" instead`},
		{"single capital letter", "webServers", `use "webservers" instead`},
		{"not a string", 42, "to be string"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			_, errs := validateSecurityGroupName(tt.value, mkSecurityGroupName)

			if tt.wantErr == "" {
				require.Empty(t, errs)

				return
			}

			require.Len(t, errs, 1)
			require.ErrorContains(t, errs[0], tt.wantErr)
		})
	}
}

// TestSecurityGroupNameSchemaRejectsUppercase tests that the name attribute is wired to the validator.
func TestSecurityGroupNameSchemaRejectsUppercase(t *testing.T) {
	t.Parallel()

	nameSchema := SecurityGroup().Schema[mkSecurityGroupName]
	require.NotNil(t, nameSchema.ValidateDiagFunc)

	require.False(t, nameSchema.ValidateDiagFunc("mytestfilter", nil).HasError())
	require.True(t, nameSchema.ValidateDiagFunc("MyTestFilter", nil).HasError())
}
