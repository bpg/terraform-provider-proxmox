/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package ssh

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestShellQuote(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		in   string
		want string
	}{
		{"plain", "file.txt", `'file.txt'`},
		{"empty", "", `''`},
		{"semicolon and space", "a; touch /tmp/x", `'a; touch /tmp/x'`},
		{"command substitution", "$(id)", `'$(id)'`},
		{"backticks", "`id`", "'`id`'"},
		{"single quote", "it's", `'it'"'"'s'`},
		{"double quote", `say "hi"`, `'say "hi"'`},
		{"newline", "a\nb", "'a\nb'"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			require.Equal(t, tt.want, ShellQuote(tt.in))
		})
	}
}
