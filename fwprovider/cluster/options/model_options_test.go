/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package options

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/bpg/terraform-provider-proxmox/proxmox/cluster"
)

// TestFromAPI_BandwidthLimit covers the parsing loop for the comma-separated `bwlimit`
// API value, especially the guard against malformed segments that would otherwise
// index out of bounds on `SplitN` returning a length-1 slice.
func TestFromAPI_BandwidthLimit(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name          string
		input         string
		expectErr     bool
		wantClone     *int64
		wantDefault   *int64
		wantMigration *int64
		wantMove      *int64
		wantRestore   *int64
	}{
		{
			name:      "happy path — all five keys",
			input:     "clone=100,default=200,migration=300,move=400,restore=500",
			wantClone: new(int64(100)), wantDefault: new(int64(200)),
			wantMigration: new(int64(300)), wantMove: new(int64(400)), wantRestore: new(int64(500)),
		},
		{
			name:      "trailing comma — must not panic",
			input:     "clone=100,",
			wantClone: new(int64(100)),
		},
		{
			name:      "bare token without equals — skipped, no panic",
			input:     "clone=100,garbage",
			wantClone: new(int64(100)),
		},
		{
			name:  "empty string — no keys set",
			input: "",
		},
		{
			name:      "unknown key — ignored, others still parsed",
			input:     "clone=100,future_key=42",
			wantClone: new(int64(100)),
		},
		{
			name:      "non-numeric value — error returned",
			input:     "clone=abc",
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			input := tt.input
			resp := &cluster.OptionsResponseData{}
			resp.BandwidthLimit = &input

			m := &clusterOptionsModel{}
			err := m.fromAPI(resp)

			if tt.expectErr {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			assertInt64(t, "clone", m.BandwidthLimitClone, tt.wantClone)
			assertInt64(t, "default", m.BandwidthLimitDefault, tt.wantDefault)
			assertInt64(t, "migration", m.BandwidthLimitMigration, tt.wantMigration)
			assertInt64(t, "move", m.BandwidthLimitMove, tt.wantMove)
			assertInt64(t, "restore", m.BandwidthLimitRestore, tt.wantRestore)
		})
	}
}

// assertInt64 asserts that the Terraform Int64 value matches the expected pointer:
// nil → null, non-nil → value equals the pointee.
func assertInt64(t *testing.T, name string, got attrInt64, want *int64) {
	t.Helper()

	if want == nil {
		assert.Truef(t, got.IsNull(), "%s: expected null, got %s", name, got.String())
		return
	}

	assert.Equalf(t, *want, got.ValueInt64(), "%s: value mismatch", name)
}

// attrInt64 is a minimal interface satisfied by terraform-plugin-framework's types.Int64,
// used here to keep the assertion helper decoupled from the concrete package import.
type attrInt64 interface {
	IsNull() bool
	ValueInt64() int64
	String() string
}
