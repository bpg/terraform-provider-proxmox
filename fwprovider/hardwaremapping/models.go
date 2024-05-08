/*
	This Source Code Form is subject to the terms of the Mozilla Public
	License, v. 2.0. If a copy of the MPL was not distributed with this
	file, You can obtain one at https://mozilla.org/MPL/2.0/.
*/

package hardwaremapping

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/types"

	customtypes "github.com/bpg/terraform-provider-proxmox/fwprovider/types/hardwaremapping"
	apitypes "github.com/bpg/terraform-provider-proxmox/proxmox/cluster/mapping"
	proxmoxtypes "github.com/bpg/terraform-provider-proxmox/proxmox/types/hardwaremapping"
)

const (
	// schemaAttrNameComment is the name of the schema attribute for the comment of a hardware mapping.
	// Note that the Proxmox VE API attribute is named "description", but we map it as a comment since this naming is
	// generally across the Proxmox VE web UI and API documentations. This still follows the [Terraform "best practices"]
	// as it improves the user experience by matching the field name to the naming used in the human-facing interfaces.
	//
	// [Terraform "best practices"]: https://developer.hashicorp.com/terraform/plugin/best-practices/hashicorp-provider-design-principles#resource-and-attribute-schema-should-closely-match-the-underlying-api
	schemaAttrNameComment = "comment"

	// schemaAttrNameMap is the name of the schema attribute for the map of a hardware mapping.
	schemaAttrNameMap = "map"

	// schemaAttrNameMapDeviceID is the name of the schema attribute for the device ID in a map of a hardware mapping.
	schemaAttrNameMapDeviceID = "id"

	// schemaAttrNameMapIOMMUGroup is the name of the schema attribute for the IOMMU group in a map of a hardware mapping.
	schemaAttrNameMapIOMMUGroup = "iommu_group"

	// schemaAttrNameMapNode is the name of the schema attribute for the node in a map of a hardware mapping.
	schemaAttrNameMapNode = "node"

	// schemaAttrNameMapPath is the name of the schema attribute for the path in a map of a hardware mapping.
	schemaAttrNameMapPath = "path"

	// schemaAttrNameMapSubsystemID is the name of the schema attribute for the subsystem ID in a map of a hardware
	// mapping.
	schemaAttrNameMapSubsystemID = "subsystem_id"

	// schemaAttrNameMediatedDevices is the name of the schema attribute for the mediated devices in a map of a hardware
	// mapping.
	schemaAttrNameMediatedDevices = "mediated_devices"

	// schemaAttrNameName is the name of the schema attribute for the name of a hardware mapping.
	schemaAttrNameName = "name"

	// schemaAttrNameTerraformID is the name of the schema attribute for the Terraform ID of a hardware mapping.
	schemaAttrNameTerraformID = "id"

	// schemaAttrNameType is the name of the schema attribute for the [proxmoxtypes.Type].
	schemaAttrNameType = "type"

	// schemaAttrNameCheckNode is the name of the schema attribute for the "check node" option of a
	// dataSource.
	schemaAttrNameCheckNode = "check_node"

	// schemaAttrNameChecks is the name of the schema attribute for the node checks diagnostics of a hardware mapping data
	// source.
	// Note that the Proxmox VE API attribute for [proxmoxtypes.TypeUSB] is named "errors", but we map it as "checks"
	// since this naming is generally across the Proxmox VE web UI and API documentations, including the attribute for
	// [proxmoxtypes.TypePCI].
	// This still follows the [Terraform "best practices"] as it improves the user experience by matching the field name
	// to the naming used in the human-facing interfaces.
	//
	// [Terraform "best practices"]: https://developer.hashicorp.com/terraform/plugin/best-practices/hashicorp-provider-design-principles#resource-and-attribute-schema-should-closely-match-the-underlying-api
	schemaAttrNameChecks = "checks"

	// schemaAttrNameChecksDiagsMappingID is the name of the schema attribute for a node check diagnostic mapping ID of a
	// dataSource.
	schemaAttrNameChecksDiagsMappingID = "mapping_id"

	// schemaAttrNameChecksDiagsMessage is the name of the schema attribute for a node check diagnostic message of a
	// dataSource.
	schemaAttrNameChecksDiagsMessage = "message"

	// schemaAttrNameChecksDiagsSeverity is the name of the schema attribute for a node check diagnostic severity of a
	// dataSource.
	schemaAttrNameChecksDiagsSeverity = "severity"

	// schemaAttrNameHWMIDs is the name of the schema attribute for the hardware mapping IDs of a
	// dataSource.
	schemaAttrNameHWMIDs = "ids"
)

