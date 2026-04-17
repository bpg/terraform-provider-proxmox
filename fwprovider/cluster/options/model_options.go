/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package options

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/bpg/terraform-provider-proxmox/fwprovider/attribute"
	customtypes "github.com/bpg/terraform-provider-proxmox/fwprovider/types"
	"github.com/bpg/terraform-provider-proxmox/proxmox/cluster"
)

const (
	// ClusterOptionsNextIDLowerMaximum is the maximum number for the "lower" range for the next VM ID option.
	// Note that this value is not documented in the section about the cluster options in the Proxmox VE API explorer but
	// [in the sections about QEMU (POST)] as well as [the dedicated Proxmox VE documentations about QEMU/KVM].
	//
	// [in the sections about QEMU (POST)]: https://pve.proxmox.com/pve-docs/api-viewer/#/nodes/{node}/qemu
	// [the dedicated Proxmox VE documentations about QEMU/KVM]: https://pve.proxmox.com/pve-docs/pve-admin-guide.html#_strong_qm_strong_qemu_kvm_virtual_machine_manager
	ClusterOptionsNextIDLowerMaximum = 999999999

	// ClusterOptionsNextIDLowerMinimum is the minimum number for the "lower" range for the next VM ID option.
	// Note that this value is not documented in the section about the cluster options in the Proxmox VE API explorer but
	// [in the sections about QEMU (POST)] as well as [the dedicated Proxmox VE documentations about QEMU/KVM].
	//
	// [in the sections about QEMU (POST)]: https://pve.proxmox.com/pve-docs/api-viewer/#/nodes/{node}/qemu
	// [the dedicated Proxmox VE documentations about QEMU/KVM]: https://pve.proxmox.com/pve-docs/pve-admin-guide.html#_strong_qm_strong_qemu_kvm_virtual_machine_manager
	ClusterOptionsNextIDLowerMinimum = 100
)

type clusterOptionsModel struct {
	ID                      types.String               `tfsdk:"id"`
	BandwidthLimitClone     types.Int64                `tfsdk:"bandwidth_limit_clone"`
	BandwidthLimitDefault   types.Int64                `tfsdk:"bandwidth_limit_default"`
	BandwidthLimitMigration types.Int64                `tfsdk:"bandwidth_limit_migration"`
	BandwidthLimitMove      types.Int64                `tfsdk:"bandwidth_limit_move"`
	BandwidthLimitRestore   types.Int64                `tfsdk:"bandwidth_limit_restore"`
	Console                 types.String               `tfsdk:"console"`
	CrsHA                   types.String               `tfsdk:"crs_ha"`
	CrsHARebalanceOnStart   types.Bool                 `tfsdk:"crs_ha_rebalance_on_start"`
	Description             types.String               `tfsdk:"description"`
	EmailFrom               types.String               `tfsdk:"email_from"`
	HAShutdownPolicy        types.String               `tfsdk:"ha_shutdown_policy"`
	HTTPProxy               types.String               `tfsdk:"http_proxy"`
	Keyboard                types.String               `tfsdk:"keyboard"`
	Language                types.String               `tfsdk:"language"`
	MacPrefix               types.String               `tfsdk:"mac_prefix"`
	MaxWorkers              types.Int64                `tfsdk:"max_workers"`
	MigrationNetwork        customtypes.IPCIDRValue    `tfsdk:"migration_cidr"`
	MigrationType           types.String               `tfsdk:"migration_type"`
	NextID                  *clusterOptionsNextIDModel `tfsdk:"next_id"`
	Notify                  *clusterOptionsNotifyModel `tfsdk:"notify"`
}

type clusterOptionsNextIDModel struct {
	Lower types.Int64 `tfsdk:"lower"`
	Upper types.Int64 `tfsdk:"upper"`
}

type clusterOptionsNotifyModel struct {
	HAFencingMode        types.String `tfsdk:"ha_fencing_mode"`
	HAFencingTarget      types.String `tfsdk:"ha_fencing_target"`
	PackageUpdates       types.String `tfsdk:"package_updates"`
	PackageUpdatesTarget types.String `tfsdk:"package_updates_target"`
	Replication          types.String `tfsdk:"replication"`
	ReplicationTarget    types.String `tfsdk:"replication_target"`
}

