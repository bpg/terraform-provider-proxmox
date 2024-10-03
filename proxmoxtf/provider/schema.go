/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package provider

import (
	"os"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

const (
	mkProviderEndpoint            = "endpoint"
	mkProviderInsecure            = "insecure"
	mkProviderMinTLS              = "min_tls"
	mkProviderAuthTicket          = "auth_ticket"
	mkProviderCSRFPreventionToken = "csrf_prevention_token" // #nosec G101
	mkProviderAPIToken            = "api_token"
	mkProviderOTP                 = "otp"
	mkProviderPassword            = "password"
	mkProviderUsername            = "username"
	mkProviderTmpDir              = "tmp_dir"
	mkProviderRandomVMIDs         = "random_vm_ids"
	mkProviderRandomVMIDStart     = "random_vm_id_start"
	mkProviderRandomVMIDEnd       = "random_vm_id_end"
	mkProviderSSH                 = "ssh"
	mkProviderSSHUsername         = "username"
	mkProviderSSHPassword         = "password"
	mkProviderSSHAgent            = "agent"
	mkProviderSSHAgentSocket      = "agent_socket"
	mkProviderSSHPrivateKey       = "private_key"
	mkProviderSSHSocks5Server     = "socks5_server"
	mkProviderSSHSocks5Username   = "socks5_username"
	mkProviderSSHSocks5Password   = "socks5_password"

	mkProviderSSHNode        = "node"
	mkProviderSSHNodeName    = "name"
	mkProviderSSHNodeAddress = "address"
	mkProviderSSHNodePort    = "port"
)

func createSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		mkProviderEndpoint: {
			Type:         schema.TypeString,
			Optional:     true,
			Description:  "The endpoint for the Proxmox VE API.",
			ValidateFunc: validation.IsURLWithHTTPorHTTPS,
		},
		mkProviderInsecure: {
			Type:        schema.TypeBool,
			Optional:    true,
			Description: "Whether to skip the TLS verification step.",
		},
		mkProviderMinTLS: {
			Type:     schema.TypeString,
			Optional: true,
			Description: "The minimum required TLS version for API calls." +
				"Supported values: `1.0|1.1|1.2|1.3`. Defaults to `1.3`.",
		},
		mkProviderAuthTicket: {
			Type:         schema.TypeString,
			Optional:     true,
			Sensitive:    true,
			Description:  "The pre-authenticated Ticket for the Proxmox VE API.",
			ValidateFunc: validation.StringIsNotEmpty,
		},
		mkProviderCSRFPreventionToken: {
			Type:         schema.TypeString,
			Optional:     true,
			Sensitive:    true,
			Description:  "The pre-authenticated CSRF Prevention Token for the Proxmox VE API.",
			ValidateFunc: validation.StringIsNotEmpty,
		},
		mkProviderAPIToken: {
			Type:         schema.TypeString,
			Optional:     true,
			Sensitive:    true,
			Description:  "The API token for the Proxmox VE API.",
			ValidateFunc: validation.StringIsNotEmpty,
		},
		mkProviderOTP: {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "The one-time password for the Proxmox VE API.",
			Deprecated: "The `otp` attribute is deprecated and will be removed in a future release. " +
				"Please use the `api_token` attribute instead.",
		},
		mkProviderUsername: {
			Type:         schema.TypeString,
			Optional:     true,
			Description:  "The username for the Proxmox VE API.",
			ValidateFunc: validation.StringIsNotEmpty,
		},
		mkProviderPassword: {
			Type:         schema.TypeString,
			Optional:     true,
			Sensitive:    true,
			Description:  "The password for the Proxmox VE API.",
			ValidateFunc: validation.StringIsNotEmpty,
		},
		mkProviderSSH: {
			Type:        schema.TypeList,
			Optional:    true,
			MaxItems:    1,
			Description: "The SSH configuration for the Proxmox nodes.",
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					mkProviderSSHUsername: {
						Type:     schema.TypeString,
						Optional: true,
						Description: "The username used for the SSH connection. " +
							"Defaults to the value of the `username` field of the " +
							"`provider` block.",
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
						Description: "The password used for the SSH connection. " +
							"Defaults to the value of the `password` field of the " +
							"`provider` block.",
						DefaultFunc: schema.MultiEnvDefaultFunc(
							[]string{"PROXMOX_VE_SSH_PASSWORD", "PM_VE_SSH_PASSWORD"},
							nil,
						),
						ValidateFunc: validation.StringIsNotEmpty,
					},
					mkProviderSSHAgent: {
						Type:     schema.TypeBool,
						Optional: true,
						Description: "Whether to use the SSH agent for authentication. Takes precedence over " +
							"the `private_key` and `password` fields. Defaults to the value of the " +
							"`PROXMOX_VE_SSH_AGENT` environment variable, or `false` if not set.",
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
							"environment variable.",
						DefaultFunc: schema.MultiEnvDefaultFunc(
							[]string{"SSH_AUTH_SOCK", "PROXMOX_VE_SSH_AUTH_SOCK", "PM_VE_SSH_AUTH_SOCK"},
							nil,
						),
						ValidateFunc: validation.StringIsNotEmpty,
					},
					mkProviderSSHPrivateKey: {
						Type:      schema.TypeString,
						Optional:  true,
						Sensitive: true,
						Description: "The unencrypted private key (in PEM format) used for the SSH connection. " +
							"Defaults to the value of the `PROXMOX_VE_SSH_PRIVATE_KEY` environment variable.",
					},
					mkProviderSSHSocks5Server: {
						Type:     schema.TypeString,
						Optional: true,
						Description: "The address:port of the SOCKS5 proxy server. " +
							"Defaults to the value of the `PROXMOX_VE_SSH_SOCKS5_SERVER` environment variable.",
						DefaultFunc: schema.MultiEnvDefaultFunc(
							[]string{"PROXMOX_VE_SSH_SOCKS5_SERVER"},
							nil,
						),
						ValidateFunc: validation.StringIsNotEmpty,
					},
					mkProviderSSHSocks5Username: {
						Type:     schema.TypeString,
						Optional: true,
						Description: "The username for the SOCKS5 proxy server. " +
							"Defaults to the value of the `PROXMOX_VE_SSH_SOCKS5_USERNAME` environment variable.",
						DefaultFunc: schema.MultiEnvDefaultFunc(
							[]string{"PROXMOX_VE_SSH_SOCKS5_USERNAME"},
							nil,
						),
						ValidateFunc: validation.StringIsNotEmpty,
					},
					mkProviderSSHSocks5Password: {
						Type:      schema.TypeString,
						Optional:  true,
						Sensitive: true,
						Description: "The password for the SOCKS5 proxy server. " +
							"Defaults to the value of the `PROXMOX_VE_SSH_SOCKS5_PASSWORD` environment variable.",
						DefaultFunc: schema.MultiEnvDefaultFunc(
							[]string{"PROXMOX_VE_SSH_SOCKS5_PASSWORD"},
							nil,
						),
						ValidateFunc: validation.StringIsNotEmpty,
					},
					mkProviderSSHNode: {
						Type:        schema.TypeList,
						Optional:    true,
						MinItems:    0,
						Description: "Overrides for SSH connection configuration for a Proxmox VE node.",
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								mkProviderSSHNodeName: {
									Type:         schema.TypeString,
									Required:     true,
									Description:  "The name of the Proxmox VE node.",
									ValidateFunc: validation.StringIsNotEmpty,
								},
								mkProviderSSHNodeAddress: {
									Type:         schema.TypeString,
									Required:     true,
									Description:  "The address of the Proxmox VE node.",
									ValidateFunc: validation.StringIsNotEmpty,
								},
								mkProviderSSHNodePort: {
									Type:         schema.TypeInt,
									Optional:     true,
									Description:  "The port of the Proxmox VE node.",
									Default:      22,
									ValidateFunc: validation.IsPortNumber,
								},
							},
						},
					},
				},
			},
		},
		mkProviderTmpDir: {
			Type:         schema.TypeString,
			Optional:     true,
			Description:  "The alternative temporary directory.",
			ValidateFunc: validation.StringIsNotEmpty,
		},
		mkProviderRandomVMIDs: {
			Type:        schema.TypeBool,
			Optional:    true,
			Description: "Whether to generate random VM / Container IDs.",
		},
		mkProviderRandomVMIDStart: {
			Type:         schema.TypeInt,
			Optional:     true,
			Description:  "The starting number for random VM / Container IDs.",
			ValidateFunc: validation.IntBetween(100, 999999999),
		},
		mkProviderRandomVMIDEnd: {
			Type:         schema.TypeInt,
			Optional:     true,
			Description:  "The ending number for random VM / Container IDs.",
			ValidateFunc: validation.IntBetween(100, 999999999),
		},
	}
}
