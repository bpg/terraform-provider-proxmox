/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package metrics

import (
	"github.com/bpg/terraform-provider-proxmox/proxmox/cluster/metrics"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type metricsServerModel struct {
	ID                  types.String `tfsdk:"id"`
	Name                types.String `tfsdk:"name"`
	Disable             types.Bool   `tfsdk:"disable"`
	MTU                 types.Int64  `tfsdk:"mtu"`
	Port                types.Int64  `tfsdk:"port"`
	Server              types.String `tfsdk:"server"`
	Timeout             types.Int64  `tfsdk:"timeout"`
	Type                types.String `tfsdk:"type"`
	InfluxAPIPathPrefix types.String `tfsdk:"influx_api_path_prefix"`
	InfluxBucket        types.String `tfsdk:"influx_bucket"`
	InfluxDBProto       types.String `tfsdk:"influx_db_proto"`
	InfluxMaxBodySize   types.Int64  `tfsdk:"influx_max_body_size"`
	InfluxOrganization  types.String `tfsdk:"influx_organization"`
	InfluxToken         types.String `tfsdk:"influx_token"`
	InfluxVerify        types.Bool   `tfsdk:"influx_verify"`
	GraphitePath        types.String `tfsdk:"graphite_path"`
	GraphiteProto       types.String `tfsdk:"graphite_proto"`
}

// importFromAPI takes data from metrics server PVE API response and set fields based on it.
// Note: API response does not contain name so it must be passed directly.
func (m *metricsServerModel) importFromAPI(name string, data *metrics.ServerData) {
	m.ID = types.StringValue(name)
	m.Name = types.StringValue(name)

	var disable *bool

	if data.Disable != nil {
		if *data.Disable == 1 {
			*disable = true
		} else {
			*disable = false
		}
	}

	m.Disable = types.BoolPointerValue(disable)

	m.MTU = types.Int64PointerValue(data.MTU)
	m.Port = types.Int64PointerValue(data.Port)
	m.Server = types.StringPointerValue(data.Server)
	m.Timeout = types.Int64PointerValue(data.Timeout)
	m.Type = types.StringPointerValue(data.Type)
	m.InfluxAPIPathPrefix = types.StringPointerValue(data.APIPathPrefix)
	m.InfluxBucket = types.StringPointerValue(data.Bucket)
	m.InfluxDBProto = types.StringPointerValue(data.InfluxDBProto)
	m.InfluxMaxBodySize = types.Int64PointerValue(data.MaxBodySize)
	m.InfluxOrganization = types.StringPointerValue(data.Organization)
	m.InfluxToken = types.StringPointerValue(data.Token)

	var influxVerify *bool

	if data.Verify != nil {
		if *data.Verify == 1 {
			*influxVerify = true
		} else {
			*influxVerify = false
		}
	}

	m.InfluxVerify = types.BoolPointerValue(influxVerify)

	m.GraphitePath = types.StringPointerValue(data.Path)
	m.GraphiteProto = types.StringPointerValue(data.Proto)
}

// toAPIRequestBody creates metrics server request data for PUT and POST requests.
func (m *metricsServerModel) toAPIRequestBody() *metrics.ServerRequestData {
	data := &metrics.ServerRequestData{}

	data.ID = m.Name.ValueStringPointer()

	if !m.Disable.IsNull() {
		var disable *int64
		if m.Disable.ValueBool() {
			*disable = 1
		} else {
			*disable = 0
		}

		data.Disable = disable
	}

	data.MTU = m.MTU.ValueInt64Pointer()
	data.Port = m.Port.ValueInt64Pointer()
	data.Server = m.Server.ValueStringPointer()
	data.Timeout = m.Timeout.ValueInt64Pointer()
	data.Type = m.Type.ValueStringPointer()
	data.APIPathPrefix = m.InfluxAPIPathPrefix.ValueStringPointer()
	data.Bucket = m.InfluxBucket.ValueStringPointer()
	data.InfluxDBProto = m.InfluxDBProto.ValueStringPointer()
	data.MaxBodySize = m.InfluxMaxBodySize.ValueInt64Pointer()
	data.Organization = m.InfluxOrganization.ValueStringPointer()
	data.Token = m.InfluxToken.ValueStringPointer()

	if !m.InfluxVerify.IsNull() {
		var influxVerify *int64
		if m.InfluxVerify.ValueBool() {
			*influxVerify = 1
		} else {
			*influxVerify = 0
		}

		data.Verify = influxVerify
	}

	data.Path = m.GraphitePath.ValueStringPointer()
	data.Proto = m.GraphiteProto.ValueStringPointer()

	return data
}

type metricsServerDatasourceModel struct {
	ID      types.String `tfsdk:"id"`
	Name    types.String `tfsdk:"name"`
	Disable types.Bool   `tfsdk:"disable"`
	Port    types.Int64  `tfsdk:"port"`
	Server  types.String `tfsdk:"server"`
	Type    types.String `tfsdk:"type"`
}

// importFromAPI takes data from metrics server PVE API response and set fields based on it.
// Note: API response does not contain name so it must be passed directly.
func (m *metricsServerDatasourceModel) importFromAPI(name string, data *metrics.ServerData) {
	m.ID = types.StringValue(name)
	m.Name = types.StringValue(name)

	var disable *bool

	if data.Disable != nil {
		if *data.Disable == 1 {
			*disable = true
		} else {
			*disable = false
		}
	}

	m.Disable = types.BoolPointerValue(disable)
	m.Port = types.Int64PointerValue(data.Port)
	m.Server = types.StringPointerValue(data.Server)
	m.Type = types.StringPointerValue(data.Type)
}
