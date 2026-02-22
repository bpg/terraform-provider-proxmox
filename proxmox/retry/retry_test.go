package retry

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/bpg/terraform-provider-proxmox/proxmox/api"
)

func TestIsTransientAPIError(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		err  error
		want bool
	}{
		{"nil", nil, false},
		{"HTTP 500", &api.HTTPError{Code: http.StatusInternalServerError}, true},
		{"HTTP 503", &api.HTTPError{Code: http.StatusServiceUnavailable}, true},
		{"HTTP 400", &api.HTTPError{Code: http.StatusBadRequest}, false},
		{"HTTP 403", &api.HTTPError{Code: http.StatusForbidden}, false},
		{"HTTP 404", &api.HTTPError{Code: http.StatusNotFound}, false},
		{"got no worker upid", fmt.Errorf("got no worker upid"), true},
		{"got timeout", fmt.Errorf("got timeout"), true},
		{"wrapped got no worker upid", fmt.Errorf("error: %w", fmt.Errorf("got no worker upid")), true},
		{"wrapped got timeout", fmt.Errorf("error: %w", fmt.Errorf("got timeout")), true},
		{"generic error", errors.New("something else"), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.want, IsTransientAPIError(tt.err))
		})
	}
}

func TestErrorContains(t *testing.T) {
	t.Parallel()

	check := ErrorContains("already exists")
	assert.True(t, check(fmt.Errorf("container 100 already exists")))
	assert.False(t, check(fmt.Errorf("something else")))
	assert.False(t, check(nil))
}

func TestDoTask_Success(t *testing.T) {
	t.Parallel()

	op := NewTaskOperation("test-op",
		WithRetryIf(IsTransientAPIError),
	)

	dispatched := 0
	waited := 0

	err := op.DoTask(t.Context(),
		func() (*string, error) {
			dispatched++
			id := "UPID:test:1234"

			return &id, nil
		},
		func(_ context.Context, taskID string) error {
			waited++

			assert.Equal(t, "UPID:test:1234", taskID)

			return nil
		},
	)

	require.NoError(t, err)
	assert.Equal(t, 1, dispatched)
	assert.Equal(t, 1, waited)
}

func TestDoTask_RetryOnTransientError(t *testing.T) {
	t.Parallel()

	op := NewTaskOperation("test-op",
		WithBaseDelay(1*time.Millisecond),
		WithRetryIf(IsTransientAPIError),
	)

	dispatched := 0

	err := op.DoTask(t.Context(),
		func() (*string, error) {
			dispatched++
			if dispatched == 1 {
				return nil, &api.HTTPError{Code: http.StatusInternalServerError}
			}

			id := "UPID:test:1234"

			return &id, nil
		},
		func(_ context.Context, _ string) error { return nil },
	)

	require.NoError(t, err)
	assert.Equal(t, 2, dispatched)
}

func TestDoTask_AlreadyDoneOnRetry(t *testing.T) {
	t.Parallel()

	op := NewTaskOperation("test-op",
		WithBaseDelay(1*time.Millisecond),
		WithRetryIf(IsTransientAPIError),
		WithAlreadyDoneCheck(ErrorContains("already exists")),
	)

	dispatched := 0

	err := op.DoTask(t.Context(),
		func() (*string, error) {
			dispatched++
			if dispatched == 1 {
				return nil, &api.HTTPError{Code: http.StatusInternalServerError}
			}
			// Second attempt: "already exists" because first actually succeeded
			return nil, fmt.Errorf("container 100 already exists")
		},
		func(_ context.Context, _ string) error {
			t.Fatal("waitFn should not be called when already done")
			return nil
		},
	)

	require.NoError(t, err)
	assert.Equal(t, 2, dispatched)
}

func TestDoTask_AlreadyDoneIgnoredOnFirstAttempt(t *testing.T) {
	t.Parallel()

	op := NewTaskOperation("test-op",
		WithBaseDelay(1*time.Millisecond),
		WithRetryIf(IsTransientAPIError),
		WithAlreadyDoneCheck(ErrorContains("already exists")),
	)

	err := op.DoTask(t.Context(),
		func() (*string, error) {
			return nil, fmt.Errorf("container 100 already exists")
		},
		func(_ context.Context, _ string) error { return nil },
	)

	// "already exists" on first attempt is a real error, not suppressed
	require.Error(t, err)
	assert.Contains(t, err.Error(), "already exists")
}

