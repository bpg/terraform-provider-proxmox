/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package resource

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/require"
)

func TestVMReadSMBIOSState(t *testing.T) {
	t.Parallel()

	cloneBlock := []any{map[string]any{"vm_id": 100}}
	configuredBlock := []any{map[string]any{mkSMBIOSSKU: "test-sku"}}

	// smbios as decoded from the Proxmox config, with defaults for the unset fields.
	pveSMBIOS := func(sku string) map[string]any {
		return map[string]any{
			mkSMBIOSFamily:       dvSMBIOSFamily,
			mkSMBIOSManufacturer: dvSMBIOSManufacturer,
			mkSMBIOSProduct:      dvSMBIOSProduct,
			mkSMBIOSSerial:       dvSMBIOSSerial,
			mkSMBIOSSKU:          sku,
			mkSMBIOSVersion:      dvSMBIOSVersion,
			mkSMBIOSUUID:         "11111111-2222-3333-4444-555555555555",
		}
	}

	tests := []struct {
		name    string
		raw     map[string]any
		pve     map[string]any
		wantSKU *string // nil means the block is expected to be absent
	}{
		{
			name:    "clone with configured block refreshes state from Proxmox when SKU is missing",
			raw:     map[string]any{mkClone: cloneBlock, mkSMBIOS: configuredBlock},
			pve:     pveSMBIOS(""),
			wantSKU: ptr(""),
		},
		{
			name:    "clone with configured block keeps matching SKU",
			raw:     map[string]any{mkClone: cloneBlock, mkSMBIOS: configuredBlock},
			pve:     pveSMBIOS("test-sku"),
			wantSKU: ptr("test-sku"),
		},
		{
			name: "clone with configured block and no SMBIOS in Proxmox drops the block from state",
			raw:  map[string]any{mkClone: cloneBlock, mkSMBIOS: configuredBlock},
			pve:  nil,
		},
		{
			name: "clone without configured block ignores template SMBIOS",
			raw:  map[string]any{mkClone: cloneBlock},
			pve:  pveSMBIOS("from-template"),
		},
		{
			name:    "no clone refreshes state from Proxmox",
			raw:     map[string]any{mkSMBIOS: configuredBlock},
			pve:     pveSMBIOS(""),
			wantSKU: ptr(""),
		},
		{
			name: "no clone and no SMBIOS in Proxmox leaves the block absent",
			raw:  map[string]any{},
			pve:  nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			d := schema.TestResourceDataRaw(t, VM().Schema, tt.raw)
			clone := d.Get(mkClone).([]any)

			diags := vmReadSMBIOSState(d, tt.pve, clone)
			require.False(t, diags.HasError(), "unexpected diagnostics: %v", diags)

			got := d.Get(mkSMBIOS).([]any)
			if tt.wantSKU == nil {
				require.Empty(t, got)

				return
			}

			require.Len(t, got, 1)
			require.Equal(t, *tt.wantSKU, got[0].(map[string]any)[mkSMBIOSSKU])
		})
	}
}

func ptr[T any](v T) *T {
	return &v
}
