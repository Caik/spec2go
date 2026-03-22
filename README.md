# 📋 spec2go

A Go implementation of the [Specification Pattern](https://en.wikipedia.org/wiki/Specification_pattern) for composable, reusable business rules.

**Stop scattering validation logic across your codebase.** spec2go lets you define small, testable business rules as *specifications* and combine them into *policies* — making your domain logic explicit, reusable, and easy to reason about.

[![CI](https://github.com/caik/spec2go/actions/workflows/ci.yml/badge.svg)](https://github.com/caik/spec2go/actions/workflows/ci.yml)
[![codecov](https://codecov.io/gh/caik/spec2go/graph/badge.svg)](https://codecov.io/gh/caik/spec2go)
[![Go Reference](https://pkg.go.dev/badge/github.com/caik/spec2go/pkg/spec.svg)](https://pkg.go.dev/github.com/caik/spec2go/pkg/spec)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Go](https://img.shields.io/badge/Go-1.21+-blue.svg)](https://go.dev)

## 📑 Table of Contents

- [Overview](#-overview)
- [Features](#-features)
- [Installation](#-installation)
- [Core Concepts](#-core-concepts)
- [Examples](#-examples)
- [Building](#-building)
- [Contributing](#-contributing)
- [License](#%EF%B8%8F-license)

## 🔍 Overview

spec2go provides a clean, type-safe way to define and evaluate business rules using Go generics. Instead of scattering validation logic throughout your codebase, you define atomic **specifications** that can be composed into **policies**.

```go
// Define failure reasons as typed constants
type Reason string

const (
    Underage        Reason = "UNDERAGE"
    EmailNotVerified Reason = "EMAIL_NOT_VERIFIED"
)

// Define specifications
isAdult := spec.New("IsAdult", func(u User) bool { return u.Age >= 18 }, Underage)
hasVerifiedEmail := spec.New("HasVerifiedEmail", func(u User) bool { return u.EmailVerified }, EmailNotVerified)

// Build a policy
registrationPolicy := spec.NewPolicy[User, Reason]().
    With(isAdult).
    With(hasVerifiedEmail)

// Evaluate
result := registrationPolicy.EvaluateFailFast(user)

if result.AllPassed() {
    // proceed
} else {
    // handle result.FailureReasons()
}
```

## ✨ Features

- 🧩 **Composable** — Build complex rules from simple, reusable specifications
- 🔒 **Type-safe** — Failure reasons use Go generics (`R comparable`), not bare strings
- ⚡ **Two evaluation modes** — `EvaluateFailFast` (stops on first failure) or `EvaluateAll` (collects all failures)
- 🔗 **Logical operators** — `AllOf`, `AnyOf`, `AnyOfAll`, `Not` for combining specifications
- 🐹 **Idiomatic Go** — Generics, interfaces, package-level functions, `fmt.Stringer`
- 🚫 **Zero dependencies** — Pure Go standard library only

## 📦 Installation

```bash
go get github.com/caik/spec2go
```

```go
import "github.com/caik/spec2go/pkg/spec"
```

## 📚 Core Concepts

### 📌 Specification

A single, atomic condition that evaluates a context and returns pass/fail with a reason.

**Simple spec via `New`:**

```go
minimumAge := spec.New("MinimumAge",
    func(a LoanApplication) bool { return a.ApplicantAge >= 18 },
    ReasonApplicantTooYoung,
)
```

**Custom struct spec (for complex logic or multiple failure reasons):**

```go
type DocumentCheck struct {
    spec.NamedSpec[Claim, ClaimReason]
}

func (s DocumentCheck) Evaluate(ctx Claim) spec.SpecificationResult[ClaimReason] {
    var missing []ClaimReason
	
    if !ctx.HasIDDocument { 
		missing = append(missing, MissingID) 
	}
    
	if !ctx.HasProofOfLoss { 
		missing = append(missing, MissingProofOfLoss) 
	}
	
    if len(missing) == 0 {
        return spec.Pass[ClaimReason](s.Name())
    }
	
    return spec.Fail(s.Name(), missing...)
}
```

### 📜 Policy

An ordered collection of specifications evaluated as a unit:

```go
loanPolicy := spec.NewPolicy[LoanApplication, Reason]().
    With(minimumAge).
    With(maximumAge).
    With(creditCheck)

// Fail-fast: stop at first failure (best for performance)
result := loanPolicy.EvaluateFailFast(application)

// Evaluate all: collect every failure (best for showing all errors to the user)
result := loanPolicy.EvaluateAll(application)
```

### 🔗 Composites

Combine specifications with logical operators:

```go
// AND — all must pass (always evaluates all)
fullyVerified := spec.AllOf("FullyVerified", emailVerified, phoneVerified)

// OR — at least one must pass (short-circuits on first pass)
hasPayment := spec.AnyOf("HasPayment", hasCreditCard, hasBankAccount)

// OR — at least one must pass (evaluates all, no short-circuit)
hasPayment := spec.AnyOfAll("HasPayment", hasCreditCard, hasBankAccount)

// NOT — inverts the result
notBlocked := spec.Not("NotBlocked", ReasonBlocked, isBlockedCountry)
```

Composites can be nested arbitrarily and produce human-readable expressions:

```go
fmt.Println(loanPolicy.String())
// (MinimumAge AND (GoodCreditScore OR (SufficientIncome AND IsEmployed)))
```

### 💡 Failure Reasons

Since Go has no enums, define failure reasons as typed constants — any `comparable` type works:

```go
// Typed string constants (recommended — prints meaningfully without extra code)
type LoanReason string

const (
    AgeTooYoung        LoanReason = "AGE_TOO_YOUNG"
    InsufficientIncome LoanReason = "INSUFFICIENT_INCOME"
    PoorCreditScore    LoanReason = "POOR_CREDIT_SCORE"
)

// Typed int with iota (more memory-efficient for large sets)
type LoanReason int

const (
    AgeTooYoung LoanReason = iota
    InsufficientIncome
    PoorCreditScore
)
```

## 🎯 Examples

The `examples/` directory contains complete working examples:

```bash
# Loan eligibility — basic specs + AnyOf/AllOf nesting
go run ./examples/loan/

# E-commerce order validation — custom struct spec + Not
go run ./examples/ecommerce/

# Feature access control — multiple policies + dynamic spec factories
go run ./examples/accesscontrol/
```

## 🛠️ Building

```bash
# Build all packages
go build ./...

# Run all tests
go test ./...

# Run tests with race detector and coverage
go test ./... -race -coverprofile=coverage.out

# View coverage report
go tool cover -html=coverage.out
```

## 🤝 Contributing

Contributions are welcome! Please see [CONTRIBUTING.md](CONTRIBUTING.md) for guidelines.

## ⚖️ License

[![License](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)

Released 2026 by [Carlos Henrique Severino](https://github.com/Caik)
