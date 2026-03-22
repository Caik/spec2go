package spec

// PolicyResult is the aggregated outcome of evaluating a [Policy].
type PolicyResult[R comparable] struct {
	allPassed bool
	results   []SpecificationResult[R]
}

// newPolicyResult constructs a PolicyResult with a defensive copy of results.
func newPolicyResult[R comparable](allPassed bool, results []SpecificationResult[R]) PolicyResult[R] {
	r := make([]SpecificationResult[R], len(results))
	copy(r, results)

	return PolicyResult[R]{allPassed: allPassed, results: r}
}

// AllPassed reports whether every specification in the policy passed.
func (pr PolicyResult[R]) AllPassed() bool {
	return pr.allPassed
}

// Results returns all specification results (both passed and failed).
func (pr PolicyResult[R]) Results() []SpecificationResult[R] {
	return pr.results
}

// FailedResults returns only the results where the specification failed.
func (pr PolicyResult[R]) FailedResults() []SpecificationResult[R] {
	var failed []SpecificationResult[R]

	for _, r := range pr.results {
		if !r.Passed() {
			failed = append(failed, r)
		}
	}

	return failed
}

// FailureReasons returns a flattened list of all failure reasons across all failed specifications.
func (pr PolicyResult[R]) FailureReasons() []R {
	var reasons []R

	for _, r := range pr.results {
		if !r.Passed() {
			reasons = append(reasons, r.FailureReasons()...)
		}
	}

	return reasons
}
