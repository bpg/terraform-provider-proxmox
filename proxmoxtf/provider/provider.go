/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package provider

import (
	"context"
	"errors"
	"net/url"
	"os"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/bpg/terraform-provider-proxmox/proxmox"
	"github.com/bpg/terraform-provider-proxmox/proxmoxtf"
	"github.com/bpg/terraform-provider-proxmox/proxmoxtf/datasource"
	"github.com/bpg/terraform-provider-proxmox/proxmoxtf/datasource/firewall"
	"github.com/bpg/terraform-provider-proxmox/proxmoxtf/resource"
)

const (
	dvProviderVirtualEnvironmentEndpoint = ""
	dvProviderVirtualEnvironmentOTP      = ""
	dvProviderVirtualEnvironmentPassword = ""
	dvProviderVirtualEnvironmentUsername = ""

	mkProviderVirtualEnvironment         = "virtual_environment"
	mkProviderVirtualEnvironmentEndpoint = "endpoint"
	mkProviderVirtualEnvironmentInsecure = "insecure"
	mkProviderVirtualEnvironmentOTP      = "otp"
	mkProviderVirtualEnvironmentPassword = "password"
	mkProviderVirtualEnvironmentUsername = "username"
)

// ProxmoxVirtualEnvironment returns the object for this provider.
func ProxmoxVirtualEnvironment() *schema.Provider {
	return &schema.Provider{
		ConfigureContextFunc: providerConfigure,
		DataSourcesMap: map[string]*schema.Resource{
			"proxmox_virtual_environment_cluster_alias":   firewall.DataSourceVirtualEnvironmentFirewallAlias(),
			"proxmox_virtual_environment_cluster_aliases": firewall.DataSourceVirtualEnvironmentClusterAliases(),
			"proxmox_virtual_environment_datastores":      datasource.Datastores(),
			"proxmox_virtual_environment_dns":             datasource.DNS(),
			"proxmox_virtual_environment_group":           datasource.Group(),
			"proxmox_virtual_environment_groups":          datasource.Groups(),
			"proxmox_virtual_environment_hosts":           datasource.Hosts(),
			"proxmox_virtual_environment_nodes":           datasource.Nodes(),
			"proxmox_virtual_environment_pool":            datasource.Pool(),
			"proxmox_virtual_environment_pools":           datasource.Pools(),
			"proxmox_virtual_environment_role":            datasource.Role(),
			"proxmox_virtual_environment_roles":           datasource.Roles(),
			"proxmox_virtual_environment_time":            datasource.Time(),
			"proxmox_virtual_environment_user":            datasource.User(),
			"proxmox_virtual_environment_users":           datasource.Users(),
			"proxmox_virtual_environment_version":         datasource.Version(),
		},
		ResourcesMap: map[string]*schema.Resource{
			"proxmox_virtual_environment_certificate":   resource.ResourceVirtualEnvironmentCertificate(),
			"proxmox_virtual_environment_cluster_alias": resource.ResourceVirtualEnvironmentClusterAlias(),
			"proxmox_virtual_environment_cluster_ipset": resource.ResourceVirtualEnvironmentFirewallIPSet(),
			"proxmox_virtual_environment_container":     resource.ResourceVirtualEnvironmentContainer(),
			"proxmox_virtual_environment_dns":           resource.ResourceVirtualEnvironmentDNS(),
			"proxmox_virtual_environment_file":          resource.ResourceVirtualEnvironmentFile(),
			"proxmox_virtual_environment_group":         resource.ResourceVirtualEnvironmentGroup(),
			"proxmox_virtual_environment_hosts":         resource.ResourceVirtualEnvironmentHosts(),
			"proxmox_virtual_environment_pool":          resource.ResourceVirtualEnvironmentPool(),
			"proxmox_virtual_environment_role":          resource.ResourceVirtualEnvironmentRole(),
			"proxmox_virtual_environment_time":          resource.ResourceVirtualEnvironmentTime(),
			"proxmox_virtual_environment_user":          resource.ResourceVirtualEnvironmentUser(),
			"proxmox_virtual_environment_vm":            resource.VM(),
		},
		Schema: map[string]*schema.Schema{
			mkProviderVirtualEnvironment: {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						mkProviderVirtualEnvironmentEndpoint: {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The endpoint for the Proxmox Virtual Environment API",
							DefaultFunc: schema.MultiEnvDefaultFunc(
								[]string{"PROXMOX_VE_ENDPOINT", "PM_VE_ENDPOINT"},
								dvProviderVirtualEnvironmentEndpoint,
							),
							ValidateFunc: func(v interface{}, k string) (warns []string, errs []error) {
								value := v.(string)

								if value == "" {
									return []string{}, []error{
										errors.New(
											"you must specify an endpoint for the Proxmox Virtual Environment API (valid: https://host:port)",
										),
									}
								}

								_, err := url.ParseRequestURI(value)
								if err != nil {
									return []string{}, []error{
										errors.New(
											"you must specify a valid endpoint for the Proxmox Virtual Environment API (valid: https://host:port)",
										),
									}
								}

								return []string{}, []error{}
							},
						},
						mkProviderVirtualEnvironmentInsecure: {
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
						mkProviderVirtualEnvironmentOTP: {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The one-time password for the Proxmox Virtual Environment API",
							DefaultFunc: schema.MultiEnvDefaultFunc(
								[]string{"PROXMOX_VE_OTP", "PM_VE_OTP"},
								dvProviderVirtualEnvironmentOTP,
							),
						},
						mkProviderVirtualEnvironmentPassword: {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The password for the Proxmox Virtual Environment API",
							DefaultFunc: schema.MultiEnvDefaultFunc(
								[]string{"PROXMOX_VE_PASSWORD", "PM_VE_PASSWORD"},
								dvProviderVirtualEnvironmentPassword,
							),
							ValidateFunc: func(v interface{}, k string) (warns []string, errs []error) {
								value := v.(string)

								if value == "" {
									return []string{}, []error{
										errors.New(
											"you must specify a password for the Proxmox Virtual Environment API",
										),
									}
								}

								return []string{}, []error{}
							},
						},
						mkProviderVirtualEnvironmentUsername: {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The username for the Proxmox Virtual Environment API",
							DefaultFunc: schema.MultiEnvDefaultFunc(
								[]string{"PROXMOX_VE_USERNAME", "PM_VE_USERNAME"},
								dvProviderVirtualEnvironmentUsername,
							),
							ValidateFunc: func(v interface{}, k string) (warns []string, errs []error) {
								value := v.(string)

								if value == "" {
									return []string{}, []error{
										errors.New(
											"you must specify a username for the Proxmox Virtual Environment API (valid: username@realm)",
										),
									}
								}

								return []string{}, []error{}
							},
						},
					},
				},
				MaxItems: 1,
			},
		},
	}
}

func providerConfigure(_ context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
	var err error
	var veClient *proxmox.VirtualEnvironmentClient

	// Initialize the client for the Virtual Environment, if required.
	veConfigBlock := d.Get(mkProviderVirtualEnvironment).([]interface{})

	if len(veConfigBlock) > 0 {
		veConfig := veConfigBlock[0].(map[string]interface{})

		veClient, err = proxmox.NewVirtualEnvironmentClient(
			veConfig[mkProviderVirtualEnvironmentEndpoint].(string),
			veConfig[mkProviderVirtualEnvironmentUsername].(string),
			veConfig[mkProviderVirtualEnvironmentPassword].(string),
			veConfig[mkProviderVirtualEnvironmentOTP].(string),
			veConfig[mkProviderVirtualEnvironmentInsecure].(bool),
		)
		if err != nil {
			return nil, diag.FromErr(err)
		}
	}

	config := proxmoxtf.NewProviderConfiguration(veClient)

	return config, nil
}
