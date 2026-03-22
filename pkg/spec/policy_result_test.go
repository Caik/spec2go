package spec_test

import (
	"testing"

	"github.com/caik/spec2go/pkg/spec"
)

func TestPolicyResult_AllPassedTrue(t *testing.T) {
	p := spec.NewPolicy[testCtx, testReason]().With(alwaysPass).With(isAdult)
	result := p.EvaluateAll(testCtx{Age: 20})
	if !result.AllPassed() {
		t.Error("AllPassed() = false, want true")
	}
}

func TestPolicyResult_AllPassedFalse(t *testing.T) {
	p := spec.NewPolicy[testCtx, testReason]().With(alwaysPass).With(alwaysFail)
	result := p.EvaluateAll(testCtx{})
	if result.AllPassed() {
		t.Error("AllPassed() = true, want false")
	}
}

func TestPolicyResult_FailedResults(t *testing.T) {
	p := spec.NewPolicy[testCtx, testReason]().With(isAdult).With(hasFunds).With(isVerified)
	result := p.EvaluateAll(testCtx{Age: 16, Balance: 50, Verified: false})

	failed := result.FailedResults()
	if len(failed) != 3 {
		t.Errorf("len(FailedResults()) = %d, want 3", len(failed))
	}
}

func TestPolicyResult_FailedResults_EmptyWhenAllPass(t *testing.T) {
	p := spec.NewPolicy[testCtx, testReason]().With(alwaysPass)
	result := p.EvaluateAll(testCtx{})

	if len(result.FailedResults()) != 0 {
		t.Errorf("FailedResults() should be empty when all pass")
	}
}

func TestPolicyResult_FailureReasons(t *testing.T) {
	p := spec.NewPolicy[testCtx, testReason]().With(isAdult).With(hasFunds)
	result := p.EvaluateAll(testCtx{Age: 16, Balance: 50})

	reasons := result.FailureReasons()
	if len(reasons) != 2 {
		t.Errorf("len(FailureReasons()) = %d, want 2", len(reasons))
	}
}

func TestPolicyResult_FailureReasons_EmptyWhenAllPass(t *testing.T) {
	p := spec.NewPolicy[testCtx, testReason]().With(alwaysPass)
	result := p.EvaluateAll(testCtx{})

	if len(result.FailureReasons()) != 0 {
		t.Errorf("FailureReasons() should be empty when all pass")
	}
}

func TestPolicyResult_Results_ContainsAll(t *testing.T) {
	p := spec.NewPolicy[testCtx, testReason]().With(isAdult).With(hasFunds)
	result := p.EvaluateAll(testCtx{Age: 20, Balance: 200})

	if len(result.Results()) != 2 {
		t.Errorf("len(Results()) = %d, want 2", len(result.Results()))
	}
}
