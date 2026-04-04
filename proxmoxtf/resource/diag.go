/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package resource

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"

	"github.com/bpg/terraform-provider-proxmox/proxmox/nodes/tasks"
)

// sdkDiagAccumulator adapts SDK v2 diag.Diagnostics to satisfy tasks.DiagnosticAccumulator.
type sdkDiagAccumulator struct {
	diags *diag.Diagnostics
}

func (s *sdkDiagAccumulator) AddError(summary, detail string) {
	*s.diags = append(*s.diags, diag.Diagnostic{Severity: diag.Error, Summary: summary, Detail: detail})
}

func (s *sdkDiagAccumulator) AddWarning(summary, detail string) {
	*s.diags = append(*s.diags, diag.Diagnostic{Severity: diag.Warning, Summary: summary, Detail: detail})
}

// TaskResultDiags converts a TaskResult into SDK diagnostics containing any errors and warnings.
func TaskResultDiags(result tasks.TaskResult, summary string) diag.Diagnostics {
	var diags diag.Diagnostics
	result.AddDiags(&sdkDiagAccumulator{&diags}, summary)

	return diags
}
