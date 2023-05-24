/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package provider

import (
	"context"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/bpg/terraform-provider-proxmox/proxmox"
	"github.com/bpg/terraform-provider-proxmox/proxmox/ssh"
	"github.com/bpg/terraform-provider-proxmox/proxmoxtf"
)

const (
	dvProviderOTP                = ""
	mkProviderVirtualEnvironment = "virtual_environment"
	mkProviderEndpoint           = "endpoint"
	mkProviderInsecure           = "insecure"
	mkProviderOTP                = "otp"
	mkProviderPassword           = "password"
	mkProviderUsername           = "username"
	mkProviderSSH                = "ssh"
	mkProviderSSHUsername        = "username"
	mkProviderSSHPassword        = "password"
	mkProviderSSHAgent           = "agent"
	mkProviderSSHAgentSocket     = "agent_socket"
)

// ProxmoxVirtualEnvironment returns the object for this provider.
func ProxmoxVirtualEnvironment() *schema.Provider {
	return &schema.Provider{
		ConfigureContextFunc: providerConfigure,
		DataSourcesMap:       createDatasourceMap(),
		ResourcesMap:         createResourceMap(),
		Schema:               createSchema(),
	}
}

func providerConfigure(_ context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
	var err error

	var veClient *proxmox.VirtualEnvironmentClient

	var sshClient ssh.Client

	// Legacy configuration, wrapped in the deprecated `virtual_environment` block
	veConfigBlock := d.Get(mkProviderVirtualEnvironment).([]interface{})
	//nolint:nestif
	if len(veConfigBlock) > 0 {
		veConfig := veConfigBlock[0].(map[string]interface{})

		username := veConfig[mkProviderUsername].(string)
		password := veConfig[mkProviderPassword].(string)

		veClient, err = proxmox.NewVirtualEnvironmentClient(
			veConfig[mkProviderEndpoint].(string),
			username,
			password,
			veConfig[mkProviderOTP].(string),
			veConfig[mkProviderInsecure].(bool),
		)
		if err != nil {
			return nil, diag.Errorf("error creating virtual environment client: %s", err)
		}

		veSSHConfig := veConfig[mkProviderSSH].(map[string]interface{})

		sshUsername := veSSHConfig[mkProviderSSHUsername].(string)
		if sshUsername == "" {
			sshUsername = strings.Split(username, "@")[0]
		}

		sshPassword := veSSHConfig[mkProviderSSHPassword].(string)
		if sshPassword == "" {
			sshPassword = password
		}

		sshClient, err = ssh.NewSSHClient(
			sshUsername,
			sshPassword,
			veSSHConfig[mkProviderSSHAgent].(bool),
			veSSHConfig[mkProviderSSHAgentSocket].(string),
		)
		if err != nil {
			return nil, diag.Errorf("error creating SSH client: %s", err)
		}
	} else {
		username := d.Get(mkProviderUsername).(string)
		password := d.Get(mkProviderPassword).(string)
		veClient, err = proxmox.NewVirtualEnvironmentClient(
			d.Get(mkProviderEndpoint).(string),
			username,
			password,
			d.Get(mkProviderOTP).(string),
			d.Get(mkProviderInsecure).(bool),
		)
		if err != nil {
			return nil, diag.Errorf("error creating virtual environment client: %s", err)
		}

		sshconf := map[string]interface{}{
			mkProviderSSHUsername:    username,
			mkProviderSSHPassword:    password,
			mkProviderSSHAgent:       false,
			mkProviderSSHAgentSocket: "",
		}

		sshBlock, sshSet := d.GetOk(mkProviderSSH)
		if sshSet {
			sshconf = sshBlock.(*schema.Set).List()[0].(map[string]interface{})
		}

		sshClient, err = ssh.NewSSHClient(
			sshconf[mkProviderSSHUsername].(string),
			sshconf[mkProviderSSHPassword].(string),
			sshconf[mkProviderSSHAgent].(bool),
			sshconf[mkProviderSSHAgentSocket].(string),
		)
		if err != nil {
			return nil, diag.Errorf("error creating SSH client: %s", err)
		}
	}

	config := proxmoxtf.NewProviderConfiguration(veClient, sshClient)

	return config, nil
}
