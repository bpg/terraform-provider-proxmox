/*
	This Source Code Form is subject to the terms of the Mozilla Public
	License, v. 2.0. If a copy of the MPL was not distributed with this
	file, You can obtain one at https://mozilla.org/MPL/2.0/.
*/

package fwprovider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/types"

	customtypes "github.com/bpg/terraform-provider-proxmox/fwprovider/types"
	apitypes "github.com/bpg/terraform-provider-proxmox/proxmox/cluster/mapping"
	proxmoxtypes "github.com/bpg/terraform-provider-proxmox/proxmox/types"
)

const (
	// hardwareMappingSchemaAttrNameComment is the name of the schema attribute for the comment of a hardware mapping.
	// Note that the Proxmox VE API attribute is named "description", but we map it as a comment since this naming is
	// generally across the Proxmox VE web UI and
	// API documentations. This still follows the Terraform "best practices" as it improves the user experience by
	// matching the field name to the naming used in
	// the human-facing interfaces.
	// References:
	//   1. https://developer.hashicorp.com/terraform/plugin/best-practices/hashicorp-provider-design-principles#resource-and-attribute-schema-should-closely-match-the-underlying-api
	hardwareMappingSchemaAttrNameComment = "comment"

	// hardwareMappingSchemaAttrNameMap is the name of the schema attribute for the map of a hardware mapping.
	hardwareMappingSchemaAttrNameMap = "map"

	// hardwareMappingSchemaAttrNameMapDeviceID is the name of the schema attribute for the device ID in a map of a
	// hardware mapping.
	hardwareMappingSchemaAttrNameMapDeviceID = "id"

	// hardwareMappingSchemaAttrNameMapID is the name of the schema attribute for the IOMMU group in a map of a hardware
	// mapping.
	hardwareMappingSchemaAttrNameMapIOMMUGroup = "iommu_group"

	// hardwareMappingSchemaAttrNameMapNode is the name of the schema attribute for the node in a map of a hardware
	// mapping.
	hardwareMappingSchemaAttrNameMapNode = "node"

	// hardwareMappingSchemaAttrNameMapPath is the name of the schema attribute for the path in a map of a hardware
	// mapping.
	hardwareMappingSchemaAttrNameMapPath = "path"

	// hardwareMappingSchemaAttrNameMapSubsystemID is the name of the schema attribute for the subsystem ID in a map of a
	// hardware mapping.
	hardwareMappingSchemaAttrNameMapSubsystemID = "subsystem_id"

	// hardwareMappingSchemaAttrNameMediatedDevices is the name of the schema attribute for the mediated devices in a map
	// of a hardware mapping.
	hardwareMappingSchemaAttrNameMediatedDevices = "mediated_devices"

	// hardwareMappingSchemaAttrNameName is the name of the schema attribute for the name of a hardware mapping.
	hardwareMappingSchemaAttrNameName = "name"

	// hardwareMappingSchemaAttrNameTerraformID is the name of the schema attribute for the Terraform ID of a hardware
	// mapping.
	hardwareMappingSchemaAttrNameTerraformID = "id"

	// hardwareMappingSchemaAttrNameType is the name of the schema attribute for the [proxmoxtypes.HardwareMappingType].
	hardwareMappingSchemaAttrNameType = "type"

	// hardwareMappingsSchemaAttrNameCheckNode is the name of the schema attribute for the "check node" option of a
	// hardware mappings data source.
	hardwareMappingsSchemaAttrNameCheckNode = "check_node"

	// hardwareMappingsSchemaAttrNameChecks is the name of the schema attribute for the node checks diagnostics of a
	// hardware mapping data source.
	// Note that the Proxmox VE API attribute for [proxmoxtypes.HardwareMappingTypeUSB] is named "errors", but we map it
	// as "checks" since this naming is
	// generally across the Proxmox VE web UI and API documentations, including the attribute for
	// [proxmoxtypes.HardwareMappingTypePCI].
	// This still follows the Terraform "best practices" as it improves the user experience by matching the field name to
	// the naming used in the human-facing
	// interfaces.
	hardwareMappingsSchemaAttrNameChecks = "checks"

	// hardwareMappingsSchemaAttrNameChecksDiagnosticsMappingID is the name of the schema attribute for a node check
	// diagnostic mapping ID of a hardware mappings
	// data source.
	hardwareMappingsSchemaAttrNameChecksDiagnosticsMappingID = "mapping_id"

	// hardwareMappingsSchemaAttrNameChecksDiagnosticsMessage is the name of the schema attribute for a node check
	// diagnostic message of a hardware mappings data
	// source.
	hardwareMappingsSchemaAttrNameChecksDiagnosticsMessage = "message"

	// hardwareMappingsSchemaAttrNameChecksDiagnosticsSeverity is the name of the schema attribute for a node check
	// diagnostic severity of a hardware mappings
	// data source.
	hardwareMappingsSchemaAttrNameChecksDiagnosticsSeverity = "severity"

	// hardwareMappingsSchemaAttrNameHardwareMappingIDs is the name of the schema attribute for the hardware mapping IDs
	// of a hardware mappings data source.
	hardwareMappingsSchemaAttrNameHardwareMappingIDs = "ids"
)

