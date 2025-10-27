package pools

import (
	"errors"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"strconv"
	"strings"
)

type poolMembershipModel struct {
	ID     types.String `tfsdk:"id"`
	VmID   types.Int64  `tfsdk:"vm_id"`
	PoolId types.String `tfsdk:"pool_id"`
}

const poolMembershipIDFormat = "{pool_id}/{vm_id}"

func (p poolMembershipModel) generateID() types.String {
	return types.StringValue(fmt.Sprintf("%s/%d", p.PoolId.ValueString(), p.VmID.ValueInt64()))
}

func parseMembershipModelFromID(id string) (*poolMembershipModel, error) {
	parts := strings.Split(id, "/")

	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid pool membership ID format %#v, expected: %s", id, poolMembershipIDFormat)
	}

	poolId := parts[0]
	vmId, err := strconv.ParseInt(parts[1], 10, 64)
	if err != nil {
		vmIdErr := fmt.Errorf("wrong vm_id format: %s", parts[1])
		return nil, errors.Join(vmIdErr, err)
	}

	model := poolMembershipModel{
		ID:     types.StringValue(id),
		VmID:   types.Int64Value(vmId),
		PoolId: types.StringValue(poolId),
	}

	return &model, nil
}
