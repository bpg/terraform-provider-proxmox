/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

// file deepcode ignore NoHardcodedCredentials/test: test file

package tasks

import (
	"reflect"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestParseTaskID(t *testing.T) {
	t.Parallel()

	parseTime := func(s string) time.Time {
		stime, err := time.Parse(time.RFC3339, s)
		require.NoError(t, err)

		return stime.UTC()
	}

	tests := []struct {
		name    string
		taskID  string
		want    TaskID
		wantErr bool
	}{
		{
			name:   "imgcopy task",
			taskID: "UPID:pve:00061CB3:010BA69C:64EFECB0:imgcopy::root@pam:",
			want: TaskID{
				NodeName:  "pve",
				PID:       400563,
				PStart:    17540764,
				StartTime: parseTime("2023-08-30T21:28:16-04:00"),
				Type:      "imgcopy",
				ID:        "",
				User:      "root@pam",
			},
		},
		{
			name:   "qmcreate task",
			taskID: "UPID:pve:00061CB3:010BA69C:64EFECB0:qmcreate:101:root@pam:",
			want: TaskID{
				NodeName:  "pve",
				PID:       400563,
				PStart:    17540764,
				StartTime: parseTime("2023-08-30T21:28:16-04:00"),
				Type:      "qmcreate",
				ID:        "101",
				User:      "root@pam",
			},
		},
		{
			name:   "large pid",
			taskID: "UPID:pve-srv:0012D68C:8206F38C:65144CB7:qmcreate:501:root@pam:",
			want: TaskID{
				NodeName:  "pve-srv",
				PID:       1234572,
				PStart:    2181493644,
				StartTime: parseTime("2023-09-27T15:39:35Z"),
				Type:      "qmcreate",
				ID:        "501",
				User:      "root@pam",
			},
		},
		{
			name:    "missing node",
			taskID:  "UPID::00061CB3:010BA69C:64EFECB0:qmcreate:101:root@pam:",
			wantErr: true,
		},
		{
			name:    "wrong ID format",
			taskID:  "blah",
			wantErr: true,
		},
		{
			name:    "missing pid",
			taskID:  "UPID:pve::010BA69C:64EFECB0:qmcreate:101:root@pam:",
			wantErr: true,
		},
		{
			name:    "missing parts",
			taskID:  "UPID:pve:00061CB3:010BA69C:64EFECB0::root@pam:",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		tt := tt // capture range variable
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, err := ParseTaskID(tt.taskID)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseTaskID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ParseTaskID() got = %v, want %v", got, tt.want)
			}
		})
	}
}