// hardwareMapPCIModel maps the schema data for the map of a PCI hardware mapping.
//

type hardwareMapPCIModel struct {
	// Comment is the "comment" for the map.
	// This field is optional and is omitted by the Proxmox API when not set.
	// Note that the Proxmox VE API attribute is named "description", but we map it as a comment since this naming is
	// generally across the Proxmox VE web UI and
	// API documentations. This still follows the Terraform "best practices" as it improves the user experience by
	// matching the field name to the naming used in
	// the human-facing interfaces.
	// References:
	//   1. https://developer.hashicorp.com/terraform/plugin/best-practices/hashicorp-provider-design-principles#resource-and-attribute-schema-should-closely-match-the-underlying-api
	Comment types.String `tfsdk:"comment"`

	// ID is the identifier of the map.
	ID types.String `tfsdk:"id"`

	// IOMMUGroup is the "IOMMU group" for the map.
	// This field is optional and is omitted by the Proxmox API when not set.
	IOMMUGroup types.Int64 `tfsdk:"iommu_group"`

	// Node is the "node name" for the map.
	Node types.String `tfsdk:"node"`

	// Path is the "path" for the map.
	Path customtypes.HardwareMappingPathValue `tfsdk:"path"`

	// SubsystemID is the "subsystem ID" for the map.
	// This field is not mandatory for the Proxmox API call, but causes a PCI hardware mapping to be incomplete when not
	// set.
	SubsystemID types.String `tfsdk:"subsystem_id"`
}

// hardwareMapUSBModel maps the schema data for the map of a USB hardware mapping.
type hardwareMapUSBModel struct {
	// Comment is the "comment" for the map.
	// This field is optional and is omitted by the Proxmox API when not set.
	// Note that the Proxmox VE API attribute is named "description", but we map it as a comment since this naming is
	// generally across the Proxmox VE web UI and
	// API documentations. This still follows the Terraform "best practices" as it improves the user experience by
	// matching the field name to the naming used in
	// the human-facing interfaces.
	// References:
	//   1. https://developer.hashicorp.com/terraform/plugin/best-practices/hashicorp-provider-design-principles#resource-and-attribute-schema-should-closely-match-the-underlying-api
	Comment types.String `tfsdk:"comment"`

	// ID is the identifier of the map.
	ID types.String `tfsdk:"id"`

	// Node is the "node name" for the map.
	Node types.String `tfsdk:"node"`

	// Path is the "path" for the map.
	Path customtypes.HardwareMappingPathValue `tfsdk:"path"`
}

