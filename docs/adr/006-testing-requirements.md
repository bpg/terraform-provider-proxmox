# ADR-006: Testing Requirements

## Status

Accepted

## Date

2026-02-04

## Context

Contributions need clear testing expectations. Without documented requirements, some PRs arrive with no tests, while others over-invest in mocking infrastructure that does not reflect real API behavior. The provider relies heavily on acceptance tests against a live Proxmox environment to verify correctness.

## Decision

### Test Coverage Requirements

| Test Type              | Requirement                                      | Purpose                                    |
|------------------------|--------------------------------------------------|--------------------------------------------|
| Acceptance tests       | **Required** for all new resources and bug fixes | Verify real API behavior end-to-end        |
| Unit tests             | **Recommended** for complex logic                | Test parsing, validation, model conversion |
| Example configurations | **Optional**                                     | User-facing examples in `example/`         |

### Acceptance Test Structure

Acceptance tests use the Terraform testing framework and run against a live Proxmox instance.

#### File Placement

Test files are colocated with their resource files (see [ADR-003](003-resource-file-organization.md)):

```text
fwprovider/cluster/sdn/vnet/
├── resource.go
├── model.go
└── resource_test.go      # Acceptance tests here
```

#### Build Tags

All acceptance test files must include the build tag:

```go
//go:build acceptance || all
```

#### Test Environment Setup

Use `test.InitEnvironment(t)` from `fwprovider/test/` to get provider factories and configuration helpers:

```go
func TestAccResourceSDNVNet(t *testing.T) {
    t.Parallel()

    te := test.InitEnvironment(t)
    // te.AccProviders — pre-configured provider factories
    // te.RenderConfig() — template rendering with {{.NodeName}}
}
```

#### Table-Driven Tests

Structure tests as table-driven with named scenarios:

```go
tests := []struct {
    name  string
    steps []resource.TestStep
}{
    {"create and update vnet", []resource.TestStep{
        {
            Config: te.RenderConfig(`...`),
            Check:  resource.ComposeTestCheckFunc(...),
        },
        {
            Config:            te.RenderConfig(`...updated...`),
            ImportState:       true,
            ImportStateVerify: true,
        },
    }},
}
```

#### Running Tests

Use `resource.ParallelTest` (not `resource.Test`) for acceptance tests that can run concurrently:

```go
for _, tt := range tests {
    t.Run(tt.name, func(t *testing.T) {
        resource.ParallelTest(t, resource.TestCase{
            ProtoV6ProviderFactories: te.AccProviders,
            Steps:                    tt.steps,
        })
    })
}
```

### Test Helpers

The `fwprovider/test/` package provides shared utilities:

| Helper                                      | Purpose                                         |
|---------------------------------------------|-------------------------------------------------|
| `test.InitEnvironment(t)`                   | Provider factories, config rendering, node name |
| `test.ResourceAttributes(addr, map)`        | Bulk attribute assertions                       |
| `test.NoResourceAttributesSet(addr, attrs)` | Assert attributes are null/unset                |
| `te.RenderConfig(template, vars...)`        | Render HCL with `{{.NodeName}}` substitution    |

### Test Scenarios to Cover

For each new resource, tests should cover:

1. **Create** — verify all attributes are set correctly after creation
2. **Update** — change optional fields and verify state reflects changes
3. **Import** — verify `ImportState` + `ImportStateVerify` round-trips correctly
4. **Delete** — covered implicitly by test framework cleanup

Additional scenarios when applicable:

5. **Validation errors** — use `resource.UnitTest` with `PlanOnly: true` and `ExpectError`
6. **Field removal** — verify optional fields can be unset (tests the `CheckDelete` path)

### Validation Tests

Validation logic can be tested without a live Proxmox instance using `resource.UnitTest`:

```go
func TestAccMyResource_Validators(t *testing.T) {
    t.Parallel()
    te := test.InitEnvironment(t)

    resource.UnitTest(t, resource.TestCase{
        ProtoV6ProviderFactories: te.AccProviders,
        Steps: []resource.TestStep{
            {
                PlanOnly: true,
                Config:   `resource "proxmox_virtual_environment_my_resource" "test" { ... }`,
                ExpectError: regexp.MustCompile(`Invalid Attribute Combination`),
            },
        },
    })
}
```

### Test Naming

- Function names: `TestAcc{Resource}{Scenario}` (e.g., `TestAccResourceSDNVNet`)
- Resource names in HCL: descriptive, no issue numbers (e.g., `test-vnet`, `acc_influxdb_server`)
- VM names: descriptive, no issue numbers

### API Verification

Tests passing does not guarantee correct API calls. After tests pass, verify API calls with mitmproxy using the `/debug-api` workflow. See [DEBUGGING.md](../../.dev/DEBUGGING.md).

## Consequences

### Positive

- Every new resource ships with verified behavior
- Regressions are caught by CI
- Import support is verified, not assumed
- Validation tests run without infrastructure

### Negative

- Acceptance tests require a Proxmox environment
- Tests add to PR review scope
- Parallel test execution requires careful resource naming to avoid collisions

## References

- [ADR-003: Resource File Organization](003-resource-file-organization.md)
- [Reference Examples](reference-examples.md) — acceptance test walkthrough
- [Terraform Plugin Testing](https://developer.hashicorp.com/terraform/plugin/testing)
