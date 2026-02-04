# ADR-005: Error Handling Patterns

## Status

Accepted

## Date

2026-02-04

## Context

The provider has a three-layer architecture (HTTP client → domain client → resource) and each layer handles errors differently. Without documented conventions, error message formats diverge across resources — the codebase currently mixes "Unable to", "Error", "Failed to", "Could not", and "Cannot" prefixes. Contributors need a single pattern to follow.

Additionally, retry behavior for transient failures, async task polling, and resource state waiting needs clear guidance.

## Decision

### Error Message Format

User-facing diagnostic summaries (the first argument to `resp.Diagnostics.AddError`) use the format:

```text
"Unable to [Action] [Resource]"
```

Where `[Action]` is a CRUD verb and `[Resource]` identifies the Proxmox object.

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
errors.Join(ErrResourceDoesNotExist, &HTTPError{Code: 404, Message: "..."})
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

| Scenario                                          | Strategy | Attempts      | Backoff                                 |
|---------------------------------------------------|----------|---------------|-----------------------------------------|
| Transient errors (network issues, 5xx responses)  | Retry    | Up to 3       | Linear                                  |
| Async task status checks                          | Retry    | Up to 5       | Exponential, with configurable timeout  |
| Resource state polling (creation/deletion)        | Poll     | Until timeout | Configurable interval and timeout       |

Do not retry on:

- 4xx client errors (except 408 Request Timeout)
- `ErrResourceDoesNotExist` (unless polling for creation)
- Authentication failures (401, 403)

## Consequences

### Positive

- Consistent error messages across all resources
- Error chain preserved for debugging (`err.Error()` includes all layers)
- Sentinel errors enable correct Read/Delete behavior without string matching
- Retry policies prevent flaky behavior from transient failures

### Negative

- Existing resources using other formats ("Error", "Failed to", etc.) should be migrated over time
- The `errors.Join` pattern for not-found detection relies on string matching in the HTTP layer for 500 responses

## References

- [Reference Examples](reference-examples.md) — error handling in CRUD methods
- [ADR-004: Schema Design Conventions](004-schema-design-conventions.md) — model conversion patterns
- `proxmox/api/errors.go` — sentinel error definitions
- `proxmox/api/client.go` — HTTP error wrapping
