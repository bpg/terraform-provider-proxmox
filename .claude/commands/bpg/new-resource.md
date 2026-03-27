---
name: new-resource
description: Scaffold a new Framework resource with all 3 files and provider registration
argument-hint: <package/ResourceName> e.g. cluster/firewall/AliasSet
allowed-tools:
  - Read
  - Write
  - Edit
  - Bash
  - Grep
  - Glob
  - AskUserQuestion
---

<objective>
Scaffold a new Terraform Plugin Framework resource following all project conventions.

Creates 3 files and registers the resource in the provider:

1. `resource_{name}.go` — CRUD operations with short-name alias (ADR-007)
2. `model_{name}.go` — Model struct with `toAPI()` / `fromAPI()` methods
3. `resource_{name}_test.go` — Acceptance test skeleton

**Argument format:** `<package_path/ResourceName>`

Examples:

- `cluster/firewall/AliasSet` → `fwprovider/cluster/firewall/resource_alias_set.go`
- `network/Bridge` → `fwprovider/network/resource_bridge.go`
- `access/Token` → `fwprovider/access/resource_token.go`
</objective>

<context>
Input: $ARGUMENTS
</context>

<process>

## Step 1: Parse Arguments

Parse `$ARGUMENTS` into:

- **Package path**: directory under `fwprovider/` (e.g., `cluster/firewall`)
- **Resource name (PascalCase)**: the resource struct name (e.g., `AliasSet`)
- **Resource name (snake_case)**: for file names and TypeName (e.g., `alias_set`)
- **Short TypeName**: `proxmox_{snake_case}` (e.g., `proxmox_alias_set`) — per ADR-007

If `$ARGUMENTS` is empty or unclear, ask:

```text
AskUserQuestion(
  header: "Resource Details",
  question: "What resource should I scaffold? Provide package path and name.",
  options: [
    { label: "Enter path/Name", description: "e.g. cluster/firewall/AliasSet" }
  ]
)
```

## Step 2: Gather Context

Ask the user for the API client path and key attributes:

```text
AskUserQuestion(
  header: "API Client",
  question: "What API client should this resource use?",
  options: [
    { label: "Enter client path", description: "e.g. cfg.Client.Cluster().Firewall().Aliases()" }
  ]
)
```

Then ask:

```text
AskUserQuestion(
  header: "Key Attributes",
  question: "List the key attributes for this resource (comma-separated). I'll create schema stubs for each.",
  options: [
    { label: "Enter attributes", description: "e.g. name (required, string), cidr (required, string), comment (optional, string)" }
  ]
)
```

## Step 3: Check Package Exists

```bash
ls fwprovider/{package_path}/ 2>/dev/null || echo "NEW_PACKAGE"
```

If the package doesn't exist, create it. If it does, check for existing resources to understand the package's conventions (single vs multi-resource package).

## Step 4: Determine File Naming

- **New package with single resource**: `resource.go`, `model.go`, `resource_test.go`
- **Existing package or multi-resource**: `resource_{snake_name}.go`, `model_{snake_name}.go`, `resource_{snake_name}_test.go`

## Step 5: Create Model File

Write `model_{snake_name}.go` (or `model.go` for single-resource packages):

```go
/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package {package_name}

import (
	"github.com/hashicorp/terraform-plugin-framework/types"

	// TODO: Import the API client package for request/response types
)

// {camelCase}Model maps the resource schema to Go types.
type {camelCase}Model struct {
	ID types.String `tfsdk:"id"`
	// TODO: Add fields matching schema attributes
}

// toAPI converts the Terraform model to an API request body.
func (m *{camelCase}Model) toAPI() *{api_package}.CreateRequestBody {
	data := &{api_package}.CreateRequestBody{}

	// TODO: Map model fields to API request fields
	// Required fields use .ValueString(), .ValueInt64(), etc.
	// Optional fields use .ValueStringPointer(), .ValueInt64Pointer(), etc.

	return data
}

// fromAPI populates the model from an API response.
func (m *{camelCase}Model) fromAPI(id string, data *{api_package}.ResponseData) {
	m.ID = types.StringValue(id)

	// TODO: Map API response fields to model fields
	// Use types.StringValue() for required, types.StringPointerValue() for optional
}
```

## Step 6: Create Resource File

Write `resource_{snake_name}.go` (or `resource.go`):

