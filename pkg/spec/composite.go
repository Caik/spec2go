package spec

import (
	"fmt"
	"strings"
)

// AnyOf returns a [Specification] that passes if any of the given specifications pass (OR logic).
// Evaluation short-circuits on the first passing specification.
// Panics if no specifications are provided.
func AnyOf[T any, R comparable](name string, specs ...Specification[T, R]) Specification[T, R] {
	if len(specs) == 0 {
		panic("AnyOf requires at least one specification")
	}
	expr := buildExpression(specs, " OR ")
	return Func[T, R]{
		name:       name,
		expression: expr,
		fn: func(ctx T) SpecificationResult[R] {
			var reasons []R
			for _, s := range specs {
				res := s.Evaluate(ctx)
				if res.Passed() {
					return Pass[R](name)
				}
				reasons = append(reasons, res.FailureReasons()...)
			}
			return Fail(name, reasons...)
		},
	}
}

// AnyOfAll returns a [Specification] that passes if any of the given specifications pass (OR logic).
// Unlike [AnyOf], all specifications are always evaluated — no short-circuiting.
// Panics if no specifications are provided.
func AnyOfAll[T any, R comparable](name string, specs ...Specification[T, R]) Specification[T, R] {
	if len(specs) == 0 {
		panic("AnyOfAll requires at least one specification")
	}
	expr := buildExpression(specs, " OR ")
	return Func[T, R]{
		name:       name,
		expression: expr,
		fn: func(ctx T) SpecificationResult[R] {
			anyPassed := false
			var reasons []R
			for _, s := range specs {
				res := s.Evaluate(ctx)
				if res.Passed() {
					anyPassed = true
				} else {
					reasons = append(reasons, res.FailureReasons()...)
				}
			}
			if anyPassed {
				return Pass[R](name)
			}
			return Fail(name, reasons...)
		},
	}
}

// AllOf returns a [Specification] that passes only if all of the given specifications pass (AND logic).
// All specifications are always evaluated regardless of intermediate results.
// Panics if no specifications are provided.
func AllOf[T any, R comparable](name string, specs ...Specification[T, R]) Specification[T, R] {
	if len(specs) == 0 {
		panic("AllOf requires at least one specification")
	}
	expr := buildExpression(specs, " AND ")
	return Func[T, R]{
		name:       name,
		expression: expr,
		fn: func(ctx T) SpecificationResult[R] {
			allPassed := true
			var reasons []R
			for _, s := range specs {
				res := s.Evaluate(ctx)
				if !res.Passed() {
					allPassed = false
					reasons = append(reasons, res.FailureReasons()...)
				}
			}
			if allPassed {
				return Pass[R](name)
			}
			return Fail(name, reasons...)
		},
	}
}

// Not returns a [Specification] that passes when the given specification fails, and fails otherwise.
// The provided failureReason is used when the inner specification passes (i.e., the Not fails).
func Not[T any, R comparable](name string, failureReason R, s Specification[T, R]) Specification[T, R] {
	expr := fmt.Sprintf("NOT %s", s.Expression())
	return Func[T, R]{
		name:       name,
		expression: expr,
		fn: func(ctx T) SpecificationResult[R] {
			res := s.Evaluate(ctx)
			if !res.Passed() {
				return Pass[R](name)
			}
			return Fail(name, failureReason)
		},
	}
}

func buildExpression[T any, R comparable](specs []Specification[T, R], sep string) string {
	if len(specs) == 1 {
		return specs[0].Expression()
	}
	parts := make([]string, len(specs))
	for i, s := range specs {
		parts[i] = s.Expression()
	}
	return fmt.Sprintf("(%s)", strings.Join(parts, sep))
}
