package spec_test

import (
	"testing"

	"github.com/caik/spec2go/pkg/spec"
)

// --- NewPolicy / With ---

func TestNewPolicy_StartsEmpty(t *testing.T) {
	p := spec.NewPolicy[testCtx, testReason]()
	result := p.EvaluateAll(testCtx{})
	if !result.AllPassed() {
		t.Error("empty policy should pass")
	}
}

func TestPolicy_WithReturnsSamePolicy(t *testing.T) {
	p := spec.NewPolicy[testCtx, testReason]()
	returned := p.With(alwaysPass)
	if p != returned {
		t.Error("With() should return the same policy pointer")
	}
}

func TestPolicy_FluentChain(t *testing.T) {
	p := spec.NewPolicy[testCtx, testReason]().
		With(isAdult).
		With(hasFunds).
		With(isVerified)

	result := p.EvaluateAll(testCtx{Age: 20, Balance: 200, Verified: true})
	if !result.AllPassed() {
		t.Error("expected all to pass")
	}
}

// --- EvaluateFailFast ---

func TestEvaluateFailFast_PassesWhenAllSpecsPass(t *testing.T) {
	p := spec.NewPolicy[testCtx, testReason]().With(isAdult).With(hasFunds)
	result := p.EvaluateFailFast(testCtx{Age: 20, Balance: 200})
	if !result.AllPassed() {
		t.Error("expected pass")
	}
}

func TestEvaluateFailFast_StopsAtFirstFailure(t *testing.T) {
	evalCount := 0
	second := spec.New("Second", func(c testCtx) bool {
		evalCount++
		return true
	}, blocked)

	p := spec.NewPolicy[testCtx, testReason]().With(alwaysFail).With(second)
	p.EvaluateFailFast(testCtx{})

	if evalCount != 0 {
		t.Errorf("evalCount = %d, want 0 (should stop at first failure)", evalCount)
	}
}

func TestEvaluateFailFast_ResultsContainOnlyEvaluated(t *testing.T) {
	p := spec.NewPolicy[testCtx, testReason]().With(alwaysFail).With(alwaysPass)
	result := p.EvaluateFailFast(testCtx{})

	// Only the first spec was evaluated before stopping
	if len(result.Results()) != 1 {
		t.Errorf("len(Results()) = %d, want 1", len(result.Results()))
	}
}

func TestEvaluateFailFast_ReturnsFailureReasonOfFirstFailedSpec(t *testing.T) {
	p := spec.NewPolicy[testCtx, testReason]().With(isAdult).With(hasFunds)
	result := p.EvaluateFailFast(testCtx{Age: 16, Balance: 50})

	reasons := result.FailureReasons()
	if len(reasons) != 1 || reasons[0] != tooYoung {
		t.Errorf("FailureReasons() = %v, want [%v]", reasons, tooYoung)
	}
}

// --- EvaluateAll ---

func TestEvaluateAll_PassesWhenAllSpecsPass(t *testing.T) {
	p := spec.NewPolicy[testCtx, testReason]().With(isAdult).With(hasFunds)
	result := p.EvaluateAll(testCtx{Age: 20, Balance: 200})
	if !result.AllPassed() {
		t.Error("expected pass")
	}
}

func TestEvaluateAll_EvaluatesAllSpecs(t *testing.T) {
	evalCount := 0
	s1 := spec.New("S1", func(c testCtx) bool { evalCount++; return false }, tooYoung)
	s2 := spec.New("S2", func(c testCtx) bool { evalCount++; return false }, insufficientFunds)

	p := spec.NewPolicy[testCtx, testReason]().With(s1).With(s2)
	p.EvaluateAll(testCtx{})

	if evalCount != 2 {
		t.Errorf("evalCount = %d, want 2", evalCount)
	}
}

func TestEvaluateAll_CollectsAllFailureReasons(t *testing.T) {
	p := spec.NewPolicy[testCtx, testReason]().With(isAdult).With(hasFunds).With(isVerified)
	result := p.EvaluateAll(testCtx{Age: 16, Balance: 50, Verified: false})

	if result.AllPassed() {
		t.Error("expected fail")
	}
	if len(result.FailureReasons()) != 3 {
		t.Errorf("len(FailureReasons()) = %d, want 3", len(result.FailureReasons()))
	}
}

func TestEvaluateAll_ResultsContainAll(t *testing.T) {
	p := spec.NewPolicy[testCtx, testReason]().With(alwaysFail).With(alwaysPass)
	result := p.EvaluateAll(testCtx{})

	if len(result.Results()) != 2 {
		t.Errorf("len(Results()) = %d, want 2", len(result.Results()))
	}
}

// --- String / Expression ---

func TestPolicy_StringEmpty(t *testing.T) {
	p := spec.NewPolicy[testCtx, testReason]()
	if p.String() != "()" {
		t.Errorf("String() = %q, want %q", p.String(), "()")
	}
}

func TestPolicy_StringSingleSpec(t *testing.T) {
	p := spec.NewPolicy[testCtx, testReason]().With(isAdult)
	if p.String() != "(IsAdult)" {
		t.Errorf("String() = %q, want %q", p.String(), "(IsAdult)")
	}
}

func TestPolicy_StringMultipleSpecs(t *testing.T) {
	p := spec.NewPolicy[testCtx, testReason]().With(isAdult).With(hasFunds)
	if p.String() != "(IsAdult AND HasFunds)" {
		t.Errorf("String() = %q, want %q", p.String(), "(IsAdult AND HasFunds)")
	}
}

func TestPolicy_StringWithComposite(t *testing.T) {
	financiallyQualified, err := spec.AnyOf("FinanciallyQualified", hasFunds, isVerified)
	if err != nil {
		t.Fatal(err)
	}
	p := spec.NewPolicy[testCtx, testReason]().With(isAdult).With(financiallyQualified)
	want := "(IsAdult AND (HasFunds OR IsVerified))"
	if p.String() != want {
		t.Errorf("String() = %q, want %q", p.String(), want)
	}
}
