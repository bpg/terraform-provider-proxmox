package proxmox

import (
	"testing"
)

func TestParseDiskSize(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		size    *string
		want    int
		wantErr bool
	}{
		{"handle null size", nil, 0, false},
		{"parse terabytes", strPtr("2T"), 2048, false},
		{"parse gigabytes", strPtr("2G"), 2, false},
		{"parse megabytes", strPtr("2048M"), 2, false},
		{"error on arbitrary string", strPtr("something"), -1, true},
		{"error on missing unit", strPtr("12345"), -1, true},
	}
	for _, test := range tests {
		tt := test
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got, err := ParseDiskSize(tt.size)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseDiskSize() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ParseDiskSize() got = %v, want %v", got, tt.want)
			}
		})
	}
}
