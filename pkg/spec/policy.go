package spec

import (
	"fmt"
	"strings"
)

// Policy is a named collection of [Specification]s evaluated as a unit.
// Use [NewPolicy] to create a policy, then chain [Policy.With] calls to add specifications.
type Policy[T any, R comparable] struct {
	specs []Specification[T, R]
}

// NewPolicy creates an empty Policy.
func NewPolicy[T any, R comparable]() *Policy[T, R] {
	return &Policy[T, R]{}
}

// With adds a Specification to the policy and returns the same policy for method chaining.
func (p *Policy[T, R]) With(s Specification[T, R]) *Policy[T, R] {
	p.specs = append(p.specs, s)
	return p
}

// EvaluateFailFast evaluates specifications in order and stops at the first failure.
// Use this when you only need to know whether a policy passes and want minimal evaluation.
func (p *Policy[T, R]) EvaluateFailFast(ctx T) PolicyResult[R] {
	var results []SpecificationResult[R]
	for _, s := range p.specs {
		res := s.Evaluate(ctx)
		results = append(results, res)
		if !res.Passed() {
			return newPolicyResult(false, results)
		}
	}
	return newPolicyResult(true, results)
}

// EvaluateAll evaluates every specification and collects all failures.
// Use this when you need a complete picture of all failing rules.
func (p *Policy[T, R]) EvaluateAll(ctx T) PolicyResult[R] {
	allPassed := true
	results := make([]SpecificationResult[R], 0, len(p.specs))
	for _, s := range p.specs {
		res := s.Evaluate(ctx)
		results = append(results, res)
		if !res.Passed() {
			allPassed = false
		}
	}
	return newPolicyResult(allPassed, results)
}

// String returns a SQL-like expression representing the policy's specification structure.
// Implements [fmt.Stringer].
func (p *Policy[T, R]) String() string {
	if len(p.specs) == 0 {
		return "()"
	}
	parts := make([]string, len(p.specs))
	for i, s := range p.specs {
		parts[i] = s.Expression()
	}
	return fmt.Sprintf("(%s)", strings.Join(parts, " AND "))
}