// hardwareMappingModel maps the schema data for a PCI hardware mapping.
type hardwareMappingPCIModel struct {
	// Comment is the comment of the PCI hardware mapping.
	// Note that the Proxmox VE API attribute is named "description", but we map it as a comment since this naming is
	// generally across the Proxmox VE web UI and
	// API documentations. This still follows the Terraform "best practices" as it improves the user experience by
	// matching the field name to the naming used in
	// the human-facing interfaces.
	// References:
	//   1. https://developer.hashicorp.com/terraform/plugin/best-practices/hashicorp-provider-design-principles#resource-and-attribute-schema-should-closely-match-the-underlying-api
	Comment types.String `tfsdk:"comment"`

	// ID is the Terraform identifier.
	ID types.String `tfsdk:"id"`

	// Name is the name of the PCI hardware mapping.
	Name types.String `tfsdk:"name"`

	// Map is the map of the PCI hardware mapping.
	Map []hardwareMapPCIModel `tfsdk:"map"`

	// MediatedDevices is the indicator for mediated devices of the PCI hardware mapping.
	MediatedDevices types.Bool `tfsdk:"mediated_devices"`
}

// hardwareMappingUSBModel maps the schema data for a USB hardware mapping.
type hardwareMappingUSBModel struct {
	// Comment is the comment of the USB hardware mapping.
	// Note that the Proxmox VE API attribute is named "description", but we map it as a comment since this naming is
	// generally across the Proxmox VE web UI and
	// API documentations. This still follows the Terraform "best practices" as it improves the user experience by
	// matching the field name to the naming used in
	// the human-facing interfaces.
	// References:
	//   1. https://developer.hashicorp.com/terraform/plugin/best-practices/hashicorp-provider-design-principles#resource-and-attribute-schema-should-closely-match-the-underlying-api
	Comment types.String `tfsdk:"comment"`

	// ID is the Terraform identifier.
	ID types.String `tfsdk:"id"`

	// Name is the name of the USB hardware mapping.
	Name types.String `tfsdk:"name"`

	// Map is the map of the USB hardware mapping.
	Map []hardwareMapUSBModel `tfsdk:"map"`
}

// hardwareMappingsModel maps the schema data for a hardware mappings data source.
type hardwareMappingsModel struct {
	// Checks might contain relevant hardware mapping diagnostics about incorrect configurations for the node name set
	// defined by CheckNode.
	// Note that the Proxmox VE API attribute for [proxmoxtypes.HardwareMappingTypeUSB] is named "errors", but we map it
	// as "checks" since this naming is
	// generally across the Proxmox VE web UI and API documentations, including the attribute for
	// [proxmoxtypes.HardwareMappingTypePCI].
	// Also note that the Proxmox VE API, for whatever reason, only returns one error at a time, even though the field is
	// an array.
	// This still follows the Terraform "best practices" as it improves the user experience by matching the field name to
	// the naming used in the human-facing interfaces.
	Checks []hardwareMappingsNodeCheckDiagnosticModel `tfsdk:"checks"`

	// CheckNode is the name of the node whose configuration should be checked for correctness.
	CheckNode types.String `tfsdk:"check_node"`

	// ID is the Terraform identifier.
	ID types.String `tfsdk:"id"`

	// MappingIDs is the set of hardware mapping identifiers.
	MappingIDs types.Set `tfsdk:"ids"`

	// Type is the [proxmoxtypes.HardwareMappingType].
	Type types.String `tfsdk:"type"`
}

// hardwareMappingsNodeCheckDiagnosticModel maps the schema data for hardware mapping node check diagnostic data.
type hardwareMappingsNodeCheckDiagnosticModel struct {
	// MappingID is the corresponding hardware mapping ID of this node check diagnostic entry.
	MappingID types.String `tfsdk:"mapping_id"`

	// Message is the message of the node check diagnostic entry.
	Message types.String `tfsdk:"message"`

	// Severity is the severity of the node check diagnostic entry.
	Severity types.String `tfsdk:"severity"`
}

