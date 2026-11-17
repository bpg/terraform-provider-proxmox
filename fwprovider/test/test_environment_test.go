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
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

// TestAccEnvironmentExecuteNodeCommands verifies that ExecuteNodeCommands runs as root on the test node.
func TestAccEnvironmentExecuteNodeCommands(t *testing.T) {
	te := InitEnvironment(t)

	out := te.ExecuteNodeCommands([]string{"id -u"})
	require.Equal(t, "0", strings.TrimSpace(out))
}
