/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package storage

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework-validators/setvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ resource.Resource                = &rbdStorageResource{}
	_ resource.ResourceWithConfigure   = &rbdStorageResource{}
	_ resource.ResourceWithImportState = &rbdStorageResource{}
)

// NewRBDStorageResource creates the proxmox_storage_rbd resource.
func NewRBDStorageResource() resource.Resource {
	return &rbdStorageResource{
		storageResource: &storageResource[*RBDStorageModel, RBDStorageModel]{
			storageType:  "rbd",
			resourceName: "proxmox_storage_rbd",
		},
	}
}

// rbdStorageResource is the resource implementation.
type rbdStorageResource struct {
	*storageResource[*RBDStorageModel, RBDStorageModel]
}

// Metadata returns the resource type name.
func (r *rbdStorageResource) Metadata(_ context.Context, _ resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = r.resourceName
}

// Schema defines the schema for the RBD storage resource.
func (r *rbdStorageResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	attributes := map[string]schema.Attribute{
		"content": schema.SetAttribute{
			Description: "The content types that can be stored on this storage.",
			MarkdownDescription: "The content types that can be stored on this storage. " +
				"Valid values: `images` (VM disk images), `rootdir` (container root directories). " +
				"Proxmox defaults to `images`.",
			ElementType: types.StringType,
			Optional:    true,
			Computed:    true,
			Validators: []validator.Set{
				setvalidator.ValueStringsAre(
					stringvalidator.OneOf("images", "rootdir"),
				),
			},
		},
		"pool": schema.StringAttribute{
			Description: "The name of the Ceph pool that holds the storage's images.",
			MarkdownDescription: "The name of the Ceph pool that holds the storage's images. " +
				"For an erasure-coded `proxmox_ceph_pool`, this is its `metadata_pool_name`, and `data_pool` is its " +
				"`data_pool.name`.",
			Required: true,
			Validators: []validator.String{
				stringvalidator.LengthAtLeast(1),
			},
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.RequiresReplace(),
			},
		},
		"data_pool": schema.StringAttribute{
			Description: "The name of the erasure-coded Ceph pool that holds the image data. " +
				"The image metadata stays in `pool`.",
			Optional: true,
			Validators: []validator.String{
				stringvalidator.LengthAtLeast(1),
			},
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.RequiresReplace(),
			},
		},
		"namespace": schema.StringAttribute{
			Description: "The RBD namespace within `pool` for the storage's images. " +
				"Proxmox doesn't create the namespace.",
			Optional: true,
			Validators: []validator.String{
				stringvalidator.LengthAtLeast(1),
			},
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.RequiresReplace(),
			},
		},
		"krbd": schema.BoolAttribute{
			Description: "Whether to always access the images through the kernel RBD module instead of librbd. " +
				"Containers always use the kernel module.",
			Optional: true,
			Computed: true,
			Default:  booldefault.StaticBool(false),
		},
		"shared": schema.BoolAttribute{
			Description: "Whether the storage is shared across all nodes. Proxmox always shares RBD storage.",
			Computed:    true,
		},
	}
	factory := newStorageSchemaFactory()
	factory.WithAttributes(attributes)
	factory.WithDescription("Manages an RBD storage on the Proxmox VE cluster's own Ceph cluster.")
	resp.Schema = *factory.Schema
	resp.Schema.MarkdownDescription = "Manages an RBD storage on the Proxmox VE cluster's own Ceph cluster, " +
		"for example on a pool managed by `proxmox_ceph_pool`. Storages on external Ceph clusters aren't supported yet."
}
