# ADR-002: Domain-Specific API Client Structure

## Status

Accepted

## Date

2026-02-01 (retroactive documentation)

## Context

The Proxmox VE API is organized hierarchically:

```text
/api2/json/
├── cluster/
│   ├── firewall/
│   ├── ha/
│   ├── acme/
│   ├── sdn/
│   │   ├── zones/
│   │   ├── vnets/
│   │   └── subnets/
│   └── ...
├── nodes/{node}/
│   ├── qemu/{vmid}/
│   ├── lxc/{vmid}/
│   ├── storage/
│   └── ...
├── access/
│   ├── users/
│   ├── groups/
│   └── ...
└── pools/
```

We need a consistent pattern for:

1. Making HTTP requests to these endpoints
2. Organizing client code by domain
3. Providing type-safe access to API operations

## Decision

Use a **layered domain client architecture** with the following structure:

### Layer 1: Base API Client (`proxmox/api/client.go`)

Handles all HTTP communication:

- Authentication (ticket, API token, username/password)
- Request/response serialization
- Retry logic
- Error handling
- TLS configuration

### Layer 2: Top-Level Client (`proxmox/client.go`)

Wraps API client and provides factory methods for domain clients:

```go
type Client struct {
    api    api.Client
    ssh    ssh.Client
    // ...
}

func (c *Client) Cluster() *cluster.Client { ... }
func (c *Client) Node(name string) *nodes.Client { ... }
func (c *Client) Access() *access.Client { ... }
func (c *Client) Pools() *pools.Client { ... }
```

### Layer 3: Domain Clients (`proxmox/{domain}/client.go`)

Each domain client:

1. Embeds `api.Client`
2. Implements `ExpandPath()` for URL construction
3. Provides factory methods for sub-domain clients
4. Contains CRUD methods for domain-specific operations

**Pattern:**

```go
// proxmox/cluster/client.go
package cluster

import "github.com/bpg/terraform-provider-proxmox/proxmox/api"

type Client struct {
    api.Client
}

// ExpandPath prepends the cluster base path
func (c *Client) ExpandPath(path string) string {
    return fmt.Sprintf("cluster/%s", path)
}

// SDNZones returns a client for SDN zone operations
func (c *Client) SDNZones() *zones.Client {
    return &zones.Client{Client: c.Client}
}
```

### Layer 4: Sub-Domain Clients

For deeply nested APIs (e.g., SDN zones → vnets → subnets):

```go
// proxmox/cluster/sdn/zones/client.go
package zones

type Client struct {
    api.Client
}

func (c *Client) ExpandPath(path string) string {
    return fmt.Sprintf("cluster/sdn/zones/%s", path)
}

// VNets returns a client for VNet operations within this zone
func (c *Client) VNets(zoneID string) *vnets.Client {
    return &vnets.Client{
        Client: c.Client,
        ZoneID: zoneID,
    }
}
```

## Consequences

### Positive

- Clear organization matching Proxmox API structure
- Type-safe client access
- Easy to add new API endpoints
- Shared authentication/connection handling
- IDE-friendly discovery via factory methods

### Negative

- Some boilerplate for new domain clients
- Deep nesting can be verbose: `client.Cluster().SDNZones().VNets(id).Subnets()`

### Implementation Guidelines

1. **New API endpoint?** Add methods to existing domain client if it fits
2. **New API domain?** Create new package in `proxmox/`
3. **Factory methods** should return pointer to new client
4. **Path expansion** handles URL construction, callers pass relative paths
5. **Response types** defined in same package as client

## Example: Adding a New API Endpoint

To add support for a new Proxmox API at `/api2/json/cluster/foo/bar`:

1. Create `proxmox/cluster/foo/client.go`:

```go
package foo

import "github.com/bpg/terraform-provider-proxmox/proxmox/api"

type Client struct {
    api.Client
}

func (c *Client) ExpandPath(path string) string {
    return fmt.Sprintf("cluster/foo/%s", path)
}

type Bar struct {
    ID   string `json:"id"`
    Name string `json:"name"`
}

func (c *Client) GetBar(ctx context.Context, id string) (*Bar, error) {
    var result Bar
    err := c.DoRequest(ctx, http.MethodGet, c.ExpandPath(id), nil, &result)
    if err != nil {
        return nil, fmt.Errorf("error getting bar with id %s: %w", id, err)
    }
    return &result, nil
}
```

2. Add factory method to `proxmox/cluster/client.go`:

```go
func (c *Client) Foo() *foo.Client {
    return &foo.Client{Client: c.Client}
}
```

## References

- `proxmox/api/client.go` — Base HTTP client implementation
- `proxmox/cluster/client.go` — Example domain client
- `proxmox/nodes/client.go` — Example with node parameter
- [ADR-003: Resource File Organization](003-resource-file-organization.md) — corresponding `fwprovider/` structure
- [ADR-005: Error Handling](005-error-handling.md) — error patterns in domain clients
