/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL can be obtained
 * at http://mozilla.org/MPL/2.0/.
 */

package vms

import (
	"encoding/json"
	"net/url"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCustomAgentFreezeFSAliases(t *testing.T) {
	t.Parallel()

	for _, agentString := range []string{
		"enabled=1,freeze-fs=0",
		"enabled=1,freeze-fs-on-backup=0",
		"enabled=1,guest-fsfreeze=0",
	} {
		t.Run(agentString, func(t *testing.T) {
			t.Parallel()

			var agent CustomAgent
			require.NoError(t, json.Unmarshal([]byte(`"`+agentString+`"`), &agent))
			require.NotNil(t, agent.FreezeFS)
			require.False(t, bool(*agent.FreezeFS))

			values := url.Values{}
			require.NoError(t, agent.EncodeValues("agent", &values))
			require.Contains(t, values.Get("agent"), "freeze-fs-on-backup=0")
		})
	}
}
