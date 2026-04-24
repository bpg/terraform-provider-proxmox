/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package memory

import (
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/bpg/terraform-provider-proxmox/fwprovider/attribute"
	"github.com/bpg/terraform-provider-proxmox/proxmox/nodes/vms"
)

// Model represents the memory configuration model.
//
// Mapping to Proxmox API:
//   - Size → memory (total available RAM)
//   - Balloon → balloon (guaranteed minimum RAM via balloon device; 0 disables balloon driver)
//   - Shares → shares (CPU scheduler priority for memory ballooning)
//   - Hugepages → hugepages (use hugepages for VM memory)
//   - KeepHugepages → keephugepages (don't release hugepages on shutdown)
type Model struct {
	Size          types.Int64  `tfsdk:"size"`
	Balloon       types.Int64  `tfsdk:"balloon"`
	Shares        types.Int64  `tfsdk:"shares"`
	Hugepages     types.String `tfsdk:"hugepages"`
	KeepHugepages types.Bool   `tfsdk:"keep_hugepages"`
}

func attributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"size":           types.Int64Type,
		"balloon":        types.Int64Type,
		"shares":         types.Int64Type,
		"hugepages":      types.StringType,
		"keep_hugepages": types.BoolType,
	}
}

// NullValue returns a properly typed null Value.
func NullValue() Value {
	return types.ObjectNull(attributeTypes())
}

// toAPI writes the memory-related fields onto the shared create/update body. Unlike the ADR-004
// reference shape (`toAPI() *SomeStruct`), memory has no dedicated API struct — its fields are
// independent top-level PVE keys. The signature therefore takes the body and populates in place,
// keeping the `toAPI()` naming contract while acknowledging the write-through shape.
// Null/unknown fields produce nil pointers so the request omits them entirely.
func (m *Model) toAPI(body *vms.CreateRequestBody) {
	if v := attribute.Int64PtrFromValue(m.Size); v != nil {
		n := int(*v)
		body.DedicatedMemory = &n
	}

	if v := attribute.Int64PtrFromValue(m.Balloon); v != nil {
		n := int(*v)
		body.FloatingMemory = &n
	}

	if v := attribute.Int64PtrFromValue(m.Shares); v != nil {
		n := int(*v)
		body.FloatingMemoryShares = &n
	}

	body.Hugepages = attribute.StringPtrFromValue(m.Hugepages)
	body.KeepHugepages = attribute.CustomBoolPtrFromValue(m.KeepHugepages)
}
