# 🤝 Contributing to spec2go

Thank you for your interest in contributing to spec2go! 🎉 This document provides guidelines and instructions for contributing.

## 🚀 Getting Started

1. **Fork** the repository on GitHub
2. **Clone** your fork locally:
   ```bash
   git clone https://github.com/YOUR_USERNAME/spec2go.git
   cd spec2go
   ```
3. **Build** the project:
   ```bash
   go build ./...
   ```
4. **Run tests** to confirm everything works:
   ```bash
   go test ./...
   ```

## 💻 Development Workflow

### 🌿 Branching

Create a feature branch from `main`:

```bash
git checkout -b feature/your-feature-name
```

Use descriptive branch names: `feature/add-xor-composite`, `fix/nil-failure-reasons`, `docs/improve-readme`

### ✏️ Making Changes

1. Write your code following the existing style
2. Add tests for new functionality — aim for 100% coverage
3. Ensure all tests pass with the race detector:
   ```bash
   go test ./... -race -count=1
   ```
4. Ensure code is formatted:
   ```bash
   gofmt -s -w .
   ```
5. Ensure vet passes:
   ```bash
   go vet ./...
   ```
6. Ensure `go.mod` is tidy:
   ```bash
   go mod tidy
   ```

### 💬 Commit Messages

We follow [Conventional Commits](https://www.conventionalcommits.org/) for clear, structured commit history. This also drives automatic semantic versioning — please follow the format carefully.

**Format:**

```
<type>[optional scope]: <description>

[optional body]
```

**Types:**

| Type | Description | Version bump |
|---|---|---|
| `feat` | New feature | minor |
| `fix` | Bug fix | patch |
| `perf` | Performance improvement | patch |
| `docs` | Documentation only | none |
| `test` | Adding or updating tests | none |
| `refactor` | Code change that neither fixes nor adds features | patch |
| `ci` | CI/CD changes | none |
| `chore` | Maintenance (dependencies, build) | none |

Append `!` for breaking changes (triggers a major bump):

```
feat!: rename EvaluateFailFast to EvaluateFirst
```

**Examples:**

```
feat: add XorOf composite specification
fix: handle empty failure reasons slice in Fail constructor
docs: add custom struct spec example to README
test: add table-driven boundary tests for AllOf
perf: avoid allocation in Pass constructor
refactor: simplify buildExpression helper
feat!: rename NamedSpec field N to Name
```

### 🔀 Pull Requests

1. Push your branch to your fork
2. Open a Pull Request against `main`
3. Fill out the PR description:
   - What does this PR do?
   - Why is this change needed?
   - How was it tested?
4. Wait for CI to pass and address any feedback

## 🎨 Code Style

- Follow standard Go conventions ([Effective Go](https://go.dev/doc/effective_go), [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments))
- Use `gofmt -s` for formatting (enforced in CI)
- Exported symbols must have godoc comments
- Keep functions small and focused
- Prefer table-driven tests for multiple input variants

## 🧪 Testing

- Write tests for all new functionality
- Place tests in `pkg/spec/` alongside the source files, using `package spec_test`
- Shared test fixtures go in `fixtures_test.go`
- Use descriptive test names:
  ```go
  func TestAllOf_CollectsAllFailureReasonsOnMultipleFailures(t *testing.T) { ... }
  ```
- Check coverage after changes:
  ```bash
  go test ./... -coverprofile=coverage.out && go tool cover -func=coverage.out
  ```

## 🐛 Reporting Issues

When reporting bugs, please include:

1. A clear description of the issue
2. Steps to reproduce
3. Expected vs actual behavior
4. Go version (`go version`)
5. Minimal code example if applicable

## 💡 Feature Requests

Feature requests are welcome! Please:

1. Check existing issues first to avoid duplicates
2. Describe the use case and motivation
3. Provide examples of how it would be used

## ❓ Questions?

Feel free to open an issue for questions or discussions.

---

Thank you for contributing! 🎉
