# ADR-006: Testing Requirements

## Status

Accepted

## Date

2026-02-04 (retroactive documentation)

## Context

Contributions need clear testing expectations. Without documented requirements, some PRs arrive with no tests, while others over-invest in mocking infrastructure that does not reflect real API behavior. The provider relies heavily on acceptance tests against a live Proxmox environment to verify correctness.

## Decision

### Test Coverage Requirements

| Test Type              | Requirement                                      | Purpose                                             |
| ---------------------- | ------------------------------------------------ | --------------------------------------------------- |
| Acceptance tests       | **Required** for all new resources and bug fixes | Verify real API behavior end-to-end                 |
| Unit tests             | **Recommended** for complex logic                | Test parsing, validation, model conversion          |
| Example configurations | **Optional**                                     | Legacy tests and user-facing examples in `example/` |

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
| ------------------------------------------- | ----------------------------------------------- |
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

> [!NOTE]
> Meeting the requirements above (create + update + import + field removal) corresponds to a D6 score of 2/3 in the [compliance scoring rubric](reference-examples.md#d6-testing-quality-weight-3). For Grade A compliance (D6=3), also add table-driven structure with named scenarios, validation tests, and `test.ResourceAttributes` / `test.NoResourceAttributesSet` bulk assertions.

### Functional Coverage Requirement

Acceptance tests must cover **all major use cases** for the resource, not just the happy path. PRs that only test one basic scenario will be rejected during review. For example, a backup job resource that only tests `all = true` but never tests targeting specific VMs by ID is incomplete — the VM targeting is a core use case.

When planning tests, identify the distinct operational modes of the resource and ensure each has at least one test scenario:

- Different input combinations (e.g., `all` vs `vmid` vs `pool` for backup targets)
- List/set attributes with multiple elements (not just empty or single-element)
- Compound string attributes that round-trip through the API (e.g., `prune_backups`)
- Nested object attributes (e.g., `fleecing`, `performance`)
- Import with non-trivial state (e.g., list attributes that must survive import round-trip)

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
                Config:   `resource "proxmox_my_resource" "test" { ... }`,
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

Tests passing does not guarantee correct API calls. After tests pass, verify API calls with mitmproxy using the "debug API" workflow. See [DEBUGGING.md](../../.dev/DEBUGGING.md).

**Proving the absence of an API call needs a positive control.** When a change's effect is that a call is _no longer made_ (a skip flag, a removed redundant request), "grepped the flow log, found 0 matches" is weak evidence — 0 also results from a wrong endpoint string or a test that never reached the code path. Make it conclusive: confirm the exact endpoint string from the client code first, assert that the _other_ expected calls are present in the same capture (proving the capture worked), and run a control scenario where the call _should_ fire, confirming a non-zero count. The 0-vs-N contrast across two near-identical runs is the proof; a bare 0 is not.

**Write-only and secret values need out-of-band assertions.** Write-only attribute values are never in state, and PVE redacts many secrets from GET responses — state checks and direct API reads both come back empty whether or not the secret arrived. Many such secrets are observable as root-only files under `/etc/pve/priv/` on the PVE host (e.g. `metricserver/<id>.pw`); assert via root SSH:

```go
out := te.ExecuteNodeCommands([]string{
    fmt.Sprintf("cat /etc/pve/priv/metricserver/%s.pw 2>/dev/null || echo MISSING", name),
})
exists := !strings.Contains(out, "MISSING")
```

Assert the file exists after create-with-secret and is gone after the secret's removal — this catches delete-path gaps that state-only checks cannot. Before reaching for mitmproxy, probe `ls -R /etc/pve/priv/` on the test host for a per-resource credentials file.

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

### Common Mistakes

- Forgetting `//go:build acceptance || all` build tag — test won't run in CI.
- Using `resource.Test` instead of `resource.ParallelTest` — slows test suite.
- Using `t.Run` without `t.Parallel()` at the top-level test function.
- Including issue numbers in test names or resource names.
- Assuming passing tests means correct API behavior — always verify with mitmproxy.

## References

- [ADR-003: Resource File Organization](003-resource-file-organization.md) — file placement
- [ADR-005: Error Handling](005-error-handling.md) — error patterns tested by acceptance tests
- [Reference Examples](reference-examples.md) — acceptance test walkthrough
- [Terraform Plugin Testing](https://developer.hashicorp.com/terraform/plugin/testing)
