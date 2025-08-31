/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package main

import (
	"context"
	"flag"
	"log"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov5"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6/tf6server"
	"github.com/hashicorp/terraform-plugin-mux/tf5to6server"
	"github.com/hashicorp/terraform-plugin-mux/tf6muxserver"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/bpg/terraform-provider-proxmox/fwprovider"
	"github.com/bpg/terraform-provider-proxmox/proxmoxtf/provider"
)

// Format Terraform code for use in documentation.
// If you do not have Terraform installed, you can remove the formatting command, but it is suggested
// to ensure the documentation is formatted properly.
//go:generate terraform fmt -recursive ./examples/
// Generate documentation.
//go:generate go tool github.com/hashicorp/terraform-plugin-docs/cmd/tfplugindocs generate --provider-dir ./ --rendered-website-dir ./build/docs-gen --provider-name "terraform-provider-proxmox" --rendered-provider-name "terraform-provider-proxmox" //nolint:lll

// Temporary: while migrating to the TF framework, we need to copy the generated docs to the right place
// for the resources / data sources that have been migrated.
//go:generate cp -R ./build/docs-gen/guides/. ./docs/guides/
//go:generate cp ./build/docs-gen/data-sources/virtual_environment_acme_account.md ./docs/data-sources/
//go:generate cp ./build/docs-gen/data-sources/virtual_environment_acme_accounts.md ./docs/data-sources/
//go:generate cp ./build/docs-gen/data-sources/virtual_environment_acme_plugin.md ./docs/data-sources/
//go:generate cp ./build/docs-gen/data-sources/virtual_environment_acme_plugins.md ./docs/data-sources/
//go:generate cp ./build/docs-gen/data-sources/virtual_environment_apt_repository.md ./docs/data-sources/
//go:generate cp ./build/docs-gen/data-sources/virtual_environment_apt_standard_repository.md ./docs/data-sources/
//go:generate cp ./build/docs-gen/data-sources/virtual_environment_datastores.md ./docs/data-sources/
//go:generate cp ./build/docs-gen/data-sources/virtual_environment_hagroup.md ./docs/data-sources/
//go:generate cp ./build/docs-gen/data-sources/virtual_environment_hagroups.md ./docs/data-sources/
//go:generate cp ./build/docs-gen/data-sources/virtual_environment_hardware_mapping_dir.md ./docs/data-sources/
//go:generate cp ./build/docs-gen/data-sources/virtual_environment_hardware_mapping_pci.md ./docs/data-sources/
//go:generate cp ./build/docs-gen/data-sources/virtual_environment_hardware_mapping_usb.md ./docs/data-sources/
//go:generate cp ./build/docs-gen/data-sources/virtual_environment_hardware_mappings.md ./docs/data-sources/
//go:generate cp ./build/docs-gen/data-sources/virtual_environment_haresource.md ./docs/data-sources/
//go:generate cp ./build/docs-gen/data-sources/virtual_environment_haresources.md ./docs/data-sources/
//go:generate cp ./build/docs-gen/data-sources/virtual_environment_sdn_zones.md ./docs/data-sources/
//go:generate cp ./build/docs-gen/data-sources/virtual_environment_sdn_zone_simple.md ./docs/data-sources/
//go:generate cp ./build/docs-gen/data-sources/virtual_environment_sdn_zone_vlan.md ./docs/data-sources/
//go:generate cp ./build/docs-gen/data-sources/virtual_environment_sdn_zone_qinq.md ./docs/data-sources/
//go:generate cp ./build/docs-gen/data-sources/virtual_environment_sdn_zone_vxlan.md ./docs/data-sources/
//go:generate cp ./build/docs-gen/data-sources/virtual_environment_sdn_zone_evpn.md ./docs/data-sources/
//go:generate cp ./build/docs-gen/data-sources/virtual_environment_version.md ./docs/data-sources/
//go:generate cp ./build/docs-gen/data-sources/virtual_environment_vm2.md ./docs/data-sources/
//go:generate cp ./build/docs-gen/data-sources/virtual_environment_metrics_server.md ./docs/data-sources/
//go:generate cp ./build/docs-gen/resources/virtual_environment_acl.md ./docs/resources/
//go:generate cp ./build/docs-gen/resources/virtual_environment_acme_account.md ./docs/resources/
//go:generate cp ./build/docs-gen/resources/virtual_environment_acme_dns_plugin.md ./docs/resources/
//go:generate cp ./build/docs-gen/resources/virtual_environment_apt_repository.md ./docs/resources/
//go:generate cp ./build/docs-gen/resources/virtual_environment_apt_standard_repository.md ./docs/resources/
//go:generate cp ./build/docs-gen/resources/virtual_environment_cluster_options.md ./docs/resources/
//go:generate cp ./build/docs-gen/resources/virtual_environment_download_file.md ./docs/resources/
//go:generate cp ./build/docs-gen/resources/virtual_environment_hagroup.md ./docs/resources/
//go:generate cp ./build/docs-gen/resources/virtual_environment_hardware_mapping_dir.md ./docs/resources/
//go:generate cp ./build/docs-gen/resources/virtual_environment_hardware_mapping_pci.md ./docs/resources/
//go:generate cp ./build/docs-gen/resources/virtual_environment_hardware_mapping_usb.md ./docs/resources/
//go:generate cp ./build/docs-gen/resources/virtual_environment_haresource.md ./docs/resources/
//go:generate cp ./build/docs-gen/resources/virtual_environment_network_linux_bridge.md ./docs/resources/
//go:generate cp ./build/docs-gen/resources/virtual_environment_network_linux_vlan.md ./docs/resources/
//go:generate cp ./build/docs-gen/resources/virtual_environment_sdn_applier.md ./docs/resources/
//go:generate cp ./build/docs-gen/resources/virtual_environment_sdn_zone_simple.md ./docs/resources/
//go:generate cp ./build/docs-gen/resources/virtual_environment_sdn_zone_vlan.md ./docs/resources/
//go:generate cp ./build/docs-gen/resources/virtual_environment_sdn_zone_qinq.md ./docs/resources/
//go:generate cp ./build/docs-gen/resources/virtual_environment_sdn_zone_vxlan.md ./docs/resources/
//go:generate cp ./build/docs-gen/resources/virtual_environment_sdn_zone_evpn.md ./docs/resources/
//go:generate cp ./build/docs-gen/resources/virtual_environment_storage_directory.md ./docs/resources/
//go:generate cp ./build/docs-gen/resources/virtual_environment_storage_lvmthin.md ./docs/resources/
//go:generate cp ./build/docs-gen/resources/virtual_environment_storage_lvm.md ./docs/resources/
//go:generate cp ./build/docs-gen/resources/virtual_environment_storage_nfs.md ./docs/resources/
//go:generate cp ./build/docs-gen/resources/virtual_environment_storage_pbs.md ./docs/resources/
//go:generate cp ./build/docs-gen/resources/virtual_environment_storage_smb.md ./docs/resources/
//go:generate cp ./build/docs-gen/resources/virtual_environment_storage_zfspool.md ./docs/resources/
//go:generate cp ./build/docs-gen/resources/virtual_environment_user_token.md ./docs/resources/
//go:generate cp ./build/docs-gen/resources/virtual_environment_vm2.md ./docs/resources/
//go:generate cp ./build/docs-gen/resources/virtual_environment_metrics_server.md ./docs/resources/