// haData returns HA settings parameter string for API, HA settings are
// defined, otherwise empty string is returned.
func (m *clusterOptionsModel) haData() string {
	var haDataParams []string

	if attribute.IsDefined(m.HAShutdownPolicy) {
		haDataParams = append(haDataParams, fmt.Sprintf("shutdown_policy=%s", m.HAShutdownPolicy.ValueString()))
	}

	if len(haDataParams) > 0 {
		return strings.Join(haDataParams, ",")
	}

	return ""
}

// migrationData returns migration settings parameter string for API, if any of migration
// settings are defined, otherwise empty string is returned.
func (m *clusterOptionsModel) migrationData() string {
	var migrationDataParams []string

	if attribute.IsDefined(m.MigrationType) {
		migrationDataParams = append(migrationDataParams, fmt.Sprintf("type=%s", m.MigrationType.ValueString()))
	}

	if attribute.IsDefined(m.MigrationNetwork) {
		migrationDataParams = append(migrationDataParams, fmt.Sprintf("network=%s", m.MigrationNetwork.ValueString()))
	}

	if len(migrationDataParams) > 0 {
		return strings.Join(migrationDataParams, ",")
	}

	return ""
}

// nextIDData returns settings for the "next-id" parameter string of the Proxmox VE API, if defined, otherwise an empty
// string is returned.
func (m *clusterOptionsModel) nextIDData() string {
	var nextIDDataParams []string

	if m.NextID == nil {
		return ""
	}

	if attribute.IsDefined(m.NextID.Lower) {
		nextIDDataParams = append(nextIDDataParams, fmt.Sprintf("lower=%d", m.NextID.Lower.ValueInt64()))
	}

	if attribute.IsDefined(m.NextID.Upper) {
		nextIDDataParams = append(nextIDDataParams, fmt.Sprintf("upper=%d", m.NextID.Upper.ValueInt64()))
	}

	if len(nextIDDataParams) > 0 {
		return strings.Join(nextIDDataParams, ",")
	}

	return ""
}

// notifyData returns settings for the "notify" parameter string of the Proxmox VE API, if defined, otherwise an empty
// string is returned.
func (m *clusterOptionsModel) notifyData() string {
	var notifyDataParams []string

	if m.Notify == nil {
		return ""
	}

	if attribute.IsDefined(m.Notify.HAFencingMode) {
		notifyDataParams = append(notifyDataParams, fmt.Sprintf("fencing=%s", m.Notify.HAFencingMode.ValueString()))
	}

	if attribute.IsDefined(m.Notify.HAFencingTarget) {
		notifyDataParams = append(
			notifyDataParams,
			fmt.Sprintf("target-fencing=%s", m.Notify.HAFencingTarget.ValueString()),
		)
	}

	if attribute.IsDefined(m.Notify.PackageUpdates) {
		notifyDataParams = append(
			notifyDataParams,
			fmt.Sprintf("package-updates=%s", m.Notify.PackageUpdates.ValueString()),
		)
	}

	if attribute.IsDefined(m.Notify.PackageUpdatesTarget) {
		notifyDataParams = append(
			notifyDataParams,
			fmt.Sprintf("target-package-updates=%s", m.Notify.PackageUpdatesTarget.ValueString()),
		)
	}

	if attribute.IsDefined(m.Notify.Replication) {
		notifyDataParams = append(notifyDataParams, fmt.Sprintf("replication=%s", m.Notify.Replication.ValueString()))
	}

	if attribute.IsDefined(m.Notify.ReplicationTarget) {
		notifyDataParams = append(
			notifyDataParams,
			fmt.Sprintf("target-replication=%s", m.Notify.ReplicationTarget.ValueString()),
		)
	}

	if len(notifyDataParams) > 0 {
		return strings.Join(notifyDataParams, ",")
	}

	return ""
}

// crsData returns cluster resource scheduling settings parameter string for API, if any of cluster resource scheduling
// settings are defined, otherwise empty string is returned.
func (m *clusterOptionsModel) crsData() string {
	var crsDataParams []string

	if attribute.IsDefined(m.CrsHA) {
		crsDataParams = append(crsDataParams, fmt.Sprintf("ha=%s", m.CrsHA.ValueString()))
	}

	if attribute.IsDefined(m.CrsHARebalanceOnStart) {
		var haRebalanceOnStart string
		if m.CrsHARebalanceOnStart.ValueBool() {
			haRebalanceOnStart = "1"
		} else {
			haRebalanceOnStart = "0"
		}

		crsDataParams = append(crsDataParams, fmt.Sprintf("ha-rebalance-on-start=%s", haRebalanceOnStart))
	}

	if len(crsDataParams) > 0 {
		return strings.Join(crsDataParams, ",")
	}

	return ""
}

