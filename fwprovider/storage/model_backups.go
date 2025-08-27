/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package storage

import (
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
