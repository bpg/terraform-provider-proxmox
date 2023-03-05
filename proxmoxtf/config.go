/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package proxmoxtf

import (
	"errors"

	"github.com/bpg/terraform-provider-proxmox/proxmox"
)

type ProviderConfiguration struct {
	veClient *proxmox.VirtualEnvironmentClient
}

func NewProviderConfiguration(veClient *proxmox.VirtualEnvironmentClient) ProviderConfiguration {
	return ProviderConfiguration{
		veClient: veClient,
	}
}

func (c *ProviderConfiguration) GetVEClient() (*proxmox.VirtualEnvironmentClient, error) {
	if c.veClient == nil {
		return nil, errors.New(
			"you must specify the virtual environment details in the provider configuration",
		)
	}

	return c.veClient, nil
}