// bandwidthData returns bandwidth limit settings parameter string for API, if any of bandwidth
// limit settings are defined, otherwise empty string is returned.
func (m *clusterOptionsModel) bandwidthData() string {
	var bandwidthParams []string

	if attribute.IsDefined(m.BandwidthLimitClone) {
		bandwidthParams = append(bandwidthParams, fmt.Sprintf("clone=%d", m.BandwidthLimitClone.ValueInt64()))
	}

	if attribute.IsDefined(m.BandwidthLimitDefault) {
		bandwidthParams = append(bandwidthParams, fmt.Sprintf("default=%d", m.BandwidthLimitDefault.ValueInt64()))
	}

	if attribute.IsDefined(m.BandwidthLimitMigration) {
		bandwidthParams = append(bandwidthParams, fmt.Sprintf("migration=%d", m.BandwidthLimitMigration.ValueInt64()))
	}

	if attribute.IsDefined(m.BandwidthLimitMove) {
		bandwidthParams = append(bandwidthParams, fmt.Sprintf("move=%d", m.BandwidthLimitMove.ValueInt64()))
	}

	if attribute.IsDefined(m.BandwidthLimitRestore) {
		bandwidthParams = append(bandwidthParams, fmt.Sprintf("restore=%d", m.BandwidthLimitRestore.ValueInt64()))
	}

	if len(bandwidthParams) > 0 {
		return strings.Join(bandwidthParams, ",")
	}

	return ""
}

// checkCompositeDelete adds an API field name to the delete list if the plan's serialized value
// is empty but the state's was not. This is the string-based equivalent of attribute.CheckDelete
// for composite fields where multiple model fields serialize into a single API parameter.
func checkCompositeDelete(planData, stateData string, toDelete *[]string, apiName string) {
	if planData == "" && stateData != "" {
		*toDelete = append(*toDelete, apiName)
	}
}

// toAPI converts the Terraform model to a cluster options API request body.
func (m *clusterOptionsModel) toAPI() *cluster.OptionsRequestData {
	body := &cluster.OptionsRequestData{}

	if !m.EmailFrom.IsUnknown() {
		body.EmailFrom = m.EmailFrom.ValueStringPointer()
	}

	if !m.Keyboard.IsUnknown() {
		body.Keyboard = m.Keyboard.ValueStringPointer()
	}

	if !m.Language.IsUnknown() {
		body.Language = m.Language.ValueStringPointer()
	}

	if !m.MaxWorkers.IsUnknown() {
		body.MaxWorkers = m.MaxWorkers.ValueInt64Pointer()
	}

	nextIDData := m.nextIDData()
	if nextIDData != "" {
		body.NextID = &nextIDData
	}

	notifyData := m.notifyData()
	if notifyData != "" {
		body.Notify = &notifyData
	}

	if !m.Console.IsUnknown() {
		body.Console = m.Console.ValueStringPointer()
	}

	if !m.HTTPProxy.IsUnknown() {
		body.HTTPProxy = m.HTTPProxy.ValueStringPointer()
	}

	if !m.MacPrefix.IsUnknown() {
		body.MacPrefix = m.MacPrefix.ValueStringPointer()
	}

	if !m.Description.IsUnknown() {
		body.Description = m.Description.ValueStringPointer()
	}

	haData := m.haData()
	if haData != "" {
		body.HASettings = &haData
	}

	bandwidthData := m.bandwidthData()
	if bandwidthData != "" {
		body.BandwidthLimit = &bandwidthData
	}

	crsData := m.crsData()
	if crsData != "" {
		body.ClusterResourceScheduling = &crsData
	}

	migrationData := m.migrationData()
	if migrationData != "" {
		body.Migration = &migrationData
	}

	return body
}

