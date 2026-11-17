//go:build acceptance || all

//testacc:tier=light
//testacc:resource=misc

/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package test

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestAccSSHExecuteNodeCommandsError(t *testing.T) {
	t.Parallel()

	te := InitEnvironment(t)

	tests := []struct {
		name     string
		command  string
		contains []string
	}{
		{"no output", "exit 3", []string{"exited with status 3"}},
		{"with output", "echo boom >&2; exit 4", []string{"exited with status 4", "boom"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			_, err := te.SSHClient().ExecuteNodeCommands(t.Context(), te.NodeName, []string{tt.command})
			require.Error(t, err)

			for _, s := range tt.contains {
				require.ErrorContains(t, err, s)
			}
		})
	}
}
