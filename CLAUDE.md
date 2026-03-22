# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Commands

```bash
# Run all tests
go test ./...

# Run tests with race detection, coverage, and shuffle (matches CI)
go test ./... -v -race -count=1 -shuffle=on -coverprofile=coverage.txt

# Run a single test
go test ./pkg/spec -run TestName

# Check coverage (CI enforces ≥95%)
go tool cover -func=coverage.txt

# Format code 
gofmt -w -s .

# Vet
go vet ./...

# Build
go build ./...
```

## Architecture

`spec2go` is a zero-dependency Go library implementing the [Specification Pattern](https://en.wikipedia.org/wiki/Specification_pattern). All library code lives in `pkg/spec/`.

### Core Interface

```go
type Specification[T any, R comparable] interface {
    Evaluate(ctx T) SpecificationResult[R]
    Name() string
    Expression() string
}
```

- `T` = context type being evaluated (e.g., `User`, `Order`)
- `R` = comparable failure reason type (e.g., a typed `string` const)

### Key Files

| File | Role |
|------|------|
| `spec.go` | `Specification` interface, `Func` struct, `New()` constructor, `NamedSpec` embed helper |
| `result.go` | `SpecificationResult[R]` with `Pass()`/`Fail()` constructors |
| `policy.go` | `Policy` — groups specs, chains via `With()`, evaluates via `EvaluateFailFast()` or `EvaluateAll()` |
| `policy_result.go` | `PolicyResult[R]` — aggregated outcome with `AllPassed`, `FailedResults`, `FailureReasons` |
| `composite.go` | Logical operators: `AllOf()`, `AnyOf()`, `AnyOfAll()`, `Not()` — all return `error` if given no specs |

### Two Ways to Create a Spec

1. **Simple** — single predicate + one failure reason:
   ```go
   spec.New("IsAdult", func(u User) bool { return u.Age >= 18 }, Underage)
   ```

2. **Custom struct** — embed `NamedSpec` for complex multi-reason logic:
   ```go
   type DocumentCheck struct { spec.NamedSpec[Claim, ClaimReason] }
   func (s DocumentCheck) Evaluate(ctx Claim) spec.SpecificationResult[ClaimReason] { ... }
   ```

### Two Evaluation Modes

| Method | Behavior |
|--------|----------|
| `EvaluateFailFast()` | Stops at first failure |
| `EvaluateAll()` | Evaluates all specs, collects all failures |

### Composite Operators

- `AllOf()` — AND, always evaluates all
- `AnyOf()` — OR, short-circuits on first pass
- `AnyOfAll()` — OR, no short-circuit
- `Not()` — inverts a spec (requires explicit failure reason)

### Expression Generation

Each spec produces a human-readable expression for debugging:
```
(MinimumAge AND (GoodCreditScore OR (SufficientIncome AND IsEmployed)))
```

### Testing Conventions

- Tests are co-located with source in `pkg/spec/`
- `fixtures_test.go` defines shared `testCtx`, `testReason`, and pre-built specs for reuse
- CI enforces **95% coverage** — keep coverage high when adding features

### Examples

Three runnable examples live in `examples/` (no tests, just demonstrations):
- `loan/` — nested composite specs
- `ecommerce/` — custom struct-based specs
- `accesscontrol/` — dynamic spec factories

### Release Process

Releases are automated via `.github/workflows/release.yml`. Version bumps are detected from commit message prefixes (`feat:` → minor, `fix:` → patch, `BREAKING CHANGE` → major). Changes only to tests, docs, examples, or `.github/` do not trigger a release.
