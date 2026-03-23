/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package migration

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
)

// PrefixMoveState returns a StateMover that migrates state from an old resource
// type name to the current one. Use for same-schema renames where only the
// resource type name prefix changes (ADR-007 Phase 2).
//
// sourceSchema must be provided so the Framework can decode SourceRawState into
// SourceState. Since the schema is identical between old and new names, we copy
// the decoded state directly to the target.
//
// This helper does not check SourceSchemaVersion or SourceProviderAddress — both
// are assumed to match since this is a same-provider, same-schema rename.
func PrefixMoveState(oldTypeName string, sourceSchema schema.Schema) resource.StateMover {
	return resource.StateMover{
		SourceSchema: &sourceSchema,
		StateMover: func(_ context.Context, req resource.MoveStateRequest, resp *resource.MoveStateResponse) {
			if req.SourceTypeName != oldTypeName {
				return
			}

			// SourceState is populated by the Framework because we provided SourceSchema.
			// Schema is identical — copy the decoded state to the target.
			resp.TargetState = tfsdk.State{
				Schema: resp.TargetState.Schema,
				Raw:    req.SourceState.Raw,
			}
		},
	}
}

// DeprecationMessage returns the standard deprecation message for a resource
// being renamed from its old proxmox_virtual_environment_* name to proxmox_*.
func DeprecationMessage(newTypeName string) string {
	return fmt.Sprintf("Use %s instead. This resource will be removed in v1.0.", newTypeName)
}
