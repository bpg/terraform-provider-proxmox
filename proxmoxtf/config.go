/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package proxmoxtf

import (
	"errors"

	"github.com/bpg/terraform-provider-proxmox/proxmox"
	"github.com/bpg/terraform-provider-proxmox/proxmox/api"
	"github.com/bpg/terraform-provider-proxmox/proxmox/ssh"
)

// ProviderConfiguration is the configuration for the provider.
type ProviderConfiguration struct {
	apiClient api.Client
	sshClient ssh.Client
}

// NewProviderConfiguration creates a new provider configuration.
func NewProviderConfiguration(
	veClient api.Client,
	sshClient ssh.Client,
) ProviderConfiguration {
	return ProviderConfiguration{
		apiClient: veClient,
		sshClient: sshClient,
	}
}

// GetAPIClient returns the virtual environment client.
// Deprecated: Use API() instead.
func (c *ProviderConfiguration) GetAPIClient() (api.Client, error) {
	if c.apiClient == nil {
		return nil, errors.New(
			"you must specify the API access details in the provider configuration",
		)
	}

	return c.apiClient, nil
}

// GetSSHClient returns the SSH client.
// Deprecated: Use API() instead.
func (c *ProviderConfiguration) GetSSHClient() (ssh.Client, error) {
	if c.sshClient == nil {
		return nil, errors.New(
			"you must specify the SSH access details in the provider configuration",
		)
	}

	return c.sshClient, nil
}

// GetAPI returns the Proxmox API client.
func (c *ProviderConfiguration) GetAPI() (proxmox.API, error) {
	a, err := c.GetAPIClient()
	if err != nil {
		return nil, err
	}

	s, err := c.GetSSHClient()
	if err != nil {
		return nil, err
	}

	return proxmox.NewClient(a, s), nil
}
