/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package storage

import (
	"fmt"
	"math"

	"github.com/bpg/terraform-provider-proxmox/proxmox/storage"
	proxmox_types "github.com/bpg/terraform-provider-proxmox/proxmox/types"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// BackupModel maps the backup block schema.
type BackupModel struct {
	MaxProtectedBackups types.Int64 `tfsdk:"max_protected_backups"`
	KeepAll             types.Bool  `tfsdk:"keep_all"`
	KeepLast            types.Int64 `tfsdk:"keep_last"`
	KeepHourly          types.Int64 `tfsdk:"keep_hourly"`
	KeepDaily           types.Int64 `tfsdk:"keep_daily"`
	KeepWeekly          types.Int64 `tfsdk:"keep_weekly"`
	KeepMonthly         types.Int64 `tfsdk:"keep_monthly"`
	KeepYearly          types.Int64 `tfsdk:"keep_yearly"`
}

func (m BackupModel) toAPI() (storage.DataStoreWithBackups, error) {
	var backups storage.DataStoreWithBackups

	intPtrFromInt64 := func(v int64) (*int, error) {
		if v > math.MaxInt || v < math.MinInt {
			return nil, fmt.Errorf("value out of range: %d", v)
		}

		i := int(v)

		return &i, nil
	}

	if !m.MaxProtectedBackups.IsNull() && !m.MaxProtectedBackups.IsUnknown() {
		v := proxmox_types.CustomInt64(m.MaxProtectedBackups.ValueInt64())
		backups.MaxProtectedBackups = &v
	}

	if !m.KeepAll.IsNull() && !m.KeepAll.IsUnknown() && m.KeepAll.ValueBool() {
		v := proxmox_types.CustomBool(true)
		backups.KeepAll = &v
	}

	setKeepCount := func(tf types.Int64, target **int) error {
		if tf.IsNull() || tf.IsUnknown() {
			return nil
		}

		ptr, err := intPtrFromInt64(tf.ValueInt64())
		if err != nil {
			return err
		}

		*target = ptr

		return nil
	}

	if err := setKeepCount(m.KeepLast, &backups.KeepLast); err != nil {
		return storage.DataStoreWithBackups{}, err
	}

	if err := setKeepCount(m.KeepHourly, &backups.KeepHourly); err != nil {
		return storage.DataStoreWithBackups{}, err
	}

	if err := setKeepCount(m.KeepDaily, &backups.KeepDaily); err != nil {
		return storage.DataStoreWithBackups{}, err
	}

	if err := setKeepCount(m.KeepWeekly, &backups.KeepWeekly); err != nil {
		return storage.DataStoreWithBackups{}, err
	}

	if err := setKeepCount(m.KeepMonthly, &backups.KeepMonthly); err != nil {
		return storage.DataStoreWithBackups{}, err
	}

	if err := setKeepCount(m.KeepYearly, &backups.KeepYearly); err != nil {
		return storage.DataStoreWithBackups{}, err
	}

	if backups.KeepAll != nil && (backups.KeepLast != nil ||
		backups.KeepHourly != nil ||
		backups.KeepDaily != nil ||
		backups.KeepWeekly != nil ||
		backups.KeepMonthly != nil ||
		backups.KeepYearly != nil) {
		return storage.DataStoreWithBackups{}, fmt.Errorf("keep_all conflicts with other keep_* settings")
	}

	return backups, nil
}
