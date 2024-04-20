package vm

import (
	"context"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"

	proxmoxtypes "github.com/bpg/terraform-provider-proxmox/proxmox/types"
)

type vmModel struct {
	Description types.String   `tfsdk:"description"`
	ID          types.Int64    `tfsdk:"id"`
	Name        types.String   `tfsdk:"name"`
	NodeName    types.String   `tfsdk:"node_name"`
	Tags        types.Set      `tfsdk:"tags"`
	Timeouts    timeouts.Value `tfsdk:"timeouts"`
}

func (m *vmModel) tagsString(ctx context.Context, diags diag.Diagnostics) *string {
	if m.Tags.IsNull() {
		return nil
	}

	elems := make([]types.String, 0, len(m.Tags.Elements()))
	d := m.Tags.ElementsAs(ctx, &elems, false)
	diags.Append(d...)

	if d.HasError() {
		return nil
	}

	var sanitizedTags []string
	for _, el := range elems {
		if el.IsNull() || el.IsUnknown() {
			continue
		}
		sanitizedTag := strings.TrimSpace(el.ValueString())
		if len(sanitizedTag) > 0 {
			sanitizedTags = append(sanitizedTags, sanitizedTag)
		}
	}

	return proxmoxtypes.StrPtr(strings.Join(sanitizedTags, ";"))
}
