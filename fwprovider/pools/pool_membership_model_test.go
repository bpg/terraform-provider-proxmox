package pools

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestParsePoolMembershipId(t *testing.T) {
	t.Parallel()
	tests := []struct {
		testName       string
		id             string
		expectedPoolId string
		expectedVmId   int64
		expectError    bool
	}{
		{"correct id", "test-pool/102", "test-pool", 102, false},
		{"wrong vm id format", "test-pool/asdlasd", "", 0, true},
		{"missing pool id", "102", "", 0, true},
		{"wrong id format", "test-pool/lxc/102", "", 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.testName, func(t *testing.T) {
			model, err := parseMembershipModelFromID(tt.id)
			if tt.expectError {
				assert.NotNil(t, err)
			} else {
				assert.Nil(t, err)
				assert.NotNil(t, model)
				assert.Equal(t, tt.id, model.ID.ValueString())
				assert.Equal(t, tt.expectedPoolId, model.PoolId.ValueString())
				assert.Equal(t, tt.expectedVmId, model.VmID.ValueInt64())
			}
		})
	}
}
