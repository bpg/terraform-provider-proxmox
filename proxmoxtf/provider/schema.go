/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package provider

import (
	"errors"
	"net/url"
	"os"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func createSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
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
	}
}
