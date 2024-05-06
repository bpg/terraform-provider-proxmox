//go:build tools
// +build tools

/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package tools

// Manage tool dependencies via go.mod.
//
// https://github.com/golang/go/wiki/Modules#how-can-i-track-tool-dependencies-for-a-module
// https://github.com/golang/go/issues/25922
import (
	_ "github.com/golangci/golangci-lint/cmd/golangci-lint"
	_ "github.com/goreleaser/goreleaser"
	_ "github.com/hashicorp/terraform-plugin-docs/cmd/tfplugindocs"
)

// Format Terraform code for use in documentation.
// If you do not have Terraform installed, you can remove the formatting command, but it is suggested
// to ensure the documentation is formatted properly.
//go:generate terraform fmt -recursive ../examples/
// Generate documentation.
//go:generate go run github.com/hashicorp/terraform-plugin-docs/cmd/tfplugindocs generate --provider-dir ../ --rendered-website-dir ./build/docs-gen

// Temporary: while migrating to the TF framework, we need to copy the generated docs to the right place
// for the resources / data sources that have been migrated.
//go:generate cp -R ../build/docs-gen/guides/ ../docs/guides/
//go:generate cp ../build/docs-gen/data-sources/virtual_environment_version.md ../docs/data-sources/
//go:generate cp ../build/docs-gen/data-sources/virtual_environment_hagroup.md ../docs/data-sources/
//go:generate cp ../build/docs-gen/data-sources/virtual_environment_hagroups.md ../docs/data-sources/
//go:generate cp ../build/docs-gen/data-sources/virtual_environment_hardware_mapping_pci.md ../docs/data-sources/
//go:generate cp ../build/docs-gen/data-sources/virtual_environment_hardware_mapping_usb.md ../docs/data-sources/
//go:generate cp ../build/docs-gen/data-sources/virtual_environment_hardware_mappings.md ../docs/data-sources/
//go:generate cp ../build/docs-gen/data-sources/virtual_environment_haresource.md ../docs/data-sources/
//go:generate cp ../build/docs-gen/data-sources/virtual_environment_haresources.md ../docs/data-sources/
//go:generate cp ../build/docs-gen/resources/virtual_environment_network_linux_bridge.md ../docs/resources/
//go:generate cp ../build/docs-gen/resources/virtual_environment_network_linux_vlan.md ../docs/resources/
//go:generate cp ../build/docs-gen/resources/virtual_environment_hagroup.md ../docs/resources/
//go:generate cp ../build/docs-gen/resources/virtual_environment_hardware_mapping_pci.md ../docs/resources/
//go:generate cp ../build/docs-gen/resources/virtual_environment_hardware_mapping_usb.md ../docs/resources/
//go:generate cp ../build/docs-gen/resources/virtual_environment_haresource.md ../docs/resources/
//go:generate cp ../build/docs-gen/resources/virtual_environment_cluster_options.md ../docs/resources/
//go:generate cp ../build/docs-gen/resources/virtual_environment_download_file.md ../docs/resources/
//go:generate cp ../build/docs-gen/resources/virtual_environment_user_token.md ../docs/resources/
//go:generate cp ../build/docs-gen/resources/virtual_environment_vm2.md ../docs/resources/
