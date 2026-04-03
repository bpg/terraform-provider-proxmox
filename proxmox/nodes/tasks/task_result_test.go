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
