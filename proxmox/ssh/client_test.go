package ssh

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBuildCommandLine(t *testing.T) {
	c := &client{}

	testCases := []struct {
		commands []string
		expected string
	}{
		{
			commands: []string{},
			expected: "",
		},
		{
			commands: []string{"echo 'Hello, World!'"},
			expected: `fun(){ echo 'Hello, World!' }; if [ $(sudo -n echo tfpve 2>&1 | grep "tfpve" | wc -l) -gt 0 ]; then sudo fun; else fun; fi`,
		},
		{
			commands: []string{"echo 'Hello'", "echo 'World'"},
			expected: `fun(){ echo 'Hello' && echo 'World' }; if [ $(sudo -n echo tfpve 2>&1 | grep "tfpve" | wc -l) -gt 0 ]; then sudo fun; else fun; fi`,
		},
	}

	for _, tc := range testCases {
		actual := c.buildCommandLine(tc.commands)
		assert.Equal(t, tc.expected, actual, "Unexpected command line")
	}
}
