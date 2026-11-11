/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package resource

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestVMCPUTypeSchema validates cpu.type through VM().Schema. Model names are
// not allowlisted: PVE accepts or rejects them. The cases cover the legacy
// names, every model reported missing on PVE 9.2.4, a custom model, and a
// synthetic future name.
func TestVMCPUTypeSchema(t *testing.T) {
	t.Parallel()

	vmSchema := VM().Schema

	cpuTypes := []string{
		"host",
		"qemu64",
		"EPYC-Genoa",
		"custom-abc",
		"FutureArch-v9",
		"ClearwaterForest",
		"ClearwaterForest-v2",
		"ClearwaterForest-v3",
		"DiamondRapids",
		"EPYC-Genoa-v2",
		"EPYC-Milan-v2",
		"EPYC-Milan-v3",
		"EPYC-Rome-v5",
		"EPYC-Turin",
		"EPYC-v5",
		"GraniteRapids",
		"GraniteRapids-v2",
		"GraniteRapids-v3",
		"GraniteRapids-v4",
		"GraniteRapids-v5",
		"Icelake-Server-v7",
		"SapphireRapids-v2",
		"SapphireRapids-v3",
		"SapphireRapids-v4",
		"SapphireRapids-v5",
		"SapphireRapids-v6",
		"SierraForest",
		"SierraForest-v2",
		"SierraForest-v3",
		"SierraForest-v4",
		"SierraForest-v5",
	}

	for _, cpuType := range cpuTypes {
		t.Run(cpuType, func(t *testing.T) {
			t.Parallel()

			diags := schema.InternalMap(vmSchema).Validate(terraform.NewResourceConfigRaw(map[string]any{
				mkNodeName: "pve",
				mkCPU: []any{
					map[string]any{
						mkCPUType: cpuType,
					},
				},
			}))
			require.False(t, diags.HasError(), "cpu.type %q: %v", cpuType, diags)
		})
	}

	cpuBlock, ok := vmSchema[mkCPU].Elem.(*schema.Resource)
	require.True(t, ok)

	typeSchema := cpuBlock.Schema[mkCPUType]
	require.Equal(t, schema.TypeString, typeSchema.Type)
	require.True(t, typeSchema.Optional)
	require.False(t, typeSchema.Required)
	require.Equal(t, "qemu64", typeSchema.Default)
	require.Nil(t, typeSchema.ValidateDiagFunc)

	resourceData := schema.TestResourceDataRaw(t, vmSchema, map[string]any{
		mkNodeName: "pve",
		mkCPU: []any{
			map[string]any{},
		},
	})
	require.Equal(t, "qemu64", resourceData.Get(mkCPU+".0."+mkCPUType))
}

func TestMachineType(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		value string
		valid bool
	}{
		{"empty is valid", "", true},
		{"invalid", "invalid", false},
		{"valid q35", "q35", true},
		{"valid q35 with viommu", "q35,viommu=virtio", true},
		{"invalid q35 with viommu", "q35,viommu=invalid", false},
		{"valid pc-q35", "pc-q35-2.3", true},
		{"valid i440fx", "pc-i440fx-3.1+pve0", true},
		{"valid virt", "virt", true},
		{"invalid i440fx", "i440fx", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			f := MachineTypeValidator()
			res := f(tt.value, nil)

			if tt.valid {
				require.Empty(t, res, "validate: '%s'", tt.value)
			} else {
				require.NotEmpty(t, res, "validate: '%s'", tt.value)
			}
		})
	}
}

func TestVmHostname(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		value string
		valid bool
	}{
		{"empty", "", true},
		{"underscores", "my_name", false},
		{"trailing dot", "my-name.com.", false},
		{"starts with alphanumeric", "-my-name.com", false},
		{"ends with alphanumeric", "my-name.com!", false},
		{"single letter", "a", true},
		{"domain name", "my-name.com", true},
		{"multi domain", "my-name.com.edu.net.xyz.dev", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			f := HostnameValidator()
			res := f(tt.value, nil)

			valid := !res.HasError()
			assert.Equal(t, tt.valid, valid)
		})
	}
}
