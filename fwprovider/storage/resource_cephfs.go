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
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ resource.Resource                = &cephFSStorageResource{}
	_ resource.ResourceWithConfigure   = &cephFSStorageResource{}
	_ resource.ResourceWithImportState = &cephFSStorageResource{}
)

// NewCephFSStorageResource creates the proxmox_storage_cephfs resource.
func NewCephFSStorageResource() resource.Resource {
	return &cephFSStorageResource{
		storageResource: &storageResource[*CephFSStorageModel, CephFSStorageModel]{
			storageType:  "cephfs",
			resourceName: "proxmox_storage_cephfs",
		},
	}
}

// cephFSStorageResource is the resource implementation.
type cephFSStorageResource struct {
	*storageResource[*CephFSStorageModel, CephFSStorageModel]
}

// Metadata returns the resource type name.
func (r *cephFSStorageResource) Metadata(_ context.Context, _ resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = r.resourceName
}

// Schema defines the schema for the CephFS storage resource.
func (r *cephFSStorageResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	attributes := map[string]schema.Attribute{
		"content": schema.SetAttribute{
			Description: "The content types that can be stored on this storage.",
			MarkdownDescription: "The content types that can be stored on this storage. " +
				"Valid values: `backup` (VM backups), `import` (VM disk images for import), `iso` (ISO images), " +
				"`snippets` (cloud-init, hook scripts, etc.), `vztmpl` (container templates). " +
				"Proxmox defaults to `backup`.",
			ElementType: types.StringType,
			Optional:    true,
			Computed:    true,
			Validators: []validator.Set{
				setvalidator.ValueStringsAre(
					stringvalidator.OneOf("backup", "import", "iso", "snippets", "vztmpl"),
				),
			},
		},
		"fs_name": schema.StringAttribute{
			Description: "The name of the CephFS to mount. Proxmox uses the cluster's default filesystem when unset.",
			Optional:    true,
			Validators: []validator.String{
				stringvalidator.LengthAtLeast(1),
			},
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.RequiresReplace(),
			},
		},
		"subdir": schema.StringAttribute{
			Description: "The directory within the filesystem to mount. Proxmox mounts the filesystem's root when unset.",
			Optional:    true,
			Validators: []validator.String{
				stringvalidator.LengthAtLeast(1),
			},
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.RequiresReplace(),
			},
		},
		"fuse": schema.BoolAttribute{
			Description: "Whether to mount the filesystem through FUSE (`ceph-fuse`) instead of the kernel client.",
			Optional:    true,
			Computed:    true,
			Default:     booldefault.StaticBool(false),
			PlanModifiers: []planmodifier.Bool{
				boolplanmodifier.RequiresReplace(),
			},
		},
		"options": schema.StringAttribute{
			Description: "The mount options (see `man mount.ceph`). Proxmox applies them when it next mounts the storage.",
			Optional:    true,
		},
		"shared": schema.BoolAttribute{
			Description: "Whether the storage is shared across all nodes. Proxmox always shares CephFS storage.",
			Computed:    true,
		},
	}

	factory := newStorageSchemaFactory()
	factory.WithAttributes(attributes)
	factory.WithDirCreationOptions()
	factory.WithDescription("Manages a CephFS storage on the Proxmox VE cluster's own Ceph cluster.")
	factory.WithBackupBlock()
	resp.Schema = *factory.Schema
	resp.Schema.MarkdownDescription = "Manages a CephFS storage on the Proxmox VE cluster's own Ceph cluster. " +
		"Proxmox mounts an enabled storage at `/mnt/pve/<id>` on each node it's available on. " +
		"Storages on external Ceph clusters aren't supported yet."
}