// these will be set by the goreleaser configuration
// to appropriate values for the compiled binary.
var version = "dev"

func main() {
	ctx := context.Background()

	var debug bool

	flag.BoolVar(&debug, "debug", false, "set to true to run the provider with support for debuggers like delve")
	flag.Parse()

	upgradedSdkServer, err := tf5to6server.UpgradeServer(
		ctx,
		func() tfprotov5.ProviderServer {
			return schema.NewGRPCProviderServer(
				provider.ProxmoxVirtualEnvironment(),
			)
		},
	)
	if err != nil {
		log.Fatal(err)
	}

	providers := []func() tfprotov6.ProviderServer{
		providerserver.NewProtocol6(fwprovider.New(version)()),
		func() tfprotov6.ProviderServer {
			return upgradedSdkServer
		},
	}

	muxServer, err := tf6muxserver.NewMuxServer(ctx, providers...)
	if err != nil {
		log.Fatal(err)
	}

	var serveOpts []tf6server.ServeOpt

	if debug {
		serveOpts = append(serveOpts, tf6server.WithManagedDebug())
	}

	// Remove any date and time prefix in log package function output to
	// prevent duplicate timestamp and incorrect log level setting
	log.SetFlags(log.Flags() &^ (log.Ldate | log.Ltime))

	err = tf6server.Serve(
		"registry.terraform.io/bpg/proxmox",
		muxServer.ProviderServer,
		serveOpts...,
	)
	if err != nil {
		log.Fatal(err)
	}
}
