package storage

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// BackupModel maps the backup block schema.
type BackupModel struct {
	MaxProtectedBackups types.Int64 `tfsdk:"max_protected_backups"`
	KeepLast            types.Int64 `tfsdk:"keep_last"`
	KeepHourly          types.Int64 `tfsdk:"keep_hourly"`
	KeepDaily           types.Int64 `tfsdk:"keep_daily"`
	KeepWeekly          types.Int64 `tfsdk:"keep_weekly"`
	KeepMonthly         types.Int64 `tfsdk:"keep_monthly"`
	KeepYearly          types.Int64 `tfsdk:"keep_yearly"`
}
