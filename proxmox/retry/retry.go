package retry

import (
	"context"
	"errors"
	"net/http"
	"strings"
	"time"

	retrylib "github.com/avast/retry-go/v5"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"github.com/bpg/terraform-provider-proxmox/proxmox/api"
)

// Operation configures retry behavior for a Proxmox API operation.
type Operation struct {
	name      string
	attempts  uint
	baseDelay time.Duration
	retryIf   func(error) bool
	// Category-specific behavior, set by preset constructors.
	logOnRetry     bool
	lastErrorOnly  bool
	useBackoff     bool
	isAlreadyDone  func(error) bool
	untilSucceeded bool
}

// Option configures an Operation.
type Option func(*Operation)

// WithAttempts sets the maximum number of retry attempts.
func WithAttempts(n uint) Option {
	return func(o *Operation) { o.attempts = n }
}

// WithBaseDelay sets the base delay between retry attempts.
func WithBaseDelay(d time.Duration) Option {
	return func(o *Operation) { o.baseDelay = d }
}

// WithRetryIf sets the predicate that determines whether an error is retryable.
func WithRetryIf(fn func(error) bool) Option {
	return func(o *Operation) { o.retryIf = fn }
}

// WithAlreadyDoneCheck sets a predicate that detects when a retry attempt finds
// evidence that the first attempt already succeeded (e.g. "already exists").
// Only checked on retry attempts, not the first attempt.
func WithAlreadyDoneCheck(fn func(error) bool) Option {
	return func(o *Operation) { o.isAlreadyDone = fn }
}

// IsTransientAPIError returns true for HTTP 5xx, "got no worker upid",
// and "got timeout" errors.
func IsTransientAPIError(err error) bool {
	if err == nil {
		return false
	}

	var httpErr *api.HTTPError
	if errors.As(err, &httpErr) {
		return httpErr.Code >= http.StatusInternalServerError
	}

	return strings.Contains(err.Error(), "got no worker upid") ||
		strings.Contains(err.Error(), "got timeout")
}

// ErrorContains returns a predicate that checks if the error message contains substr.
func ErrorContains(substr string) func(error) bool {
	return func(err error) bool {
		if err == nil {
			return false
		}

		return strings.Contains(err.Error(), substr)
	}
}

// NewTaskOperation creates an Operation for async UPID-based task operations
// (create, clone, delete, start, resize). Defaults: exponential backoff,
// retry logging, LastErrorOnly(false), WaitForTask errors are unrecoverable,
// 3 attempts.
func NewTaskOperation(name string, opts ...Option) *Operation {
	o := &Operation{
		name:          name,
		attempts:      3,
		baseDelay:     1 * time.Second,
		logOnRetry:    true,
		lastErrorOnly: false,
		useBackoff:    true,
	}

	for _, opt := range opts {
		opt(o)
	}

	return o
}

// DoTask executes an async task operation with retry. The dispatchFn is retried
// according to the operation's retry configuration. Errors from waitFn are
// wrapped in retry.Unrecoverable to prevent re-dispatch after a wait failure.
// The isAlreadyDone check (if configured) is only applied on retry attempts.
func (o *Operation) DoTask(
	ctx context.Context,
	dispatchFn func() (*string, error),
	waitFn func(ctx context.Context, taskID string) error,
) error {
	retrying := false

	retryOpts := []retrylib.Option{
		retrylib.Context(ctx),
		retrylib.Attempts(o.attempts),
		retrylib.Delay(o.baseDelay),
		retrylib.LastErrorOnly(o.lastErrorOnly),
	}

	if o.useBackoff {
		retryOpts = append(retryOpts, retrylib.DelayType(retrylib.BackOffDelay))
	}

	if o.retryIf != nil {
		retryOpts = append(retryOpts, retrylib.RetryIf(o.retryIf))
	}

	retryOpts = append(retryOpts, retrylib.OnRetry(func(n uint, err error) {
		retrying = true

		if o.logOnRetry {
			tflog.Warn(ctx, "retrying "+o.name, map[string]any{
				"attempt": n,
				"error":   err.Error(),
			})
		}
	}))

	//nolint:wrapcheck // errors from dispatchFn/waitFn pass through intentionally
	return retrylib.New(retryOpts...).Do(func() error {
		taskID, err := dispatchFn()
		if err != nil {
			if retrying && o.isAlreadyDone != nil && o.isAlreadyDone(err) {
				return nil
			}

			return err
		}

		// nil taskID with nil error means the operation is already done
		// (e.g., "already running") — skip the wait.
		if taskID == nil {
			return nil
		}

		if err := waitFn(ctx, *taskID); err != nil {
			return retrylib.Unrecoverable(err)
		}

		return nil
	})
}

