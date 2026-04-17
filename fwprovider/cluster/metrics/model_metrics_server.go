/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package metrics

import (
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/bpg/terraform-provider-proxmox/fwprovider/attribute"
	"github.com/bpg/terraform-provider-proxmox/proxmox/cluster/metrics"
	proxmoxtypes "github.com/bpg/terraform-provider-proxmox/proxmox/types"
)

type metricsServerModel struct {
	ID                     types.String `tfsdk:"id"`
	Name                   types.String `tfsdk:"name"`
	Disable                types.Bool   `tfsdk:"disable"`
	MTU                    types.Int64  `tfsdk:"mtu"`
	Port                   types.Int64  `tfsdk:"port"`
	Server                 types.String `tfsdk:"server"`
	Timeout                types.Int64  `tfsdk:"timeout"`
	Type                   types.String `tfsdk:"type"`
	InfluxAPIPathPrefix    types.String `tfsdk:"influx_api_path_prefix"`
	InfluxBucket           types.String `tfsdk:"influx_bucket"`
	InfluxDBProto          types.String `tfsdk:"influx_db_proto"`
	InfluxMaxBodySize      types.Int64  `tfsdk:"influx_max_body_size"`
	InfluxOrganization     types.String `tfsdk:"influx_organization"`
	InfluxToken            types.String `tfsdk:"influx_token"`
	InfluxVerify           types.Bool   `tfsdk:"influx_verify"`
	GraphitePath           types.String `tfsdk:"graphite_path"`
	GraphiteProto          types.String `tfsdk:"graphite_proto"`
	OTelProto              types.String `tfsdk:"opentelemetry_proto"`
	OTelPath               types.String `tfsdk:"opentelemetry_path"`
	OTelTimeout            types.Int64  `tfsdk:"opentelemetry_timeout"`
	OTelHeaders            types.String `tfsdk:"opentelemetry_headers"`
	OTelVerifySSL          types.Bool   `tfsdk:"opentelemetry_verify_ssl"`
	OTelMaxBodySize        types.Int64  `tfsdk:"opentelemetry_max_body_size"`
	OTelResourceAttributes types.String `tfsdk:"opentelemetry_resource_attributes"`
	OTelCompression        types.String `tfsdk:"opentelemetry_compression"`
}

// fromAPI populates the model from a metrics server PVE API response.
// The API response does not contain the name, so it must be passed directly.
func (m *metricsServerModel) fromAPI(name string, data *metrics.ServerData) {
	m.ID = types.StringValue(name)
	m.Name = types.StringValue(name)

	// `disable` has a schema Default of false; the API omits it when false, so
	// normalize here. `verify-certificate` and `otel-verify-ssl` are type-specific
	// and the caller preserves plan/state values; map API-provided values through
	// and leave null otherwise so the caller's override path stays straightforward.
	m.Disable = boolOrDefault(data.Disable, false)
	m.InfluxVerify = types.BoolPointerValue(data.Verify.PointerBool())
	m.OTelVerifySSL = types.BoolPointerValue(data.OTelVerifySSL.PointerBool())

	m.MTU = types.Int64PointerValue(data.MTU)
	m.Port = types.Int64Value(data.Port)
	m.Server = types.StringValue(data.Server)
	m.Timeout = types.Int64PointerValue(data.Timeout)
	m.Type = types.StringPointerValue(data.Type)
	m.InfluxAPIPathPrefix = types.StringPointerValue(data.APIPathPrefix)
	m.InfluxBucket = types.StringPointerValue(data.Bucket)
	m.InfluxDBProto = types.StringPointerValue(data.InfluxDBProto)
	m.InfluxMaxBodySize = types.Int64PointerValue(data.MaxBodySize)
	m.InfluxOrganization = types.StringPointerValue(data.Organization)
	m.InfluxToken = types.StringPointerValue(data.Token)
	m.GraphitePath = types.StringPointerValue(data.Path)
	m.GraphiteProto = types.StringPointerValue(data.Proto)
	m.OTelProto = types.StringPointerValue(data.OTelProto)
	m.OTelPath = types.StringPointerValue(data.OTelPath)
	m.OTelTimeout = types.Int64PointerValue(data.OTelTimeout)
	m.OTelHeaders = types.StringPointerValue(data.OTelHeaders)
	m.OTelMaxBodySize = types.Int64PointerValue(data.OTelMaxBodySize)
	m.OTelResourceAttributes = types.StringPointerValue(data.OTelResourceAttributes)
	m.OTelCompression = types.StringPointerValue(data.OTelCompression)
}

// boolOrDefault returns the value pointed to by b, falling back to def when nil.
// Used for API fields that the server omits when they equal the default.
func boolOrDefault(b *proxmoxtypes.CustomBool, def bool) types.Bool {
	if v := b.PointerBool(); v != nil {
		return types.BoolValue(*v)
	}

	return types.BoolValue(def)
}

