package spec

// SpecificationResult is the outcome of evaluating a single [Specification].
// A nil failureReasons slice indicates the specification passed.
type SpecificationResult[R comparable] struct {
	name           string
	failureReasons []R
}

// Name returns the name of the specification that produced this result.
func (r SpecificationResult[R]) Name() string {
	return r.name
}

// Passed reports whether the specification passed (no failure reasons).
func (r SpecificationResult[R]) Passed() bool {
	return r.failureReasons == nil
}

// FailureReasons returns a copy of the failure reasons, or nil if the specification passed.
func (r SpecificationResult[R]) FailureReasons() []R {
	if r.failureReasons == nil {
		return nil
	}

	out := make([]R, len(r.failureReasons))
	copy(out, r.failureReasons)

	return out
}

// Pass constructs a passing SpecificationResult for the given name.
func Pass[R comparable](name string) SpecificationResult[R] {
	return SpecificationResult[R]{name: name}
}

// Fail constructs a failing SpecificationResult with one or more reasons.
func Fail[R comparable](name string, reasons ...R) SpecificationResult[R] {
	r := make([]R, len(reasons))
	copy(r, reasons)

	return SpecificationResult[R]{name: name, failureReasons: r}
}
