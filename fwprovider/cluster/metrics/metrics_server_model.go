/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package metrics

import "github.com/hashicorp/terraform-plugin-framework/types"

type metricsServerResourceModel struct {
	ID                  types.String `tfsdk:"id"`
	Disable             types.Bool   `tfsdk:"disable"`
	MTU                 types.Int64  `tfsdk:"mtu"`
	Port                types.Int64  `tfsdk:"port"`
	Server              types.String `tfsdk:"server"`
	Timeout             types.Int64  `tfsdk:"timeout"`
	Type                types.String `tfsdk:"type"`
	InfluxAPIPathPrefix types.String `tfsdk:"influx_api_path_prefix"`
	InfluxBucket        types.String `tfsdk:"influx_bucket"`
	InfluxDBProto       types.String `tfsdk:"influx_db_proto"`
	InfluxMaxBodySize   types.String `tfsdk:"influx_max_body_size"`
	InfluxOrganization  types.String `tfsdk:"influx_organization"`
	InfluxToken         types.String `tfsdk:"influx_token"`
	InfluxVerify        types.Bool   `tfsdk:"influx_verify"`
	GraphitePath        types.String `tfsdk:"graphite_path"`
	GraphiteProto       types.String `tfsdk:"graphite_proto"`
}
