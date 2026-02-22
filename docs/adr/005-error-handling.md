# ADR-005: Error Handling Patterns

## Status

Accepted

## Date

2026-02-04 (retroactive documentation)

## Context

The provider has a three-layer architecture (HTTP client → domain client → resource) and each layer handles errors differently. Without documented conventions, error message formats diverge across resources — the codebase currently mixes "Unable to", "Error", "Failed to", "Could not", and "Cannot" prefixes. Contributors need a single pattern to follow.

Additionally, retry behavior for transient failures, async task polling, and resource state waiting needs clear guidance.

## Decision

### Error Message Format

User-facing diagnostic summaries (the first argument to `resp.Diagnostics.AddError`) **must** use the format:

```text
"Unable to [Action] [Resource]"
```

Where `[Action]` is a CRUD verb and `[Resource]` identifies the Proxmox object.

> **Note:** Legacy resources still use other prefixes ("Error", "Failed to", "Could not"). All **new code** must use the `"Unable to"` format. Legacy resources will be migrated over time.

Standard summaries for the resource lifecycle:

| Operation                      | Summary Format                               | Example                                          |
|--------------------------------|----------------------------------------------|--------------------------------------------------|
| Create                         | `"Unable to Create [Resource]"`              | `"Unable to Create SDN VNet"`                    |
| Read-back after create         | `"Unable to Read [Resource] After Creation"` | `"Unable to Read SDN VNet After Creation"`       |
| Read                           | `"Unable to Read [Resource]"`                | `"Unable to Read SDN VNet"`                      |
| Update                         | `"Unable to Update [Resource]"`              | `"Unable to Update SDN VNet"`                    |
| Read-back after update         | `"Unable to Read [Resource] After Update"`   | `"Unable to Read SDN VNet After Update"`         |
| Delete                         | `"Unable to Delete [Resource]"`              | `"Unable to Delete SDN VNet"`                    |
| Import                         | `"Unable to Import [Resource]"`              | `"Unable to Import SDN VNet"`                    |
| Not found (import/datasource)  | `"[Resource] Not Found"`                     | `"SDN VNet Not Found"`                           |

The detail string (second argument) should be `err.Error()`, which carries the full error chain from the API layer.

### Three-Layer Error Architecture

Errors flow through three layers. Each layer adds context while preserving the original error for inspection.

#### Layer 1: HTTP Client (`proxmox/api/`)

The base API client (`client.go`) handles HTTP-level concerns:

- Returns `HTTPError{Code, Message}` for non-2xx responses
- Parses the Proxmox error response body to extract field-level errors and messages
- Joins `ErrResourceDoesNotExist` sentinel with the `HTTPError` for 404 responses and 500 responses containing "does not exist"

```go
// Sentinel errors
const ErrNoDataObjectInResponse Error = "the server did not include a data object in the response"
const ErrResourceDoesNotExist Error = "the requested resource does not exist"

// For not-found cases, both errors are joined so callers can check either:
// This happens for HTTP 404 AND for HTTP 500 responses containing "does not exist":
errors.Join(ErrResourceDoesNotExist, &HTTPError{Code: 404, Message: "..."})
errors.Join(ErrResourceDoesNotExist, &HTTPError{Code: 500, Message: "...does not exist..."})
```

#### Layer 2: Domain Client (`proxmox/{domain}/`)

Domain client methods wrap errors with operational context using `fmt.Errorf` and the `%w` verb to preserve the error chain:

```go
func (c *Client) GetZone(ctx context.Context, id string) (*ZoneData, error) {
    resBody := &struct{ Data *ZoneData `json:"data"` }{}
    err := c.DoRequest(ctx, http.MethodGet, c.ExpandPath(id), nil, resBody)
    if err != nil {
        return nil, fmt.Errorf("error reading SDN zone %s: %w", id, err)
    }

    if resBody.Data == nil {
        return nil, api.ErrNoDataObjectInResponse
    }

    return resBody.Data, nil
}
```

#### Layer 3: Resource (`fwprovider/`)

Resources check sentinel errors to determine behavior, then convert remaining errors to Terraform diagnostics:

**Read — remove from state when resource is gone:**

```go
if errors.Is(err, api.ErrResourceDoesNotExist) {
    resp.State.RemoveResource(ctx)
    return
}
resp.Diagnostics.AddError("Unable to Read SDN VNet", err.Error())
```

**Delete — ignore already-gone resources:**

```go
err := r.client.SDNVnets(state.ID.ValueString()).DeleteVnet(ctx)
if err != nil && !errors.Is(err, api.ErrResourceDoesNotExist) {
    resp.Diagnostics.AddError("Unable to Delete SDN VNet", err.Error())
}
```

