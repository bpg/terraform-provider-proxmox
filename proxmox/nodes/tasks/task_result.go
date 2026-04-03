/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package tasks

// TaskResult holds the outcome of a Proxmox task, including any warnings
// from the task log. This allows callers to surface warnings as Terraform
// diagnostics rather than silently ignoring them.
type TaskResult struct {
	warnings []string
	err      error
}

// TaskOK returns a successful result with no warnings.
func TaskOK() *TaskResult {
	return &TaskResult{}
}

// TaskFailed returns a result with an error and no warnings.
func TaskFailed(err error) *TaskResult {
	return &TaskResult{err: err}
}

// TaskOKWithWarnings returns a successful result that carries warning lines
// from the task log.
func TaskOKWithWarnings(warnings []string) *TaskResult {
	return &TaskResult{warnings: warnings}
}

// TaskFailedWithWarnings returns a result with both an error and warning lines.
func TaskFailedWithWarnings(err error, warnings []string) *TaskResult {
	return &TaskResult{err: err, warnings: warnings}
}

// Err returns the error if the task failed, or nil on success.
func (r *TaskResult) Err() error {
	return r.err
}

// HasWarnings returns true if the task produced warning lines.
func (r *TaskResult) HasWarnings() bool {
	return len(r.warnings) > 0
}

// Warnings returns the warning lines from the task log.
func (r *TaskResult) Warnings() []string {
	return r.warnings
}
