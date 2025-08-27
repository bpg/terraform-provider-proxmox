/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package storage

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/bpg/terraform-provider-proxmox/proxmox/storage"
	"github.com/hashicorp/terraform-plugin-framework-validators/boolvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure the implementation satisfies the expected interfaces.
var _ resource.Resource = &pbsStorageResource{}

// NewProxmoxBackupServerStorageResource is a helper function to simplify the provider implementation.
func NewProxmoxBackupServerStorageResource() resource.Resource {
	return &pbsStorageResource{
		storageResource: &storageResource[
			*PBSStorageModel,
			PBSStorageModel,
		]{
			storageType:  "pbs",
			resourceName: "proxmox_virtual_environment_storage_pbs",
		},
	}
}

// pbsStorageResource is the resource implementation.
type pbsStorageResource struct {
	*storageResource[*PBSStorageModel, PBSStorageModel]
}

// Metadata returns the resource type name.
func (r *pbsStorageResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = r.resourceName
}

// Create is the generic create function.
func (r *pbsStorageResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan PBSStorageModel

	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	requestBody, err := plan.toCreateAPIRequest(ctx)
	if err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("Error creating API request for %s storage", r.storageType), err.Error())
		return
	}

	responseData, err := r.client.Storage().CreateDatastore(ctx, requestBody)
	if err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("Error creating %s storage", r.storageType), err.Error())
		return
	}

	plan.Shared = types.BoolValue(false)

	if !plan.GenerateEncryptionKey.IsNull() && plan.GenerateEncryptionKey.ValueBool() {
		var encryptionKey storage.EncryptionKey

		err := json.Unmarshal([]byte(*responseData.Config.EncryptionKey), &encryptionKey)
		if err != nil {
			resp.Diagnostics.AddError(fmt.Sprintf("Error unmarshaling encryption key for %s storage", r.storageType), err.Error())
			return
		}

		plan.GeneratedEncryptionKey = types.StringValue(*responseData.Config.EncryptionKey)
		plan.EncryptionKeyFingerprint = types.StringValue(encryptionKey.Fingerprint)
	} else {
		plan.GeneratedEncryptionKey = types.StringNull()
	}

	if !plan.EncryptionKey.IsNull() {
		var encryptionKey storage.EncryptionKey

		err := json.Unmarshal([]byte(*responseData.Config.EncryptionKey), &encryptionKey)
		if err != nil {
			resp.Diagnostics.AddError(fmt.Sprintf("Error unmarshaling encryption key for %s storage", r.storageType), err.Error())
			return
		}

		plan.EncryptionKey = types.StringValue(*responseData.Config.EncryptionKey)
		plan.EncryptionKeyFingerprint = types.StringValue(encryptionKey.Fingerprint)
	}

	diags = resp.State.Set(ctx, &plan)
	resp.Diagnostics.Append(diags...)
}

// Schema defines the schema for the Proxmox Backup Server storage resource.
func (r *pbsStorageResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	attributes := map[string]schema.Attribute{
		"server": schema.StringAttribute{
			Description: "The IP address or DNS name of the Proxmox Backup Server.",
			Required:    true,
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.RequiresReplace(),
			},
		},
		"datastore": schema.StringAttribute{
			Description: "The name of the datastore on the Proxmox Backup Server.",
			Required:    true,
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.RequiresReplace(),
			},
		},
		"username": schema.StringAttribute{
			Description: "The username for authenticating with the Proxmox Backup Server.",
			Required:    true,
		},
		"password": schema.StringAttribute{
			Description: "The password for authenticating with the Proxmox Backup Server.",
			Required:    true,
			Sensitive:   true,
		},
		"namespace": schema.StringAttribute{
			Description: "The namespace to use on the Proxmox Backup Server.",
			Optional:    true,
		},
		"fingerprint": schema.StringAttribute{
			Description: "The SHA256 fingerprint of the Proxmox Backup Server's certificate.",
			Optional:    true,
		},
		"encryption_key": schema.StringAttribute{
			Description: "An existing encryption key for the datastore. This is a sensitive value. Conflicts with `generate_encryption_key`.",
			Optional:    true,
			Sensitive:   true,
			Validators: []validator.String{
				stringvalidator.ConflictsWith(path.MatchRoot("generate_encryption_key")),
			},
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.RequiresReplace(),
			},
		},
		"encryption_key_fingerprint": schema.StringAttribute{
			Description: "The SHA256 fingerprint of the encryption key currently in use.",
			Computed:    true,
		},
		"generate_encryption_key": schema.BoolAttribute{
			Description: "If set to true, Proxmox will generate a new encryption key. The key will be stored in the `generated_encryption_key` attribute. " +
				"Conflicts with `encryption_key`.",
			Optional: true,
			Validators: []validator.Bool{
				boolvalidator.ConflictsWith(path.MatchRoot("encryption_key")),
			},
			PlanModifiers: []planmodifier.Bool{
				boolplanmodifier.RequiresReplace(),
			},
		},
		"generated_encryption_key": schema.StringAttribute{
			Description: "The encryption key returned by Proxmox when `generate_encryption_key` is true.",
			Computed:    true,
			Sensitive:   true,
		},
		"shared": schema.BoolAttribute{
			Description: "Whether the storage is shared across all nodes.",
			Computed:    true,
		},
	}
	factory := NewStorageSchemaFactory()
	factory.WithAttributes(attributes)
	factory.WithDescription("Manages a Proxmox Backup Server (PBS) storage in Proxmox VE.")
	resp.Schema = *factory.Schema
}