**Import/Datasource — report not-found as an error:**

```go
if errors.Is(err, api.ErrResourceDoesNotExist) {
    resp.Diagnostics.AddError("SDN VNet Not Found",
        fmt.Sprintf("SDN VNet with ID '%s' was not found", req.ID))
    return
}
```

### Error Inspection

Use standard Go error inspection:

| Need                 | Method                                                       |
|----------------------|--------------------------------------------------------------|
| Check for sentinel   | `errors.Is(err, api.ErrResourceDoesNotExist)`                |
| Extract HTTP status  | `errors.As(err, &httpError)` then inspect `httpError.Code`   |
| Add context          | `fmt.Errorf("context: %w", err)`                             |

### Retry Policies

Retry logic is centralized in the `proxmox/retry/` package, which provides three operation constructors:

| Constructor            | Use Case                                              | Attempts  | Backoff     | Method   |
|------------------------|-------------------------------------------------------|-----------|-------------|----------|
| `NewTaskOperation`     | Async UPID-based tasks (create, clone, delete, start) | 3         | Exponential | `DoTask` |
| `NewAPICallOperation`  | Synchronous blocking API calls (e.g. PUT /config)     | 3         | Exponential | `Do`     |
| `NewPollOperation`     | Wait-for-condition loops (status, config unlock)      | Unlimited | Fixed (1s)  | `DoPoll` |

Key behaviors of `DoTask`:

- `dispatchFn` is retried; `waitFn` errors are wrapped in `Unrecoverable` (no re-dispatch after wait failure)
- `nil` taskID with `nil` error means "already done" (e.g., "already running") — skips the wait
- `WithAlreadyDoneCheck` is only applied on retry attempts, not the first attempt

Common retry predicates:

| Predicate                    | Matches                                           |
|------------------------------|---------------------------------------------------|
| `IsTransientAPIError`        | HTTP 5xx, "got no worker upid", "got timeout"     |
| `ErrorContains(substr)`      | Error message contains substring                  |

Do not retry on:

- 4xx client errors (except specific known-transient messages)
- `ErrResourceDoesNotExist` (unless polling for creation)
- Authentication failures (401, 403)

> **Delete predicate trap:** `ErrResourceDoesNotExist` can arrive via HTTP 500 (not just 404) because `proxmox/api/client.go` uses `errors.Join(ErrResourceDoesNotExist, httpError)` for 500 responses containing "does not exist". This means `IsTransientAPIError` alone will match these errors (it checks for 5xx). Delete operations must use a combined predicate:
>
> ```go
> retry.WithRetryIf(func(err error) bool {
>     return retry.IsTransientAPIError(err) && !errors.Is(err, api.ErrResourceDoesNotExist)
> })
> ```

## Consequences

### Positive

- Consistent error messages across all resources
- Error chain preserved for debugging (`err.Error()` includes all layers)
- Sentinel errors enable correct Read/Delete behavior without string matching
- Retry policies prevent flaky behavior from transient failures

### Negative

- Existing resources using other formats ("Error", "Failed to", etc.) should be migrated over time
- The `errors.Join` pattern for not-found detection relies on string matching in the HTTP layer for 500 responses

### Common Mistakes

- Using string matching (`strings.Contains(err.Error(), "not found")`) instead of `errors.Is(err, api.ErrResourceDoesNotExist)`.
- Returning errors from Delete when the resource is already gone — check `errors.Is(err, api.ErrResourceDoesNotExist)` first.
- Using `"Error"`, `"Failed to"`, or `"Could not"` prefixes in new code — use `"Unable to [Action] [Resource]"`.
- Wrapping errors with `fmt.Errorf("...: %v", err)` instead of `%w` — breaks error chain inspection.
- Forgetting `resp.State.RemoveResource(ctx)` in Read when the resource no longer exists.
- Using `IsTransientAPIError` alone as a delete retry predicate — it will retry on `ErrResourceDoesNotExist` when it arrives via HTTP 500. Always combine with `!errors.Is(err, api.ErrResourceDoesNotExist)`. See [Retry Policies](#retry-policies).

## References

- [Reference Examples](reference-examples.md) — error handling in CRUD methods
- [ADR-002: API Client Structure](002-api-client-structure.md) — domain client layer errors
- [ADR-004: Schema Design Conventions](004-schema-design-conventions.md) — model conversion patterns
- `proxmox/api/errors.go` — sentinel error definitions
- `proxmox/api/client.go` — HTTP error wrapping
