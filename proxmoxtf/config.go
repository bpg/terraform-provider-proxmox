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

// ProviderConfiguration is the configuration for the provider.
type ProviderConfiguration struct {
	veClient *proxmox.VirtualEnvironmentClient
}

// NewProviderConfiguration creates a new provider configuration.
func NewProviderConfiguration(veClient *proxmox.VirtualEnvironmentClient) ProviderConfiguration {
	return ProviderConfiguration{
		veClient: veClient,
	}
}

// GetVEClient returns the virtual environment client.
func (c *ProviderConfiguration) GetVEClient() (*proxmox.VirtualEnvironmentClient, error) {
	if c.veClient == nil {
		return nil, errors.New(
			"you must specify the virtual environment details in the provider configuration",
		)
	}

	return c.veClient, nil
}
