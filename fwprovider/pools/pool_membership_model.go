/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package pools

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/types"
)

type MembershipType = string

var ErrInvalidMembershipType = errors.New("invalid pool membership type")

const (
	MembershipTypeVm      MembershipType = "vm"
	MembershipTypeStorage MembershipType = "storage"
)

func NewMembershipType(raw string) (MembershipType, error) {
	switch raw {
	case MembershipTypeVm,
		MembershipTypeStorage:
		return raw, nil
	default:
		return "", ErrInvalidMembershipType
	}
}

type poolMembershipModel struct {
	ID        types.String `tfsdk:"id"`
	VmID      types.Int64  `tfsdk:"vm_id"`
	StorageID types.String `tfsdk:"storage_id"`
	PoolID    types.String `tfsdk:"pool_id"`
	Type      types.String `tfsdk:"type"`
}

const poolMembershipIDFormat = "{pool_id}/{type}/{member_id}"

// Proxmox API for managing resource pools does not differentiate lxc containers and vms. All of them are considered VMs.
func (p poolMembershipModel) deduceMembershipType() (MembershipType, error) {
	var membershipType MembershipType

	switch {
	case !p.VmID.IsNull():
		membershipType = MembershipTypeVm
	case !p.StorageID.IsNull():
		membershipType = MembershipTypeStorage
	default:
		return "", ErrInvalidMembershipType
	}

	return membershipType, nil
}

func (p poolMembershipModel) generateID() (types.String, error) {
	var memberId string

	switch p.Type.ValueString() {
	case MembershipTypeVm:
		memberId = strconv.FormatInt(p.VmID.ValueInt64(), 10)
	case MembershipTypeStorage:
		memberId = p.StorageID.ValueString()
	default:
		return types.String{}, ErrInvalidMembershipType
	}

	return types.StringValue(fmt.Sprintf("%s/%s/%s", p.PoolID.ValueString(), p.Type.ValueString(), memberId)), nil
}

func createMembershipModelFromID(id string) (*poolMembershipModel, error) {
	poolId, membershipRawType, memberRawId, idParseErr := parseMembershipResourceID(id)

	if idParseErr != nil {
		return nil, idParseErr
	}

	membershipType, err := NewMembershipType(membershipRawType)
	if err != nil {
		return nil, err
	}

	model := poolMembershipModel{
		ID:     types.StringValue(id),
		PoolID: types.StringValue(poolId),
		Type:   types.StringValue(membershipType),
	}

	switch membershipType {
	case MembershipTypeVm:
		vmId, err := strconv.ParseInt(memberRawId, 10, 64)
		if err != nil {
			vmIdErr := fmt.Errorf("wrong vm_id format: %s", memberRawId)
			return nil, errors.Join(vmIdErr, err)
		}

		model.VmID = types.Int64Value(vmId)

	case MembershipTypeStorage:
		model.StorageID = types.StringValue(memberRawId)
	}

	return &model, nil
}

func parseMembershipResourceID(id string) (string, string, string, error) {
	parts := strings.Split(id, "/")

	if len(parts) != 3 {
		return "", "", "", fmt.Errorf("invalid pool membership ID format %#v, expected: %s", id, poolMembershipIDFormat)
	}

	return parts[0], parts[1], parts[2], nil
}
