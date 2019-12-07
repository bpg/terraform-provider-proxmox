/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package main

import (
	"errors"

	"github.com/hashicorp/terraform/helper/schema"
)

const (
	mkDataSourceVirtualEnvironmentVersionKeyboard     = "keyboard"
	mkDataSourceVirtualEnvironmentVersionRelease      = "release"
	mkDataSourceVirtualEnvironmentVersionRepositoryID = "repository_id"
	mkDataSourceVirtualEnvironmentVersionVersion      = "version"
)

// dataSourceVirtualEnvironmentVersion retrieves version information for a Proxmox installation.
func dataSourceVirtualEnvironmentVersion() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			mkDataSourceVirtualEnvironmentVersionKeyboard: &schema.Schema{
				Type:        schema.TypeString,
				Description: "The keyboard layout",
				Computed:    true,
				ForceNew:    true,
			},
			mkDataSourceVirtualEnvironmentVersionRelease: &schema.Schema{
				Type:        schema.TypeString,
				Description: "The release information",
				Computed:    true,
				ForceNew:    true,
			},
			mkDataSourceVirtualEnvironmentVersionRepositoryID: &schema.Schema{
				Type:        schema.TypeString,
				Description: "The repository id",
				Computed:    true,
				ForceNew:    true,
			},
			mkDataSourceVirtualEnvironmentVersionVersion: &schema.Schema{
				Type:        schema.TypeString,
				Description: "The version information",
				Computed:    true,
				ForceNew:    true,
			},
		},
		Read: dataSourceVirtualEnvironmentVersionRead,
	}
}

// dataSourceVirtualEnvironmentVersionRead reads version information for a Proxmox installation.
func dataSourceVirtualEnvironmentVersionRead(d *schema.ResourceData, m interface{}) error {
	config := m.(providerConfiguration)

	if config.veClient == nil {
		return errors.New("You must specify the virtual environment details in the provider configuration to use this data source")
	}

	version, err := config.veClient.Version()

	if err != nil {
		return err
	}

	d.SetId(version.Version)

	d.Set(mkDataSourceVirtualEnvironmentVersionKeyboard, version.Keyboard)
	d.Set(mkDataSourceVirtualEnvironmentVersionRelease, version.Release)
	d.Set(mkDataSourceVirtualEnvironmentVersionRepositoryID, version.RepositoryID)
	d.Set(mkDataSourceVirtualEnvironmentVersionVersion, version.Version)

	return nil
}