func TestDoTask_NilTaskIDSkipsWait(t *testing.T) {
	t.Parallel()

	op := NewTaskOperation("test-op",
		WithBaseDelay(1*time.Millisecond),
		WithRetryIf(IsTransientAPIError),
	)

	err := op.DoTask(t.Context(),
		func() (*string, error) {
			// nil taskID, nil error = "already done" (e.g., already running)
			return nil, nil //nolint:nilnil // testing defensive nil-taskID handling
		},
		func(_ context.Context, _ string) error {
			t.Fatal("waitFn should not be called when taskID is nil")
			return nil
		},
	)

	require.NoError(t, err)
}

func TestDoTask_WaitErrorIsUnrecoverable(t *testing.T) {
	t.Parallel()

	op := NewTaskOperation("test-op",
		WithBaseDelay(1*time.Millisecond),
		WithRetryIf(IsTransientAPIError),
	)

	dispatched := 0

	err := op.DoTask(t.Context(),
		func() (*string, error) {
			dispatched++
			id := "UPID:test:1234"

			return &id, nil
		},
		func(_ context.Context, _ string) error {
			return fmt.Errorf("task failed with exit code: ERROR")
		},
	)

	require.Error(t, err)
	// Dispatch should only happen once — waitFn error is unrecoverable
	assert.Equal(t, 1, dispatched)
}

func TestDoTask_NoRetryOnNonTransient(t *testing.T) {
	t.Parallel()

	op := NewTaskOperation("test-op",
		WithBaseDelay(1*time.Millisecond),
		WithRetryIf(IsTransientAPIError),
	)

	dispatched := 0

	err := op.DoTask(t.Context(),
		func() (*string, error) {
			dispatched++
			return nil, &api.HTTPError{Code: http.StatusBadRequest, Message: "bad request"}
		},
		func(_ context.Context, _ string) error { return nil },
	)

	require.Error(t, err)
	assert.Equal(t, 1, dispatched)
}

func TestDo_Success(t *testing.T) {
	t.Parallel()

	op := NewAPICallOperation("test-op",
		WithRetryIf(ErrorContains("got timeout")),
	)

	called := 0

	err := op.Do(t.Context(), func() error {
		called++
		return nil
	})

	require.NoError(t, err)
	assert.Equal(t, 1, called)
}

func TestDo_RetryOnTimeout(t *testing.T) {
	t.Parallel()

	op := NewAPICallOperation("test-op",
		WithBaseDelay(1*time.Millisecond),
		WithRetryIf(ErrorContains("got timeout")),
	)

	called := 0

	err := op.Do(t.Context(), func() error {
		called++
		if called == 1 {
			return fmt.Errorf("got timeout")
		}

		return nil
	})

	require.NoError(t, err)
	assert.Equal(t, 2, called)
}

func TestDo_NoRetryOnNonMatchingError(t *testing.T) {
	t.Parallel()

	op := NewAPICallOperation("test-op",
		WithBaseDelay(1*time.Millisecond),
		WithRetryIf(ErrorContains("got timeout")),
	)

	called := 0

	err := op.Do(t.Context(), func() error {
		called++
		return fmt.Errorf("permission denied")
	})

	require.Error(t, err)
	assert.Equal(t, 1, called)
}

func TestDoPoll_Success(t *testing.T) {
	t.Parallel()

	stillWaiting := errors.New("still waiting")
	op := NewPollOperation("test-poll",
		WithRetryIf(func(err error) bool { return errors.Is(err, stillWaiting) }),
	)

	called := 0

	err := op.DoPoll(t.Context(), func() error {
		called++
		if called < 3 {
			return stillWaiting
		}

		return nil
	})

	require.NoError(t, err)
	assert.Equal(t, 3, called)
}

func TestDoPoll_ContextTimeout(t *testing.T) {
	t.Parallel()

	stillWaiting := errors.New("still waiting")
	op := NewPollOperation("test-poll",
		WithBaseDelay(1*time.Millisecond),
		WithRetryIf(func(err error) bool { return errors.Is(err, stillWaiting) }),
	)

	ctx, cancel := context.WithTimeout(t.Context(), 50*time.Millisecond)
	defer cancel()

	err := op.DoPoll(ctx, func() error {
		return stillWaiting
	})

	require.Error(t, err)
	assert.ErrorIs(t, err, context.DeadlineExceeded)
}

func TestDoPoll_NonRetryableError(t *testing.T) {
	t.Parallel()

	stillWaiting := errors.New("still waiting")
	op := NewPollOperation("test-poll",
		WithRetryIf(func(err error) bool { return errors.Is(err, stillWaiting) }),
	)

	called := 0

	err := op.DoPoll(t.Context(), func() error {
		called++
		return fmt.Errorf("unexpected API error")
	})

	require.Error(t, err)
	assert.Equal(t, 1, called)
}