// fromAPI populates the Terraform model from a cluster options API response.
func (m *clusterOptionsModel) fromAPI(opts *cluster.OptionsResponseData) error {
	m.BandwidthLimitClone = types.Int64Null()
	m.BandwidthLimitDefault = types.Int64Null()
	m.BandwidthLimitMigration = types.Int64Null()
	m.BandwidthLimitMove = types.Int64Null()
	m.BandwidthLimitRestore = types.Int64Null()

	//nolint:nestif
	if opts.BandwidthLimit != nil {
		for bandwidth := range strings.SplitSeq(*opts.BandwidthLimit, ",") {
			bandwidthData := strings.SplitN(bandwidth, "=", 2)
			bandwidthName := bandwidthData[0]

			bandwidthLimit, err := strconv.ParseInt(bandwidthData[1], 10, 64)
			if err != nil {
				return fmt.Errorf("failed to parse bandwidth limit: %s", *opts.BandwidthLimit)
			}

			if bandwidthName == "clone" {
				m.BandwidthLimitClone = types.Int64Value(bandwidthLimit)
			}

			if bandwidthName == "default" {
				m.BandwidthLimitDefault = types.Int64Value(bandwidthLimit)
			}

			if bandwidthName == "migration" {
				m.BandwidthLimitMigration = types.Int64Value(bandwidthLimit)
			}

			if bandwidthName == "move" {
				m.BandwidthLimitMove = types.Int64Value(bandwidthLimit)
			}

			if bandwidthName == "restore" {
				m.BandwidthLimitRestore = types.Int64Value(bandwidthLimit)
			}
		}
	}

	m.EmailFrom = types.StringPointerValue(opts.EmailFrom)
	m.Keyboard = types.StringPointerValue(opts.Keyboard)
	m.Language = types.StringPointerValue(opts.Language)

	if opts.MaxWorkers != nil {
		value := int64(*opts.MaxWorkers)
		m.MaxWorkers = types.Int64PointerValue(&value)
	} else {
		m.MaxWorkers = types.Int64Null()
	}

	m.Console = types.StringPointerValue(opts.Console)
	m.HTTPProxy = types.StringPointerValue(opts.HTTPProxy)
	m.MacPrefix = types.StringPointerValue(opts.MacPrefix)

	if opts.Description != nil && *opts.Description != "" {
		m.Description = types.StringPointerValue(opts.Description)
	} else {
		m.Description = types.StringNull()
	}

	if opts.HASettings != nil {
		m.HAShutdownPolicy = types.StringPointerValue(opts.HASettings.ShutdownPolicy)
	} else {
		m.HAShutdownPolicy = types.StringNull()
	}

	if opts.Migration != nil {
		m.MigrationType = types.StringPointerValue(opts.Migration.Type)
		m.MigrationNetwork = customtypes.NewIPCIDRPointerValue(opts.Migration.Network)
	} else {
		m.MigrationType = types.StringNull()
		m.MigrationNetwork = customtypes.NewIPCIDRPointerValue(nil)
	}

	if opts.NextID != nil {
		m.NextID = &clusterOptionsNextIDModel{}
		m.NextID.Lower = types.Int64PointerValue(opts.NextID.Lower.PointerInt64())
		m.NextID.Upper = types.Int64PointerValue(opts.NextID.Upper.PointerInt64())
	}

	if opts.Notify != nil {
		m.Notify = &clusterOptionsNotifyModel{}
		m.Notify.HAFencingMode = types.StringPointerValue(opts.Notify.HAFencingMode)
		m.Notify.HAFencingTarget = types.StringPointerValue(opts.Notify.HAFencingTarget)
		m.Notify.PackageUpdates = types.StringPointerValue(opts.Notify.PackageUpdates)
		m.Notify.PackageUpdatesTarget = types.StringPointerValue(opts.Notify.PackageUpdatesTarget)
		m.Notify.Replication = types.StringPointerValue(opts.Notify.Replication)
		m.Notify.ReplicationTarget = types.StringPointerValue(opts.Notify.ReplicationTarget)
	}

	if opts.ClusterResourceScheduling != nil {
		m.CrsHARebalanceOnStart = types.BoolPointerValue(opts.ClusterResourceScheduling.HaRebalanceOnStart.PointerBool())
		m.CrsHA = types.StringPointerValue(opts.ClusterResourceScheduling.HA)
	} else {
		m.CrsHARebalanceOnStart = types.BoolNull()
		m.CrsHA = types.StringNull()
	}

	return nil
}
