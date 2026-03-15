/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package files

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Model is the data model for the files data source.
type Model struct {
	NodeName    types.String `tfsdk:"node_name"`
	DatastoreID types.String `tfsdk:"datastore_id"`
	ContentType types.String `tfsdk:"content_type"`
	Files       []File       `tfsdk:"files"`
}

// File represents a single file in the data source output.
type File struct {
	ID          types.String `tfsdk:"id"`
	ContentType types.String `tfsdk:"content_type"`
	FileName    types.String `tfsdk:"file_name"`
	FileFormat  types.String `tfsdk:"file_format"`
	FileSize    types.Int64  `tfsdk:"file_size"`
	VMID        types.Int64  `tfsdk:"vmid"`
}
