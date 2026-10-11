/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package backup

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/require"

	"github.com/bpg/terraform-provider-proxmox/proxmox/cluster/backup"
)

func TestNotificationModeToAPI(t *testing.T) {
	t.Parallel()

	for _, mode := range []string{"auto", "legacy-sendmail", "notification-system"} {
		t.Run("create "+mode, func(t *testing.T) {
			t.Parallel()

			var diags diag.Diagnostics

			m := &backupJobModel{NotificationMode: types.StringValue(mode)}
			body := m.toAPICreate(context.Background(), &diags)

			require.False(t, diags.HasError())
			require.Equal(t, new(mode), body.NotificationMode)
		})
	}

	t.Run("create without notification_mode", func(t *testing.T) {
		t.Parallel()

		var diags diag.Diagnostics

		m := &backupJobModel{NotificationMode: types.StringNull()}
		body := m.toAPICreate(context.Background(), &diags)

		require.Nil(t, body.NotificationMode)

		raw, err := json.Marshal(body)
		require.NoError(t, err)
		require.NotContains(t, string(raw), "notification-mode")
	})

	t.Run("request body field name", func(t *testing.T) {
		t.Parallel()

		raw, err := json.Marshal(&backup.RequestBodyCommon{NotificationMode: new("notification-system")})
		require.NoError(t, err)
		require.JSONEq(t, `{"notification-mode":"notification-system"}`, string(raw))
	})
}

func TestNotificationModeUpdate(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		plan       types.String
		state      types.String
		wantMode   *string
		wantDelete bool
	}{
		{"changed", types.StringValue("legacy-sendmail"), types.StringValue("auto"), new("legacy-sendmail"), false},
		{"unchanged", types.StringValue("auto"), types.StringValue("auto"), new("auto"), false},
		{"removed from config", types.StringNull(), types.StringValue("legacy-sendmail"), nil, true},
		{"never set", types.StringNull(), types.StringNull(), nil, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			var diags diag.Diagnostics

			plan := &backupJobModel{NotificationMode: tt.plan}
			state := &backupJobModel{NotificationMode: tt.state}
			body := plan.toAPIUpdate(context.Background(), state, &diags)

			require.False(t, diags.HasError())
			require.Equal(t, tt.wantMode, body.NotificationMode)

			if tt.wantDelete {
				require.Contains(t, body.Delete, "notification-mode")
			} else {
				require.NotContains(t, body.Delete, "notification-mode")
			}
		})
	}
}

func TestNotificationModeFromAPI(t *testing.T) {
	t.Parallel()

	data := &backup.GetResponseData{ID: "job-1", Schedule: "daily", Storage: "local"}

	t.Run("resource set", func(t *testing.T) {
		t.Parallel()

		d := *data
		d.NotificationMode = new("notification-system")

		m := &backupJobModel{}
		require.False(t, m.fromAPI(context.Background(), &d).HasError())
		require.Equal(t, types.StringValue("notification-system"), m.NotificationMode)
	})

	t.Run("resource unset", func(t *testing.T) {
		t.Parallel()

		m := &backupJobModel{}
		require.False(t, m.fromAPI(context.Background(), data).HasError())
		require.True(t, m.NotificationMode.IsNull())
	})

	t.Run("data source", func(t *testing.T) {
		t.Parallel()

		d := *data
		d.NotificationMode = new("legacy-sendmail")

		m := &backupJobDatasourceModel{}
		m.fromAPI(&d)
		require.Equal(t, types.StringValue("legacy-sendmail"), m.NotificationMode)
	})

	t.Run("response field name", func(t *testing.T) {
		t.Parallel()

		var got backup.GetResponseData

		require.NoError(t, json.Unmarshal([]byte(`{"id":"j","schedule":"daily","storage":"local","notification-mode":"auto"}`), &got))
		require.Equal(t, new("auto"), got.NotificationMode)
	})
}