```go
/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package {package_name}

import (
	"context"
	"errors"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"

	"github.com/bpg/terraform-provider-proxmox/fwprovider/attribute"
	"github.com/bpg/terraform-provider-proxmox/fwprovider/config"
	"github.com/bpg/terraform-provider-proxmox/fwprovider/migration"
	"github.com/bpg/terraform-provider-proxmox/proxmox/api"
	// TODO: Import the API client package
)

// Compile-time interface checks.
var (
	_ resource.Resource                = &{camelCase}Resource{}
	_ resource.ResourceWithConfigure   = &{camelCase}Resource{}
	_ resource.ResourceWithImportState = &{camelCase}Resource{}
)

type {camelCase}Resource struct {
	client *{api_client_type}
}

// New{PascalCase}Resource creates a new resource instance.
func New{PascalCase}Resource() resource.Resource {
	return &{camelCase}Resource{}
}

func (r *{camelCase}Resource) Metadata(
	_ context.Context,
	req resource.MetadataRequest,
	resp *resource.MetadataResponse,
) {
	resp.TypeName = req.ProviderTypeName + "_{snake_case}"
}

func (r *{camelCase}Resource) Configure(
	_ context.Context,
	req resource.ConfigureRequest,
	resp *resource.ConfigureResponse,
) {
	if req.ProviderData == nil {
		return
	}

	cfg, ok := req.ProviderData.(config.Resource)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected config.Resource, got: %T", req.ProviderData),
		)

		return
	}

	r.client = {client_access} // e.g. cfg.Client.Cluster().Firewall().Aliases()
}

func (r *{camelCase}Resource) Schema(
	_ context.Context,
	_ resource.SchemaRequest,
	resp *resource.SchemaResponse,
) {
	resp.Schema = schema.Schema{
		Description:        "Manages a {human_readable_name}.",
		DeprecationMessage: migration.DeprecationMessage("proxmox_{snake_case}"),
		Attributes: map[string]schema.Attribute{
			"id": attribute.ResourceID(),
			// TODO: Add resource-specific attributes
		},
	}
}

func (r *{camelCase}Resource) Create(
	ctx context.Context,
	req resource.CreateRequest,
	resp *resource.CreateResponse,
) {
	var plan {camelCase}Model

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	reqData := plan.toAPI()

	err := r.client.Create(ctx, reqData)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Create Resource",
			err.Error(),
		)

		return
	}

	// TODO: Set plan.ID from the created resource identifier
	// plan.ID = plan.Name

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *{camelCase}Resource) Read(
	ctx context.Context,
	req resource.ReadRequest,
	resp *resource.ReadResponse,
) {
	var state {camelCase}Model

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	data, err := r.client.Get(ctx, state.ID.ValueString())
	if err != nil {
		if errors.Is(err, api.ErrResourceDoesNotExist) {
			resp.State.RemoveResource(ctx)
			return
		}

		resp.Diagnostics.AddError(
			"Unable to Read Resource",
			err.Error(),
		)

		return
	}

	readModel := &{camelCase}Model{}
	readModel.fromAPI(state.ID.ValueString(), data)

	resp.Diagnostics.Append(resp.State.Set(ctx, readModel)...)
}

func (r *{camelCase}Resource) Update(
	ctx context.Context,
	req resource.UpdateRequest,
	resp *resource.UpdateResponse,
) {
	var plan, state {camelCase}Model

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	reqData := plan.toAPI()

	// TODO: If the API supports field deletion for optional attributes:
	// var toDelete []string
	// attribute.CheckDelete(plan.Field, state.Field, &toDelete, "field")
	// reqData.Delete = toDelete

	err := r.client.Update(ctx, state.ID.ValueString(), reqData)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Update Resource",
			err.Error(),
		)

		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *{camelCase}Resource) Delete(
	ctx context.Context,
	req resource.DeleteRequest,
	resp *resource.DeleteResponse,
) {
	var state {camelCase}Model

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.Delete(ctx, state.ID.ValueString())
	if err != nil {
		if errors.Is(err, api.ErrResourceDoesNotExist) {
			return
		}

		resp.Diagnostics.AddError(
			"Unable to Delete Resource",
			err.Error(),
		)

		return
	}
}

func (r *{camelCase}Resource) ImportState(
	ctx context.Context,
	req resource.ImportStateRequest,
	resp *resource.ImportStateResponse,
) {
	data, err := r.client.Get(ctx, req.ID)
	if err != nil {
		if errors.Is(err, api.ErrResourceDoesNotExist) {
			resp.Diagnostics.AddError(
				"Resource Does Not Exist",
				fmt.Sprintf("Resource with ID %q does not exist.", req.ID),
			)

			return
		}

		resp.Diagnostics.AddError(
			"Unable to Import Resource",
			err.Error(),
		)

		return
	}

	readModel := &{camelCase}Model{}
	readModel.fromAPI(req.ID, data)

	resp.Diagnostics.Append(resp.State.Set(ctx, readModel)...)
}

// --- Short-name alias (ADR-007 Phase 2) ---

var (
	_ resource.Resource                = &{camelCase}ResourceShort{}
	_ resource.ResourceWithConfigure   = &{camelCase}ResourceShort{}
	_ resource.ResourceWithImportState = &{camelCase}ResourceShort{}
	_ resource.ResourceWithMoveState   = &{camelCase}ResourceShort{}
)

type {camelCase}ResourceShort struct {
	{camelCase}Resource
}

// New{PascalCase}ShortResource creates the short-name alias.
func New{PascalCase}ShortResource() resource.Resource {
	return &{camelCase}ResourceShort{}
}

func (r *{camelCase}ResourceShort) Metadata(
	_ context.Context,
	_ resource.MetadataRequest,
	resp *resource.MetadataResponse,
) {
	resp.TypeName = "proxmox_{snake_case}"
}

func (r *{camelCase}ResourceShort) Schema(
	ctx context.Context,
	req resource.SchemaRequest,
	resp *resource.SchemaResponse,
) {
	r.{camelCase}Resource.Schema(ctx, req, resp)
	resp.Schema.DeprecationMessage = ""
}

func (r *{camelCase}ResourceShort) MoveState(ctx context.Context) []resource.StateMover {
	schemaResp := &resource.SchemaResponse{}
	r.{camelCase}Resource.Schema(ctx, resource.SchemaRequest{}, schemaResp)

	return []resource.StateMover{
		migration.PrefixMoveState(
			"proxmox_virtual_environment_{snake_case}",
			&schemaResp.Schema,
		),
	}
}
```

