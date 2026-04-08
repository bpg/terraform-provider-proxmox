/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package tasks

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTaskResult_OK(t *testing.T) {
	t.Parallel()

	r := TaskOK()
	require.NoError(t, r.Err())
	assert.False(t, r.HasWarnings())
	assert.Empty(t, r.Warnings())
}

func TestTaskResult_WithError(t *testing.T) {
	t.Parallel()

	r := TaskFailed(fmt.Errorf("task failed"))
	require.Error(t, r.Err())
	assert.Equal(t, "task failed", r.Err().Error())
	assert.False(t, r.HasWarnings())
}

func TestTaskResult_WithWarnings(t *testing.T) {
	t.Parallel()

	r := TaskOKWithWarnings([]string{
		"WARN: Systemd 258 detected. You may need to enable nesting.",
		"TASK WARNINGS: 1",
	})
	require.NoError(t, r.Err())
	assert.True(t, r.HasWarnings())
	assert.Len(t, r.Warnings(), 2)
	assert.Contains(t, r.Warnings()[0], "Systemd 258")
}

func TestTaskResult_WithErrorAndWarnings(t *testing.T) {
	t.Parallel()

	r := TaskFailedWithWarnings(
		fmt.Errorf("task failed"),
		[]string{"WARN: some warning"},
	)
	require.Error(t, r.Err())
	assert.True(t, r.HasWarnings())
	assert.Len(t, r.Warnings(), 1)
}

// mockDiags implements DiagnosticAccumulator for testing.
type mockDiags struct {
	errors   []string
	warnings []string
}

func (m *mockDiags) AddError(summary, detail string) {
	m.errors = append(m.errors, summary+": "+detail)
}

func (m *mockDiags) AddWarning(summary, detail string) {
	m.warnings = append(m.warnings, summary+": "+detail)
}

func TestTaskResult_AddDiags_OK(t *testing.T) {
	t.Parallel()

	d := &mockDiags{}
	hasErr := TaskOK().AddDiags(d, "op")
	assert.False(t, hasErr)
	assert.Empty(t, d.errors)
	assert.Empty(t, d.warnings)
}

func TestTaskResult_AddDiags_Error(t *testing.T) {
	t.Parallel()

	d := &mockDiags{}
	hasErr := TaskFailed(fmt.Errorf("boom")).AddDiags(d, "Create VM")
	assert.True(t, hasErr)
	assert.Equal(t, []string{"Create VM: boom"}, d.errors)
	assert.Empty(t, d.warnings)
}

func TestTaskResult_AddDiags_Warnings(t *testing.T) {
	t.Parallel()

	d := &mockDiags{}
	hasErr := TaskOKWithWarnings([]string{"w1", "w2"}).AddDiags(d, "Clone VM")
	assert.False(t, hasErr)
	assert.Empty(t, d.errors)
	assert.Equal(t, []string{"Clone VM: w1", "Clone VM: w2"}, d.warnings)
}

func TestTaskResult_AddDiags_ErrorAndWarnings(t *testing.T) {
	t.Parallel()

	d := &mockDiags{}
	hasErr := TaskFailedWithWarnings(fmt.Errorf("fail"), []string{"w1"}).AddDiags(d, "Start VM")
	assert.True(t, hasErr)
	assert.Equal(t, []string{"Start VM: fail"}, d.errors)
	assert.Equal(t, []string{"Start VM: w1"}, d.warnings)
}

func TestTaskResult_AddDiagsAsWarnings(t *testing.T) {
	t.Parallel()

	d := &mockDiags{}
	TaskFailedWithWarnings(fmt.Errorf("stop failed"), []string{"w1"}).AddDiagsAsWarnings(d, "VM stop")
	assert.Empty(t, d.errors, "errors should be empty — AddDiagsAsWarnings emits everything as warnings")
	assert.Equal(t, []string{"VM stop: stop failed", "VM stop: w1"}, d.warnings)
}
