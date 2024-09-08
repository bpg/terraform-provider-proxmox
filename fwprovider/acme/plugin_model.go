/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package acme

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// acmePluginsModel maps the schema data for the ACME plugins data source.
type acmePluginsModel struct {
	Plugins []acmePluginModel `tfsdk:"plugins"`
}

type baseACMEPluginModel struct {
	// API plugin name
	API types.String `tfsdk:"api"`
	// DNS plugin data
	Data types.Map `tfsdk:"data"`
	// Prevent changes if current configuration file has a different digest.
	// This can be used to prevent concurrent modifications.
	Digest types.String `tfsdk:"digest"`
	// Plugin ID name
	Plugin types.String `tfsdk:"plugin"`
	// Extra delay in seconds to wait before requesting validation (0 - 172800)
	ValidationDelay types.Int64 `tfsdk:"validation_delay"`
}

// acmePluginModel maps the schema data for an ACME plugin.
type acmePluginModel struct {
	baseACMEPluginModel
	Type types.String `tfsdk:"type"`
}

// acmePluginCreateModel maps the schema data for an ACME plugin.
type acmePluginCreateModel struct {
	baseACMEPluginModel
	// Flag to disable the config
	Disable types.Bool `tfsdk:"disable"`
}