// preserveTypeSpecificBools copies InfluxVerify and OTelVerifySSL from src onto dst
// when dst has null values. PVE omits these fields from GET responses when they equal
// the server default, so the caller's plan (Create/Update) or prior state (Read) is
// the authoritative source. A universal schema Default would leak the value into
// toAPI for non-matching server types, which PVE rejects.
func preserveTypeSpecificBools(dst, src *metricsServerModel) {
	if dst.InfluxVerify.IsNull() {
		dst.InfluxVerify = src.InfluxVerify
	}

	if dst.OTelVerifySSL.IsNull() {
		dst.OTelVerifySSL = src.OTelVerifySSL
	}
}

// toAPI converts the Terraform model to a metrics server request body used for both POST and PUT.
func (m *metricsServerModel) toAPI() *metrics.ServerRequestData {
	data := &metrics.ServerRequestData{}

	data.ID = m.Name.ValueString()

	data.Disable = attribute.CustomBoolPtrFromValue(m.Disable)
	data.MTU = attribute.Int64PtrFromValue(m.MTU)
	data.Port = m.Port.ValueInt64()
	data.Server = m.Server.ValueString()
	data.Timeout = attribute.Int64PtrFromValue(m.Timeout)
	data.Type = attribute.StringPtrFromValue(m.Type)
	data.APIPathPrefix = attribute.StringPtrFromValue(m.InfluxAPIPathPrefix)
	data.Bucket = attribute.StringPtrFromValue(m.InfluxBucket)
	data.InfluxDBProto = attribute.StringPtrFromValue(m.InfluxDBProto)
	data.MaxBodySize = attribute.Int64PtrFromValue(m.InfluxMaxBodySize)
	data.Organization = attribute.StringPtrFromValue(m.InfluxOrganization)
	data.Token = attribute.StringPtrFromValue(m.InfluxToken)
	data.Verify = attribute.CustomBoolPtrFromValue(m.InfluxVerify)
	data.Path = attribute.StringPtrFromValue(m.GraphitePath)
	data.Proto = attribute.StringPtrFromValue(m.GraphiteProto)
	data.OTelProto = attribute.StringPtrFromValue(m.OTelProto)
	data.OTelPath = attribute.StringPtrFromValue(m.OTelPath)
	data.OTelTimeout = attribute.Int64PtrFromValue(m.OTelTimeout)
	data.OTelHeaders = attribute.StringPtrFromValue(m.OTelHeaders)
	data.OTelVerifySSL = attribute.CustomBoolPtrFromValue(m.OTelVerifySSL)
	data.OTelMaxBodySize = attribute.Int64PtrFromValue(m.OTelMaxBodySize)
	data.OTelResourceAttributes = attribute.StringPtrFromValue(m.OTelResourceAttributes)
	data.OTelCompression = attribute.StringPtrFromValue(m.OTelCompression)

	return data
}

type metricsServerDatasourceModel struct {
	ID                     types.String `tfsdk:"id"`
	Name                   types.String `tfsdk:"name"`
	Disable                types.Bool   `tfsdk:"disable"`
	Port                   types.Int64  `tfsdk:"port"`
	Server                 types.String `tfsdk:"server"`
	Type                   types.String `tfsdk:"type"`
	OTelProto              types.String `tfsdk:"opentelemetry_proto"`
	OTelPath               types.String `tfsdk:"opentelemetry_path"`
	OTelTimeout            types.Int64  `tfsdk:"opentelemetry_timeout"`
	OTelHeaders            types.String `tfsdk:"opentelemetry_headers"`
	OTelVerifySSL          types.Bool   `tfsdk:"opentelemetry_verify_ssl"`
	OTelMaxBodySize        types.Int64  `tfsdk:"opentelemetry_max_body_size"`
	OTelResourceAttributes types.String `tfsdk:"opentelemetry_resource_attributes"`
	OTelCompression        types.String `tfsdk:"opentelemetry_compression"`
}

// fromAPI populates the datasource model from a metrics server PVE API response.
// The API response does not contain the name, so it must be passed directly.
func (m *metricsServerDatasourceModel) fromAPI(name string, data *metrics.ServerData) {
	m.ID = types.StringValue(name)
	m.Name = types.StringValue(name)

	// `disable` is omitted when false; the rest pass through directly.
	m.Disable = boolOrDefault(data.Disable, false)
	m.OTelVerifySSL = types.BoolPointerValue(data.OTelVerifySSL.PointerBool())

	m.Port = types.Int64Value(data.Port)
	m.Server = types.StringValue(data.Server)
	m.Type = types.StringPointerValue(data.Type)
	m.OTelProto = types.StringPointerValue(data.OTelProto)
	m.OTelPath = types.StringPointerValue(data.OTelPath)
	m.OTelTimeout = types.Int64PointerValue(data.OTelTimeout)
	m.OTelHeaders = types.StringPointerValue(data.OTelHeaders)
	m.OTelMaxBodySize = types.Int64PointerValue(data.OTelMaxBodySize)
	m.OTelResourceAttributes = types.StringPointerValue(data.OTelResourceAttributes)
	m.OTelCompression = types.StringPointerValue(data.OTelCompression)
}
