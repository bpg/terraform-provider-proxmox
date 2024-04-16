package validators

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCPUType(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		value string
		valid bool
	}{
		{"empty", "", false},
		{"invalid", "invalid", false},
		{"valid", "host", true},
		{"valid", "qemu64", true},
		{"valid", "custom-abc", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			f := CPUTypeValidator()
			res := f(tt.value, nil)

			if tt.valid {
				require.Empty(t, res, "validate: '%s'", tt.value)
			} else {
				require.NotEmpty(t, res, "validate: '%s'", tt.value)
			}
		})
	}
}