// importFromAPI imports the contents of a PCI hardware mapping model from the Proxmox VE API's response data.
func (hm *hardwareMappingPCIModel) importFromAPI(
	_ context.Context,
	data *apitypes.HardwareMappingGetResponseData,
) {
	// Ensure that both the ID and name are in sync.
	hm.Name = hm.ID
	// The attribute is named "description" by the Proxmox VE API, but we map it as a comment since this naming is
	// generally across the Proxmox VE web UI and API
	// documentations.
	hm.Comment = types.StringPointerValue(data.Description)
	maps := make([]hardwareMapPCIModel, len(data.Map))

	for idx, pveMap := range data.Map {
		tfMap := hardwareMapPCIModel{
			ID:   pveMap.ID.ToValue(),
			Node: types.StringValue(pveMap.Node),
			Path: customtypes.NewHardwareMappingPathPointerValue(pveMap.Path),
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
func (hm *hardwareMappingPCIModel) toCreateRequest() *apitypes.HardwareMappingCreateRequestBody {
	return &apitypes.HardwareMappingCreateRequestBody{
		HardwareMappingDataBase: hm.toRequestBase(),
		ID:                      hm.ID.ValueString(),
	}
}

// toRequestBase builds the common request data structure for the PCI hardware mapping creation or update API calls.
func (hm *hardwareMappingPCIModel) toRequestBase() apitypes.HardwareMappingDataBase {
	dataBase := apitypes.HardwareMappingDataBase{
		// The attribute is named "description" by the Proxmox VE API, but we map it as a comment since this naming is
		// generally across the Proxmox VE web UI and
		// API documentations.
		Description: hm.Comment.ValueStringPointer(),
	}
	maps := make([]proxmoxtypes.HardwareMapping, len(hm.Map))

	for idx, tfMap := range hm.Map {
		pveMap := proxmoxtypes.HardwareMapping{
			// The attribute is named "description" by the Proxmox VE API, but we map it as a comment since this naming is
			// generally across the Proxmox VE web UI and
			// API documentations.
			Description: tfMap.Comment.ValueStringPointer(),
			ID:          proxmoxtypes.HardwareMappingDeviceID(tfMap.ID.ValueString()),
			IOMMUGroup:  tfMap.IOMMUGroup.ValueInt64Pointer(),
			Node:        tfMap.Node.ValueString(),
			Path:        tfMap.Path.ValueStringPointer(),
			SubsystemID: proxmoxtypes.HardwareMappingDeviceID(tfMap.SubsystemID.ValueString()),
		}
		maps[idx] = pveMap
	}

	dataBase.Map = maps
	dataBase.MediatedDevices.FromValue(hm.MediatedDevices)

	return dataBase
}

// toUpdateRequest builds the request data structure for updating an existing PCI hardware mapping.
func (hm *hardwareMappingPCIModel) toUpdateRequest(
	currentState *hardwareMappingPCIModel,
) *apitypes.HardwareMappingUpdateRequestBody {
	var del []string

	baseRequest := hm.toRequestBase()

	if hm.Comment.IsNull() && !currentState.Comment.IsNull() {
		// hardwareMappingSchemaAttrNameComment is the name of the schema attribute for the comment of a PCI hardware
		// mapping.
		// The Proxmox VE API attribute is named "description" while we name it "comment" internally since this naming is
		// generally used across the Proxmox VE web
		// UI and API documentations. This still follows the Terraform "best practices" as it improves the user experience
		// by matching the field name to the naming
		// used in the human-facing interfaces.
		del = append(del, proxmoxtypes.HardwareMappingAttrNameDescription)
	}

	if hm.MediatedDevices.IsNull() || !hm.MediatedDevices.ValueBool() {
		del = append(del, apitypes.HardwareMappingPCIMediatedDevicesAPIParamName)

		baseRequest.MediatedDevices.FromValue(types.BoolValue(false))
	}

	return &apitypes.HardwareMappingUpdateRequestBody{
		HardwareMappingDataBase: baseRequest,
		Delete:                  del,
	}
}

// importFromAPI imports the contents of a USB hardware mapping model from the Proxmox VE API's response data.
func (hm *hardwareMappingUSBModel) importFromAPI(
	_ context.Context,
	data *apitypes.HardwareMappingGetResponseData,
) {
	// Ensure that both the ID and name are in sync.
	hm.Name = hm.ID
	// The attribute is named "description" by the Proxmox VE API, but we map it as a comment since this naming is
	// generally across the Proxmox VE web UI and API
	// documentations.
	hm.Comment = types.StringPointerValue(data.Description)
	maps := make([]hardwareMapUSBModel, len(data.Map))

	for idx, pveMap := range data.Map {
		tfMap := hardwareMapUSBModel{
			ID:   pveMap.ID.ToValue(),
			Node: types.StringValue(pveMap.Node),
			Path: customtypes.NewHardwareMappingPathPointerValue(pveMap.Path),
		}
		if pveMap.Description != nil {
			tfMap.Comment = types.StringPointerValue(pveMap.Description)
		}

		maps[idx] = tfMap
	}

	hm.Map = maps
}

// toCreateRequest builds the request data structure for creating a new USB hardware mapping.
func (hm *hardwareMappingUSBModel) toCreateRequest() *apitypes.HardwareMappingCreateRequestBody {
	return &apitypes.HardwareMappingCreateRequestBody{
		HardwareMappingDataBase: hm.toRequestBase(),
		ID:                      hm.ID.ValueString(),
	}
}

// toRequestBase builds the common request data structure for the USB hardware mapping creation or update API calls.
func (hm *hardwareMappingUSBModel) toRequestBase() apitypes.HardwareMappingDataBase {
	dataBase := apitypes.HardwareMappingDataBase{
		// The attribute is named "description" by the Proxmox VE API, but we map it as a comment since this naming is
		// generally across the Proxmox VE web UI and
		// API documentations.
		Description: hm.Comment.ValueStringPointer(),
	}
	maps := make([]proxmoxtypes.HardwareMapping, len(hm.Map))

	for idx, tfMap := range hm.Map {
		pveMap := proxmoxtypes.HardwareMapping{
			ID:   proxmoxtypes.HardwareMappingDeviceID(tfMap.ID.ValueString()),
			Node: tfMap.Node.ValueString(),
			Path: tfMap.Path.ValueStringPointer(),
		}
		if !tfMap.Comment.IsNull() {
			// The attribute is named "description" by the Proxmox VE API, but we map it as a comment since this naming is
			// generally across the Proxmox VE web UI and
			// API documentations.
			pveMap.Description = tfMap.Comment.ValueStringPointer()
		}

		maps[idx] = pveMap
	}

	dataBase.Map = maps

	return dataBase
}

// toUpdateRequest builds the request data structure for updating an existing USB hardware mapping.
func (hm *hardwareMappingUSBModel) toUpdateRequest(
	currentState *hardwareMappingUSBModel,
) *apitypes.HardwareMappingUpdateRequestBody {
	var del []string

	if hm.Comment.IsNull() && !currentState.Comment.IsNull() {
		// hardwareMappingSchemaAttrNameComment is the name of the schema attribute for the comment of a USB hardware
		// mapping.
		// The Proxmox VE API attribute is named "description" while we name it "comment" internally since this naming is
		// generally used across the Proxmox VE web
		// UI and API documentations. This still follows the Terraform "best practices" as it improves the user experience
		// by matching the field name to the naming
		// used in the human-facing interfaces.
		del = append(del, proxmoxtypes.HardwareMappingAttrNameDescription)
	}

	return &apitypes.HardwareMappingUpdateRequestBody{
		HardwareMappingDataBase: hm.toRequestBase(),
		Delete:                  del,
	}
}