## Step 7: Create Test File

Write `resource_{snake_name}_test.go` (or `resource_test.go`):

```go
//go:build acceptance || all

/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package {package_name}_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/bpg/terraform-provider-proxmox/fwprovider/test"
)

func TestAccResource{PascalCase}(t *testing.T) {
	te := test.InitEnvironment(t)

	tests := []struct {
		name  string
		steps []resource.TestStep
	}{
		{"create and read back", []resource.TestStep{
			{
				Config: te.RenderConfig(`
				resource "proxmox_{snake_case}" "test" {
					# TODO: Add required attributes
				}`),
				Check: resource.ComposeTestCheckFunc(
					test.ResourceAttributes("proxmox_{snake_case}.test", map[string]string{
						// TODO: Add expected attribute values
					}),
				),
			},
		}},
		{"create and import", []resource.TestStep{
			{
				Config: te.RenderConfig(`
				resource "proxmox_{snake_case}" "test" {
					# TODO: Add required attributes
				}`),
			},
			{
				ResourceName:      "proxmox_{snake_case}.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resource.ParallelTest(t, resource.TestCase{
				ProtoV6ProviderFactories: te.AccProviders,
				Steps:                    tt.steps,
			})
		})
	}
}
```

## Step 8: Register in Provider

Edit `fwprovider/provider.go` to add both constructors to the `Resources()` method:

```go
{package_name}.New{PascalCase}Resource,
{package_name}.New{PascalCase}ShortResource, // proxmox_{snake_case}
```

Add the import if needed:

```go
"github.com/bpg/terraform-provider-proxmox/fwprovider/{package_path}"
```

Insert alphabetically by package name in both the import block and the Resources() slice.

## Step 9: Verify Build

```bash
make build
```

If the build fails due to missing API client types (expected for a scaffold), note the TODOs.

## Step 10: Summary

Report what was created:

```text
=== Resource Scaffold Created ===

Files:
  fwprovider/{package_path}/resource_{snake_name}.go  — CRUD + short-name alias
  fwprovider/{package_path}/model_{snake_name}.go     — Model with toAPI/fromAPI stubs
  fwprovider/{package_path}/resource_{snake_name}_test.go — Acceptance test skeleton

Registered in: fwprovider/provider.go
  - New{PascalCase}Resource       (proxmox_virtual_environment_{snake_case})
  - New{PascalCase}ShortResource  (proxmox_{snake_case})

TODOs remaining:
  1. Implement API client package (proxmox/{api_path}/)
  2. Fill in model fields and toAPI/fromAPI mappings
  3. Complete schema attributes with validators
  4. Fill in acceptance test configurations
  5. Add documentation template if needed (templates/resources/{snake_case}.md.tmpl)
```

</process>

<success_criteria>

- [ ] Package path and resource name parsed correctly
- [ ] Model file created with struct, toAPI(), fromAPI()
- [ ] Resource file created with all CRUD methods
- [ ] Short-name alias with MoveState (ADR-007)
- [ ] DeprecationMessage on long-name schema
- [ ] Test file created with build tag and test skeleton
- [ ] Resource registered in provider.go (both long and short names)
- [ ] `attribute.ResourceID()` used for the `id` attribute
- [ ] `config.Resource` used in Configure method
- [ ] `api.ErrResourceDoesNotExist` handled in Read (remove from state) and Delete (ignore)
- [ ] MPL-2.0 license header on all files
- [ ] Build attempted to verify no syntax errors

</success_criteria>

<tips>

- Check existing resources in the same package for naming patterns (single vs multi-resource)
- The API client package typically mirrors the resource path: `fwprovider/cluster/ha/` → `proxmox/cluster/ha/`
- For comma-separated API values, use list/set attributes with join/split in toAPI/fromAPI (ADR-004)
- Retry patterns go in the CRUD methods, not the model (ADR-005)
- If the resource needs `attribute.CheckDelete` in Update, also import `"github.com/bpg/terraform-provider-proxmox/fwprovider/attribute"`

</tips>