// modelPCIMap maps the schema data for the map of a PCI hardware mapping.
type modelPCIMap struct {
	// Comment is the "comment" for the map.
	// This field is optional and is omitted by the Proxmox VE API when not set.
	// Note that the Proxmox VE API attribute is named "description", but we map it as a comment since this naming is
	// generally across the Proxmox VE web UI and API documentations. This still follows the [Terraform "best practices"]
	// as it improves the user experience by matching the field name to the naming used in the human-facing interfaces.
	//
	// [Terraform "best practices"]: https://developer.hashicorp.com/terraform/plugin/best-practices/hashicorp-provider-design-principles#resource-and-attribute-schema-should-closely-match-the-underlying-api
	Comment types.String `tfsdk:"comment"`

	// ID is the identifier of the map.
	ID types.String `tfsdk:"id"`

	// IOMMUGroup is the "IOMMU group" for the map.
	// This field is optional and is omitted by the Proxmox VE API when not set.
	IOMMUGroup types.Int64 `tfsdk:"iommu_group"`

	// Node is the "node name" for the map.
	Node types.String `tfsdk:"node"`

	// Path is the "path" for the map.
	Path customtypes.PathValue `tfsdk:"path"`

	// SubsystemID is the "subsystem ID" for the map.
	// This field is not mandatory for the Proxmox VE API call, but causes a PCI hardware mapping to be incomplete when
	// not set.
	SubsystemID types.String `tfsdk:"subsystem_id"`
}

// modelUSBMap maps the schema data for the map of a USB hardware mapping.
type modelUSBMap struct {
	// Comment is the "comment" for the map.
	// This field is optional and is omitted by the Proxmox VE API when not set.
	// Note that the Proxmox VE API attribute is named "description", but we map it as a comment since this naming is
	// generally across the Proxmox VE web UI and API documentations. This still follows the [Terraform "best practices"]
	// as it improves the user experience by matching the field name to the naming used in the human-facing interfaces.
	//
	// [Terraform "best practices"]: https://developer.hashicorp.com/terraform/plugin/best-practices/hashicorp-provider-design-principles#resource-and-attribute-schema-should-closely-match-the-underlying-api
	Comment types.String `tfsdk:"comment"`

	// ID is the identifier of the map.
	ID types.String `tfsdk:"id"`

	// Node is the "node name" for the map.
	Node types.String `tfsdk:"node"`

	// Path is the "path" for the map.
	Path customtypes.PathValue `tfsdk:"path"`
}

// modelPCI maps the schema data for a PCI hardware mapping.
type modelPCI struct {
	// Comment is the comment of the PCI hardware mapping.
	// Note that the Proxmox VE API attribute is named "description", but we map it as a comment since this naming is
	// generally across the Proxmox VE web UI and API documentations. This still follows the [Terraform "best practices"]
	// as it improves the user experience by matching the field name to the naming used in the human-facing interfaces.
	//
	// [Terraform "best practices"]: https://developer.hashicorp.com/terraform/plugin/best-practices/hashicorp-provider-design-principles#resource-and-attribute-schema-should-closely-match-the-underlying-api
	Comment types.String `tfsdk:"comment"`

	// ID is the Terraform identifier.
	ID types.String `tfsdk:"id"`

	// Name is the name of the PCI hardware mapping.
	Name types.String `tfsdk:"name"`

	// Map is the map of the PCI hardware mapping.
	Map []modelPCIMap `tfsdk:"map"`

	// MediatedDevices is the indicator for mediated devices of the PCI hardware mapping.
	MediatedDevices types.Bool `tfsdk:"mediated_devices"`
}

// modelUSB maps the schema data for a USB hardware mapping.
type modelUSB struct {
	// Comment is the comment of the USB hardware mapping.
	// Note that the Proxmox VE API attribute is named "description", but we map it as a comment since this naming is
	// generally across the Proxmox VE web UI and API documentations. This still follows the [Terraform "best practices"]
	// as it improves the user experience by matching the field name to the naming used in the human-facing interfaces.
	//
	// [Terraform "best practices"]: https://developer.hashicorp.com/terraform/plugin/best-practices/hashicorp-provider-design-principles#resource-and-attribute-schema-should-closely-match-the-underlying-api
	Comment types.String `tfsdk:"comment"`

	// ID is the Terraform identifier.
	ID types.String `tfsdk:"id"`

	// Name is the name of the USB hardware mapping.
	Name types.String `tfsdk:"name"`

	// Map is the map of the USB hardware mapping.
	Map []modelUSBMap `tfsdk:"map"`
}

