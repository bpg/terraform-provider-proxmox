package types

import (
	"encoding/json"
	"strconv"
	"testing"
)

func TestCustomInt64_UnmarshalJSON(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		json    string
		want    int64
		wantErr bool
	}{
		{"integer", "12345", 12345, false},
		{"quoted integer", "\"67890\"", 67890, false},
		{"scientific notation", "3.47726673451418e+15", 3477266734514180, false},
		{"negative integer", "-42", -42, false},
		{"float truncation", "123.99", 123, false},
		{"invalid string", "\"abc\"", 0, true},
		{"empty string", "\"\"", 0, true},
		{"int64 max", "9223372036854775807", 9223372036854775807, false},
		{"int64 min", "-9223372036854775808", -9223372036854775808, false},
		{"int64 max quoted", "\"9223372036854775807\"", 9223372036854775807, false},
		{"2^53-1 scientific", "9.007199254740991e+15", 9007199254740991, false},
		{"2^53+1 scientific", "9.007199254740993e+15", 9007199254740993, false},
		{"value larger than int64", "1e20", 0, true},
		{"empty string", "\"\"", 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			var v CustomInt64

			err := v.UnmarshalJSON([]byte(tt.json))
			if (err != nil) != tt.wantErr {
				t.Fatalf("error=%v, wantErr=%v", err, tt.wantErr)
			}

			if err == nil && int64(v) != tt.want {
				t.Fatalf("got=%d, want=%d", int64(v), tt.want)
			}
		})
	}
}

func TestCustomInt64_UnmarshalJSON_NullPointer(t *testing.T) {
	t.Parallel()

	var body struct {
		Data []struct {
			Total *CustomInt64 `json:"total"`
		} `json:"data"`
	}

	payload := []byte(`{"data":[{"total":null}]}`)
	if err := json.Unmarshal(payload, &body); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(body.Data) != 1 || body.Data[0].Total != nil {
		t.Fatalf("expected nil total pointer when JSON is null")
	}
}

func TestCustomInt_UnmarshalJSON(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		json    string
		want    int64
		wantErr bool
	}{
		{"plain integer", "789", 789, false},
		{"string integer", "\"456\"", 456, false},
		{"scientific notation numeric", "1.23e+02", 123, false},
		{"scientific notation string", "\"1.23e+02\"", 123, false},
		{"negative integer", "-42", -42, false},
		{"float truncation", "45.67", 45, false},
		{"large via float fallback", "3000000000", 3000000000, strconv.IntSize < 64},
		{"invalid string", "\"oops\"", 0, true},
		{"int32 max", "2147483647", 2147483647, false},
		{"int32 min", "-2147483648", -2147483648, false},
		{"int32 max+1", "2147483648", 2147483648, strconv.IntSize < 64},
		{"int32 min-1", "-2147483649", -2147483649, strconv.IntSize < 64},
		{"value larger than int", "1e20", 0, true},
		{"empty string", "\"\"", 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			var v CustomInt

			err := v.UnmarshalJSON([]byte(tt.json))
			if (err != nil) != tt.wantErr {
				t.Fatalf("error=%v, wantErr=%v", err, tt.wantErr)
			}

			if !tt.wantErr && int(v) != int(tt.want) {
				t.Fatalf("got=%d, want=%d", int(v), int(tt.want))
			}
		})
	}
}

func TestCustomInt_UnmarshalJSON_NullPointer(t *testing.T) {
	t.Parallel()

	var body struct {
		Data []struct {
			Count *CustomInt `json:"count"`
		} `json:"data"`
	}

	payload := []byte(`{"data":[{"count":null}]}`)
	if err := json.Unmarshal(payload, &body); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(body.Data) != 1 || body.Data[0].Count != nil {
		t.Fatalf("expected nil count pointer when JSON is null")
	}
}

// Note: intentionally scoped to only CustomInt and CustomInt64 UnmarshalJSON tests, aligned with Proxmox API numeric formats.
