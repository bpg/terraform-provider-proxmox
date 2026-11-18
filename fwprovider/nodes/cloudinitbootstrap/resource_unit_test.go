/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package cloudinitbootstrap

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/assert"
)

func TestChangedSeedAttrs(t *testing.T) {
	t.Parallel()

	state := bootstrapModel{
		UserData:      types.StringValue("#cloud-config\n"),
		MetaData:      types.StringNull(),
		NetworkConfig: types.StringValue("version: 2\n"),
		VendorData:    types.StringNull(),
		Timeout:       types.Int64Value(1800),
	}

	same := state
	same.MetaData = types.StringValue("") // empty and unset build the same seed
	same.Timeout = types.Int64Value(60)   // not seed content
	assert.Empty(t, changedSeedAttrs(same, state))

	changed := state
	changed.UserData = types.StringValue("#cloud-config\npackages: [htop]\n")
	changed.VendorData = types.StringValue("#cloud-config\n")
	changed.NetworkConfig = types.StringNull()
	assert.Equal(t, []string{"user_data", "network_config", "vendor_data"}, changedSeedAttrs(changed, state))
}
