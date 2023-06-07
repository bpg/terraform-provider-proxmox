/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package provider

import (
	"fmt"
	"os"
	"regexp"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

const (
	dvProviderOTP                = ""
	mkProviderVirtualEnvironment = "virtual_environment"
	mkProviderEndpoint           = "endpoint"
	mkProviderInsecure           = "insecure"
	mkProviderOTP                = "otp"
	mkProviderPassword           = "password"
	mkProviderUsername           = "username"
	mkProviderAPIToken           = "api_token"
	mkProviderSSH                = "ssh"
	mkProviderSSHUsername        = "username"
	mkProviderSSHPassword        = "password"
	mkProviderSSHAgent           = "agent"
	mkProviderSSHAgentSocket     = "agent_socket"

	mkProviderSSHNode        = "node"
	mkProviderSSHNodeName    = "name"
	mkProviderSSHNodeAddress = "address"
)

func createSchema() map[string]*schema.Schema {
	providerSchema := nestedProviderSchema()
	providerSchema[mkProviderVirtualEnvironment] = &schema.Schema{
		Type:     schema.TypeList,
		Optional: true,
		Elem: &schema.Resource{
			Schema: nestedProviderSchema(),
		},
		MaxItems:   1,
		Deprecated: "Move attributes out of virtual_environment block",
	}

	return providerSchema
}

func nestedProviderSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		mkProviderEndpoint: {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "The endpoint for the Proxmox Virtual Environment API",
			DefaultFunc: schema.MultiEnvDefaultFunc(
				[]string{"PROXMOX_VE_ENDPOINT", "PM_VE_ENDPOINT"},
				nil,
			),
			AtLeastOneOf: []string{
				mkProviderEndpoint,
				fmt.Sprintf("%s.0.%s", mkProviderVirtualEnvironment, mkProviderEndpoint),
			},
			ValidateFunc: validation.IsURLWithHTTPorHTTPS,
		},
		mkProviderInsecure: {
			Type:        schema.TypeBool,
			Optional:    true,
			Description: "Whether to skip the TLS verification step",
			DefaultFunc: func() (interface{}, error) {
				for _, k := range []string{"PROXMOX_VE_INSECURE", "PM_VE_INSECURE"} {
					v := os.Getenv(k)

					if v == "true" || v == "1" {
						return true, nil
					}
				}

				return false, nil
			},
		},
		mkProviderOTP: {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "The one-time password for the Proxmox Virtual Environment API",
			DefaultFunc: schema.MultiEnvDefaultFunc(
				[]string{"PROXMOX_VE_OTP", "PM_VE_OTP"},
				dvProviderOTP,
			),
		},
		mkProviderPassword: {
			Type:        schema.TypeString,
			Optional:    true,
			Sensitive:   true,
			Description: "The password for the Proxmox Virtual Environment API",
			DefaultFunc: schema.MultiEnvDefaultFunc(
				[]string{"PROXMOX_VE_PASSWORD", "PM_VE_PASSWORD"},
				nil,
			),
			AtLeastOneOf: []string{
				mkProviderPassword,
				fmt.Sprintf("%s.0.%s", mkProviderVirtualEnvironment, mkProviderPassword),
				mkProviderAPIToken,
			},
			ValidateFunc: validation.StringIsNotEmpty,
		},
		mkProviderUsername: {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "The username for the Proxmox Virtual Environment API",
			DefaultFunc: schema.MultiEnvDefaultFunc(
				[]string{"PROXMOX_VE_USERNAME", "PM_VE_USERNAME"},
				nil,
			),
			AtLeastOneOf: []string{
				mkProviderUsername,
				fmt.Sprintf("%s.0.%s", mkProviderVirtualEnvironment, mkProviderUsername),
				mkProviderAPIToken,
			},
			ValidateFunc: validation.StringIsNotEmpty,
		},
		mkProviderAPIToken: {
			Type:        schema.TypeString,
			Optional:    true,
			Sensitive:   true,
			Description: "The API token for the Proxmox Virtual Environment API",
			DefaultFunc: schema.MultiEnvDefaultFunc(
				[]string{"PROXMOX_VE_API_TOKEN", "PM_VE_API_TOKEN"},
				nil,
			),
			ValidateDiagFunc: validation.ToDiagFunc(validation.StringMatch(
				regexp.MustCompile(`^\S+@\w+!\S+=([a-zA-Z0-9-]+)$`),
				"Must be a valid API token, e.g. 'USER@REALM!TOKENID=UUID'",
			)),
		},
		mkProviderSSH: {
			Type:        schema.TypeList,
			Optional:    true,
			MaxItems:    1,
			Description: "The SSH connection configuration to a Proxmox node",
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					mkProviderSSHUsername: {
						Type:     schema.TypeString,
						Optional: true,
						Description: fmt.Sprintf("The username used for the SSH connection, "+
							"defaults to the user specified in '%s'", mkProviderUsername),
						DefaultFunc: schema.MultiEnvDefaultFunc(
							[]string{"PROXMOX_VE_SSH_USERNAME", "PM_VE_SSH_USERNAME"},
							nil,
						),
						ValidateFunc: validation.StringIsNotEmpty,
					},
					mkProviderSSHPassword: {
						Type:      schema.TypeString,
						Optional:  true,
						Sensitive: true,
						Description: fmt.Sprintf("The password used for the SSH connection, "+
							"defaults to the password specified in '%s'", mkProviderPassword),
						DefaultFunc: schema.MultiEnvDefaultFunc(
							[]string{"PROXMOX_VE_SSH_PASSWORD", "PM_VE_SSH_PASSWORD"},
							nil,
						),
						ValidateFunc: validation.StringIsNotEmpty,
					},
					mkProviderSSHAgent: {
						Type:        schema.TypeBool,
						Optional:    true,
						Description: "Whether to use the SSH agent for the SSH authentication. Defaults to false",
						DefaultFunc: func() (interface{}, error) {
							for _, k := range []string{"PROXMOX_VE_SSH_AGENT", "PM_VE_SSH_AGENT"} {
								v := os.Getenv(k)

								if v == "true" || v == "1" {
									return true, nil
								}
							}

							return false, nil
						},
					},
					mkProviderSSHAgentSocket: {
						Type:     schema.TypeString,
						Optional: true,
						Description: "The path to the SSH agent socket. Defaults to the value of the `SSH_AUTH_SOCK` " +
							"environment variable",
						DefaultFunc: schema.MultiEnvDefaultFunc(
							[]string{"SSH_AUTH_SOCK", "PROXMOX_VE_SSH_AUTH_SOCK", "PM_VE_SSH_AUTH_SOCK"},
							nil,
						),
						ValidateFunc: validation.StringIsNotEmpty,
					},
					mkProviderSSHNode: {
						Type:        schema.TypeList,
						Optional:    true,
						MinItems:    0,
						Description: "Overrides for SSH connection configuration to a Proxmox node",
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								mkProviderSSHNodeName: {
									Type:         schema.TypeString,
									Required:     true,
									Description:  "The name of the node to connect to",
									ValidateFunc: validation.StringIsNotEmpty,
								},
								mkProviderSSHNodeAddress: {
									Type:         schema.TypeString,
									Required:     true,
									Description:  "The address that should be used to connect to the node",
									ValidateFunc: validation.IsIPAddress,
								},
							},
						},
					},
				},
			},
		},
	}
}