// model maps the schema data for a hardware mappings data source.
type model struct {
	// Checks might contain relevant hardware mapping diagnostics about incorrect configurations for the node name set
	// defined by CheckNode.
	// Note that the Proxmox VE API attribute for [proxmoxtypes.TypeUSB] is named "errors", but we map it as "checks"
	// since this naming is generally across the Proxmox VE web UI and API documentations, including the attribute for
	// [proxmoxtypes.TypePCI].
	// Also note that the Proxmox VE API, for whatever reason, only returns one error at a time, even though the field is
	// an array.
	// This still follows the [Terraform "best practices"] as it improves the user experience by matching the field name
	// to the naming used in the human-facing interfaces.
	//
	// [Terraform "best practices"]: https://developer.hashicorp.com/terraform/plugin/best-practices/hashicorp-provider-design-principles#resource-and-attribute-schema-should-closely-match-the-underlying-api
	Checks []modelNodeCheckDiag `tfsdk:"checks"`

	// CheckNode is the name of the node whose configuration should be checked for correctness.
	CheckNode types.String `tfsdk:"check_node"`

	// ID is the Terraform identifier.
	ID types.String `tfsdk:"id"`

	// MappingIDs is the set of hardware mapping identifiers.
	MappingIDs types.Set `tfsdk:"ids"`

	// Type is the [proxmoxtypes.Type].
	Type types.String `tfsdk:"type"`
}

// modelNodeCheckDiag maps the schema data for hardware mapping node check diagnostic data.
type modelNodeCheckDiag struct {
	// MappingID is the corresponding hardware mapping ID of this node check diagnostic entry.
	MappingID types.String `tfsdk:"mapping_id"`

	// Message is the message of the node check diagnostic entry.
	Message types.String `tfsdk:"message"`

	// Severity is the severity of the node check diagnostic entry.
	Severity types.String `tfsdk:"severity"`
}

// importFromAPI imports the contents of a PCI hardware mapping model from the Proxmox VE API's response data.
func (hm *modelPCI) importFromAPI(_ context.Context, data *apitypes.GetResponseData) {
	// Ensure that both the ID and name are in sync.
	hm.Name = hm.ID
	// The attribute is named "description" by the Proxmox VE API, but we map it as a comment since this naming is
	// generally across the Proxmox VE web UI and API documentations.
	hm.Comment = types.StringPointerValue(data.Description)
	maps := make([]modelPCIMap, len(data.Map))

	for idx, pveMap := range data.Map {
		tfMap := modelPCIMap{
			ID:   pveMap.ID.ToValue(),
			Node: types.StringValue(pveMap.Node),
			Path: customtypes.NewPathPointerValue(pveMap.Path),
		}

		if pveMap.Description != nil {
			tfMap.Comment = types.StringPointerValue(pveMap.Description)
		}

		if pveMap.SubsystemID != "" {
			tfMap.SubsystemID = pveMap.SubsystemID.ToValue()
		}

		if pveMap.IOMMUGroup != nil {
			tfMap.IOMMUGroup = types.Int64Value(*pveMap.IOMMUGroup)
		}

		maps[idx] = tfMap
	}

	hm.MediatedDevices = data.MediatedDevices.ToValue()
	hm.Map = maps
}

// toCreateRequest builds the request data structure for creating a new PCI hardware mapping.
func (hm *modelPCI) toCreateRequest() *apitypes.CreateRequestBody {
	return &apitypes.CreateRequestBody{
		DataBase: hm.toRequestBase(),
		ID:       hm.ID.ValueString(),
	}
}

// toRequestBase builds the common request data structure for the PCI hardware mapping creation or update API calls.
func (hm *modelPCI) toRequestBase() apitypes.DataBase {
	dataBase := apitypes.DataBase{
		// The attribute is named "description" by the Proxmox VE API, but we map it as a comment since this naming is
		// generally across the Proxmox VE web UI and API documentations.
		Description: hm.Comment.ValueStringPointer(),
	}
	maps := make([]proxmoxtypes.Map, len(hm.Map))

	for idx, tfMap := range hm.Map {
		pveMap := proxmoxtypes.Map{
			// The attribute is named "description" by the Proxmox VE API, but we map it as a comment since this naming is
			// generally across the Proxmox VE web UI and API documentations.
			Description: tfMap.Comment.ValueStringPointer(),
			ID:          proxmoxtypes.DeviceID(tfMap.ID.ValueString()),
			IOMMUGroup:  tfMap.IOMMUGroup.ValueInt64Pointer(),
			Node:        tfMap.Node.ValueString(),
			Path:        tfMap.Path.ValueStringPointer(),
			SubsystemID: proxmoxtypes.DeviceID(tfMap.SubsystemID.ValueString()),
		}
		maps[idx] = pveMap
	}

	dataBase.Map = maps
	dataBase.MediatedDevices.FromValue(hm.MediatedDevices)

	return dataBase
}

