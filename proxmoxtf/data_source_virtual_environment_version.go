/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package proxmoxtf

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const (
	mkDataSourceVirtualEnvironmentVersionKeyboardLayout = "keyboard_layout"
	mkDataSourceVirtualEnvironmentVersionRelease        = "release"
	mkDataSourceVirtualEnvironmentVersionRepositoryID   = "repository_id"
	mkDataSourceVirtualEnvironmentVersionVersion        = "version"
)

func dataSourceVirtualEnvironmentVersion() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			mkDataSourceVirtualEnvironmentVersionKeyboardLayout: {
				Type:        schema.TypeString,
				Description: "The keyboard layout",
				Computed:    true,
				ForceNew:    true,
			},
			mkDataSourceVirtualEnvironmentVersionRelease: {
				Type:        schema.TypeString,
				Description: "The release information",
				Computed:    true,
				ForceNew:    true,
			},
			mkDataSourceVirtualEnvironmentVersionRepositoryID: {
				Type:        schema.TypeString,
				Description: "The repository id",
				Computed:    true,
				ForceNew:    true,
			},
			mkDataSourceVirtualEnvironmentVersionVersion: {
				Type:        schema.TypeString,
				Description: "The version information",
				Computed:    true,
				ForceNew:    true,
			},
		},
		Read: dataSourceVirtualEnvironmentVersionRead,
	}
}

func dataSourceVirtualEnvironmentVersionRead(d *schema.ResourceData, m interface{}) error {
	config := m.(providerConfiguration)
	veClient, err := config.GetVEClient()

	if err != nil {
		return err
	}

	version, err := veClient.Version()

	if err != nil {
		return err
	}

	d.SetId("version")

	d.Set(mkDataSourceVirtualEnvironmentVersionKeyboardLayout, version.Keyboard)
	d.Set(mkDataSourceVirtualEnvironmentVersionRelease, version.Release)
	d.Set(mkDataSourceVirtualEnvironmentVersionRepositoryID, version.RepositoryID)
	d.Set(mkDataSourceVirtualEnvironmentVersionVersion, version.Version)

	return nil
}
