/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package proxmoxtf

import (
	"errors"

	"github.com/bpg/terraform-provider-proxmox/proxmox"
	"github.com/bpg/terraform-provider-proxmox/proxmox/types"
)

// ProviderConfiguration is the configuration for the provider.
type ProviderConfiguration struct {
	veClient  *proxmox.VirtualEnvironmentClient
	sshClient types.SSHClient
}

// NewProviderConfiguration creates a new provider configuration.
func NewProviderConfiguration(
	veClient *proxmox.VirtualEnvironmentClient,
	sshClient types.SSHClient,
) ProviderConfiguration {
	return ProviderConfiguration{
		veClient:  veClient,
		sshClient: sshClient,
	}
}

// GetVEClient returns the virtual environment client.
func (c *ProviderConfiguration) GetVEClient() (*proxmox.VirtualEnvironmentClient, error) {
	if c.veClient == nil {
		return nil, errors.New(
			"you must specify the API access details in the provider configuration",
		)
	}

	return c.veClient, nil
}

// GetSSHClient returns the SSH client.
func (c *ProviderConfiguration) GetSSHClient() (types.SSHClient, error) {
	if c.sshClient == nil {
		return nil, errors.New(
			"you must specify the SSH access details in the provider configuration",
		)
	}

	return c.sshClient, nil
}
