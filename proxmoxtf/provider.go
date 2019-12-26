/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package proxmoxtf

import (
	"errors"
	"net/url"
	"os"

	"github.com/danitso/terraform-provider-proxmox/proxmox"
	"github.com/hashicorp/terraform/helper/schema"
)

const (
	mkProviderVirtualEnvironment         = "virtual_environment"
	mkProviderVirtualEnvironmentEndpoint = "endpoint"
	mkProviderVirtualEnvironmentInsecure = "insecure"
	mkProviderVirtualEnvironmentPassword = "password"
	mkProviderVirtualEnvironmentUsername = "username"
)

type providerConfiguration struct {
	veClient *proxmox.VirtualEnvironmentClient
}

// Provider returns the object for this provider.
func Provider() *schema.Provider {
	return &schema.Provider{
		ConfigureFunc: providerConfigure,
		DataSourcesMap: map[string]*schema.Resource{
			"proxmox_virtual_environment_datastores": dataSourceVirtualEnvironmentDatastores(),
			"proxmox_virtual_environment_group":      dataSourceVirtualEnvironmentGroup(),
			"proxmox_virtual_environment_groups":     dataSourceVirtualEnvironmentGroups(),
			"proxmox_virtual_environment_nodes":      dataSourceVirtualEnvironmentNodes(),
			"proxmox_virtual_environment_pool":       dataSourceVirtualEnvironmentPool(),
			"proxmox_virtual_environment_pools":      dataSourceVirtualEnvironmentPools(),
			"proxmox_virtual_environment_role":       dataSourceVirtualEnvironmentRole(),
			"proxmox_virtual_environment_roles":      dataSourceVirtualEnvironmentRoles(),
			"proxmox_virtual_environment_user":       dataSourceVirtualEnvironmentUser(),
			"proxmox_virtual_environment_users":      dataSourceVirtualEnvironmentUsers(),
			"proxmox_virtual_environment_version":    dataSourceVirtualEnvironmentVersion(),
		},
		ResourcesMap: map[string]*schema.Resource{
			"proxmox_virtual_environment_file":  resourceVirtualEnvironmentFile(),
			"proxmox_virtual_environment_group": resourceVirtualEnvironmentGroup(),
			"proxmox_virtual_environment_pool":  resourceVirtualEnvironmentPool(),
			"proxmox_virtual_environment_role":  resourceVirtualEnvironmentRole(),
			"proxmox_virtual_environment_user":  resourceVirtualEnvironmentUser(),
			"proxmox_virtual_environment_vm":    resourceVirtualEnvironmentVM(),
		},
		Schema: map[string]*schema.Schema{
			mkProviderVirtualEnvironment: &schema.Schema{
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
								"",
							),
							ValidateFunc: func(v interface{}, k string) (warns []string, errs []error) {
								value := v.(string)

								if value == "" {
									return []string{}, []error{
										errors.New("You must specify an endpoint for the Proxmox Virtual Environment API (valid: https://host:port)"),
									}
								}

								_, err := url.ParseRequestURI(value)

								if err != nil {
									return []string{}, []error{
										errors.New("You must specify a valid endpoint for the Proxmox Virtual Environment API (valid: https://host:port)"),
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
						mkProviderVirtualEnvironmentPassword: {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The password for the Proxmox Virtual Environment API",
							DefaultFunc: schema.MultiEnvDefaultFunc(
								[]string{"PROXMOX_VE_PASSWORD", "PM_VE_PASSWORD"},
								"",
							),
							ValidateFunc: func(v interface{}, k string) (warns []string, errs []error) {
								value := v.(string)

								if value == "" {
									return []string{}, []error{
										errors.New("You must specify a password for the Proxmox Virtual Environment API"),
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
								"",
							),
							ValidateFunc: func(v interface{}, k string) (warns []string, errs []error) {
								value := v.(string)

								if value == "" {
									return []string{}, []error{
										errors.New("You must specify a username for the Proxmox Virtual Environment API (valid: username@realm)"),
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

func providerConfigure(d *schema.ResourceData) (interface{}, error) {
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
			veConfig[mkProviderVirtualEnvironmentInsecure].(bool),
		)

		if err != nil {
			return nil, err
		}
	}

	config := providerConfiguration{
		veClient: veClient,
	}

	return config, nil
}

func (c *providerConfiguration) GetVEClient() (*proxmox.VirtualEnvironmentClient, error) {
	if c.veClient == nil {
		return nil, errors.New("You must specify the virtual environment details in the provider configuration")
	}

	return c.veClient, nil
}