// NewAPICallOperation creates an Operation for synchronous API calls that block
// server-side (e.g. PUT /config). Defaults: exponential backoff, retry logging,
// LastErrorOnly(false), 3 attempts, 1s base delay.
func NewAPICallOperation(name string, opts ...Option) *Operation {
	o := &Operation{
		name:          name,
		attempts:      3,
		baseDelay:     1 * time.Second,
		logOnRetry:    true,
		lastErrorOnly: false,
		useBackoff:    true,
	}

	for _, opt := range opts {
		opt(o)
	}

	return o
}

// Do executes a synchronous API call with retry.
func (o *Operation) Do(ctx context.Context, fn func() error) error {
	retryOpts := []retrylib.Option{
		retrylib.Context(ctx),
		retrylib.Attempts(o.attempts),
		retrylib.Delay(o.baseDelay),
		retrylib.LastErrorOnly(o.lastErrorOnly),
	}

	if o.useBackoff {
		retryOpts = append(retryOpts, retrylib.DelayType(retrylib.BackOffDelay))
	}

	if o.retryIf != nil {
		retryOpts = append(retryOpts, retrylib.RetryIf(o.retryIf))
	}

	if o.logOnRetry {
		retryOpts = append(retryOpts, retrylib.OnRetry(func(n uint, err error) {
			tflog.Warn(ctx, "retrying "+o.name, map[string]any{
				"attempt": n,
				"error":   err.Error(),
			})
		}))
	}

	//nolint:wrapcheck // errors from fn pass through intentionally
	return retrylib.New(retryOpts...).Do(fn)
}

// NewPollOperation creates an Operation for wait-for-condition polling loops
// (e.g. WaitForVMStatus, WaitForConfigUnlock). Defaults: fixed 1s delay,
// no attempt limit (relies on context deadline), LastErrorOnly(true),
// no retry logging.
func NewPollOperation(name string, opts ...Option) *Operation {
	o := &Operation{
		name:           name,
		baseDelay:      1 * time.Second,
		lastErrorOnly:  true,
		untilSucceeded: true,
	}

	for _, opt := range opts {
		opt(o)
	}

	return o
}

// DoPoll executes a polling loop with retry until fn returns nil or a
// non-retryable error. Uses fixed delay and no attempt limit by default.
func (o *Operation) DoPoll(ctx context.Context, fn func() error) error {
	retryOpts := []retrylib.Option{
		retrylib.Context(ctx),
		retrylib.Delay(o.baseDelay),
		retrylib.DelayType(retrylib.FixedDelay),
		retrylib.LastErrorOnly(o.lastErrorOnly),
	}

	if o.untilSucceeded {
		retryOpts = append(retryOpts, retrylib.UntilSucceeded())
	} else {
		retryOpts = append(retryOpts, retrylib.Attempts(o.attempts))
	}

	if o.retryIf != nil {
		retryOpts = append(retryOpts, retrylib.RetryIf(o.retryIf))
	}

	//nolint:wrapcheck // errors from fn pass through intentionally
	return retrylib.New(retryOpts...).Do(fn)
}