// toUpdateRequest builds the request data structure for updating an existing PCI hardware mapping.
func (hm *modelPCI) toUpdateRequest(currentState *modelPCI) *apitypes.UpdateRequestBody {
	var del []string

	baseRequest := hm.toRequestBase()

	if hm.Comment.IsNull() && !currentState.Comment.IsNull() {
		// The Proxmox VE API attribute is named "description" while we name it "comment" internally since this naming is
		// generally used across the Proxmox VE web UI and API documentations.
		// This still follows theTerraform "best practices" [1] as it improves the user experience by matching the field
		// name to the naming used in the human-facing interfaces.
		// References:
		//   1. https://developer.hashicorp.com/terraform/plugin/best-practices/hashicorp-provider-design-principles#resource-and-attribute-schema-should-closely-match-the-underlying-api
		del = append(del, proxmoxtypes.AttrNameDescription)
	}

	if hm.MediatedDevices.IsNull() || !hm.MediatedDevices.ValueBool() {
		del = append(del, apitypes.APIParamNamePCIMediatedDevices)

		baseRequest.MediatedDevices.FromValue(types.BoolValue(false))
	}

	return &apitypes.UpdateRequestBody{
		DataBase: baseRequest,
		Delete:   del,
	}
}

// importFromAPI imports the contents of a USB hardware mapping model from the Proxmox VE API's response data.
func (hm *modelUSB) importFromAPI(_ context.Context, data *apitypes.GetResponseData) {
	// Ensure that both the ID and name are in sync.
	hm.Name = hm.ID
	// The attribute is named "description" by the Proxmox VE API, but we map it as a comment since this naming is
	// generally across the Proxmox VE web UI and API documentations.
	hm.Comment = types.StringPointerValue(data.Description)
	maps := make([]modelUSBMap, len(data.Map))

	for idx, pveMap := range data.Map {
		tfMap := modelUSBMap{
			ID:   pveMap.ID.ToValue(),
			Node: types.StringValue(pveMap.Node),
			Path: customtypes.NewPathPointerValue(pveMap.Path),
		}
		if pveMap.Description != nil {
			tfMap.Comment = types.StringPointerValue(pveMap.Description)
		}

		maps[idx] = tfMap
	}

	hm.Map = maps
}

// toCreateRequest builds the request data structure for creating a new USB hardware mapping.
func (hm *modelUSB) toCreateRequest() *apitypes.CreateRequestBody {
	return &apitypes.CreateRequestBody{
		DataBase: hm.toRequestBase(),
		ID:       hm.ID.ValueString(),
	}
}

// toRequestBase builds the common request data structure for the USB hardware mapping creation or update API calls.
func (hm *modelUSB) toRequestBase() apitypes.DataBase {
	dataBase := apitypes.DataBase{
		// The attribute is named "description" by the Proxmox VE API, but we map it as a comment since this naming is
		// generally across the Proxmox VE web UI and API documentations.
		Description: hm.Comment.ValueStringPointer(),
	}
	maps := make([]proxmoxtypes.Map, len(hm.Map))

	for idx, tfMap := range hm.Map {
		pveMap := proxmoxtypes.Map{
			ID:   proxmoxtypes.DeviceID(tfMap.ID.ValueString()),
			Node: tfMap.Node.ValueString(),
			Path: tfMap.Path.ValueStringPointer(),
		}
		if !tfMap.Comment.IsNull() {
			// The attribute is named "description" by the Proxmox VE API, but we map it as a comment since this naming is
			// generally across the Proxmox VE web UI and API documentations.
			pveMap.Description = tfMap.Comment.ValueStringPointer()
		}

		maps[idx] = pveMap
	}

	dataBase.Map = maps

	return dataBase
}

// toUpdateRequest builds the request data structure for updating an existing USB hardware mapping.
func (hm *modelUSB) toUpdateRequest(currentState *modelUSB) *apitypes.UpdateRequestBody {
	var del []string

	if hm.Comment.IsNull() && !currentState.Comment.IsNull() {
		// The Proxmox VE API attribute is named "description" while we name it "comment" internally since this naming is
		// generally used across the Proxmox VE web UI and API documentations.
		// This still follows the Terraform "best practices" [1] as it improves the user experience by matching the field
		// name to the naming used in the human-facing interfaces.
		// References:
		//   1. https://developer.hashicorp.com/terraform/plugin/best-practices/hashicorp-provider-design-principles#resource-and-attribute-schema-should-closely-match-the-underlying-api
		del = append(del, proxmoxtypes.AttrNameDescription)
	}

	return &apitypes.UpdateRequestBody{
		DataBase: hm.toRequestBase(),
		Delete:   del,
	}
}
