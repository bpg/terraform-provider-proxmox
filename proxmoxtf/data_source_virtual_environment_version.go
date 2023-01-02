/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package proxmoxtf

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
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
		ReadContext: dataSourceVirtualEnvironmentVersionRead,
	}
}

func dataSourceVirtualEnvironmentVersionRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	config := m.(providerConfiguration)
	veClient, err := config.GetVEClient()
	if err != nil {
		return diag.FromErr(err)
	}

	version, err := veClient.Version(ctx)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("version")

	err = d.Set(mkDataSourceVirtualEnvironmentVersionKeyboardLayout, version.Keyboard)
	diags = append(diags, diag.FromErr(err)...)
	err = d.Set(mkDataSourceVirtualEnvironmentVersionRelease, version.Release)
	diags = append(diags, diag.FromErr(err)...)
	err = d.Set(mkDataSourceVirtualEnvironmentVersionRepositoryID, version.RepositoryID)
	diags = append(diags, diag.FromErr(err)...)
	err = d.Set(mkDataSourceVirtualEnvironmentVersionVersion, version.Version)
	diags = append(diags, diag.FromErr(err)...)

	return diags
}
