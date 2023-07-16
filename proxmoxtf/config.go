/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package proxmoxtf

import (
	"errors"

	"github.com/bpg/proxmox-api"
	"github.com/bpg/proxmox-api/rest"
	"github.com/bpg/proxmox-api/ssh"
)

// ProviderConfiguration is the configuration for the provider.
type ProviderConfiguration struct {
	restClient rest.Client
	sshClient  ssh.Client
}

// NewProviderConfiguration creates a new provider configuration.
func NewProviderConfiguration(
	restClient rest.Client,
	sshClient ssh.Client,
) ProviderConfiguration {
	return ProviderConfiguration{
		restClient: restClient,
		sshClient:  sshClient,
	}
}

// GetClient returns the Proxmox API client.
func (c *ProviderConfiguration) GetClient() (proxmox.Client, error) {
	if c.restClient == nil {
		return nil, errors.New(
			"you must specify the API access details in the provider configuration",
		)
	}

	if c.sshClient == nil {
		return nil, errors.New(
			"you must specify the SSH access details in the provider configuration",
		)
	}

	return proxmox.NewClient(c.restClient, c.sshClient), nil
}
