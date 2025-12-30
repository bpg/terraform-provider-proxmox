package storage

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/require"
)

func TestBackupsKeepAllExcludesOtherKeepSettingsValidator(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	validatorUnderTest := backupsKeepAllExcludesOtherKeepSettingsValidator{}

	attributeTypes := map[string]attr.Type{
		"keep_all":     types.BoolType,
		"keep_last":    types.Int64Type,
		"keep_hourly":  types.Int64Type,
		"keep_daily":   types.Int64Type,
		"keep_weekly":  types.Int64Type,
		"keep_monthly": types.Int64Type,
		"keep_yearly":  types.Int64Type,
	}

	t.Run("errors when keep_all is true and another keep_* is set", func(t *testing.T) {
		t.Parallel()

		obj, diags := types.ObjectValue(attributeTypes, map[string]attr.Value{
			"keep_all":     types.BoolValue(true),
			"keep_last":    types.Int64Null(),
			"keep_hourly":  types.Int64Null(),
			"keep_daily":   types.Int64Value(7),
			"keep_weekly":  types.Int64Null(),
			"keep_monthly": types.Int64Null(),
			"keep_yearly":  types.Int64Null(),
		})
		require.False(t, diags.HasError())

		req := validator.ObjectRequest{
			Path:        path.Root("backups"),
			ConfigValue: obj,
		}
		resp := &validator.ObjectResponse{}

		validatorUnderTest.ValidateObject(ctx, req, resp)
		require.True(t, resp.Diagnostics.HasError())
	})

	t.Run("ok when keep_all is true and no other keep_* is set", func(t *testing.T) {
		t.Parallel()

		obj, diags := types.ObjectValue(attributeTypes, map[string]attr.Value{
			"keep_all":     types.BoolValue(true),
			"keep_last":    types.Int64Null(),
			"keep_hourly":  types.Int64Null(),
			"keep_daily":   types.Int64Null(),
			"keep_weekly":  types.Int64Null(),
			"keep_monthly": types.Int64Null(),
			"keep_yearly":  types.Int64Null(),
		})
		require.False(t, diags.HasError())

		req := validator.ObjectRequest{
			Path:        path.Root("backups"),
			ConfigValue: obj,
		}
		resp := &validator.ObjectResponse{}

		validatorUnderTest.ValidateObject(ctx, req, resp)
		require.False(t, resp.Diagnostics.HasError())
	})

	t.Run("ok when keep_all is false and other keep_* is set", func(t *testing.T) {
		t.Parallel()

		obj, diags := types.ObjectValue(attributeTypes, map[string]attr.Value{
			"keep_all":     types.BoolValue(false),
			"keep_last":    types.Int64Null(),
			"keep_hourly":  types.Int64Null(),
			"keep_daily":   types.Int64Value(7),
			"keep_weekly":  types.Int64Null(),
			"keep_monthly": types.Int64Null(),
			"keep_yearly":  types.Int64Null(),
		})
		require.False(t, diags.HasError())

		req := validator.ObjectRequest{
			Path:        path.Root("backups"),
			ConfigValue: obj,
		}
		resp := &validator.ObjectResponse{}

		validatorUnderTest.ValidateObject(ctx, req, resp)
		require.False(t, resp.Diagnostics.HasError())
	})
}
