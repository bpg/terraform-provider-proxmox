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

	"github.com/bpg/terraform-provider-proxmox/proxmoxtf/test"
)

// TestAliasInstantiation tests whether the Alias instance can be instantiated.
func TestAliasInstantiation(t *testing.T) {
	t.Parallel()
	require.NotNilf(t, Alias(), "Cannot instantiate Alias")
}

// TestAliasSchema tests the Alias Schema.
func TestAliasSchema(t *testing.T) {
	t.Parallel()

	s := Alias().Schema

	test.AssertRequiredArguments(t, s, []string{
		mkAliasName,
		mkAliasCIDR,
	})

	test.AssertOptionalArguments(t, s, []string{
		mkSelectorVMID,
		mkSelectorNodeName,
		mkAliasComment,
	})

	test.AssertValueTypes(t, s, map[string]schema.ValueType{
		mkAliasName:    schema.TypeString,
		mkAliasCIDR:    schema.TypeString,
		mkAliasComment: schema.TypeString,
	})
}

// TestAliasImport tests the parsing of the alias import ID.
func TestAliasImport(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		id          string
		expected    map[string]any
		expectedErr string
	}{
		{"cluster", "cluster/local_network", map[string]any{
			mkAliasName: "local_network", mkSelectorNodeName: "", mkSelectorVMID: 0, mkSelectorContainerID: 0,
		}, ""},
		{"vm", "vm/pve/100/local_network", map[string]any{
			mkAliasName: "local_network", mkSelectorNodeName: "pve", mkSelectorVMID: 100, mkSelectorContainerID: 0,
		}, ""},
		{"container", "container/pve/101/local_network", map[string]any{
			mkAliasName: "local_network", mkSelectorNodeName: "pve", mkSelectorVMID: 0, mkSelectorContainerID: 101,
		}, ""},
		{"unknown prefix", "local_network", nil, "invalid import ID"},
		{"empty cluster name", "cluster/", nil, "alias name is empty"},
		{"vm without name", "vm/pve/100", nil, "invalid vm import ID format"},
		{"vm without node", "vm//100/local_network", nil, "invalid vm import ID format"},
		{"container with non numeric id", "container/pve/abc/local_network", nil, "container ID must be a number"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			d := Alias().TestResourceData()
			d.SetId(tt.id)

			res, err := aliasImport(t.Context(), d, nil)
			if tt.expectedErr != "" {
				require.ErrorContains(t, err, tt.expectedErr)
				return
			}

			require.NoError(t, err)
			require.Len(t, res, 1)
			require.Equal(t, tt.expected[mkAliasName], res[0].Id())

			for k, v := range tt.expected {
				require.Equal(t, v, res[0].Get(k), k)
			}
		})
	}
}
