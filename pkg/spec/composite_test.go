package spec_test

import (
	"testing"

	"github.com/caik/spec2go/pkg/spec"
)

// --- AnyOf ---

func TestAnyOf_PassesWhenFirstSpecPasses(t *testing.T) {
	s := spec.AnyOf("combo", isAdult, hasFunds)
	result := s.Evaluate(testCtx{Age: 20, Balance: 0})
	if !result.Passed() {
		t.Error("expected pass when first spec passes")
	}
}

func TestAnyOf_PassesWhenSecondSpecPasses(t *testing.T) {
	s := spec.AnyOf("combo", isAdult, hasFunds)
	result := s.Evaluate(testCtx{Age: 16, Balance: 200})
	if !result.Passed() {
		t.Error("expected pass when second spec passes")
	}
}

func TestAnyOf_FailsWhenAllSpecsFail(t *testing.T) {
	s := spec.AnyOf("combo", isAdult, hasFunds)
	result := s.Evaluate(testCtx{Age: 16, Balance: 50})
	if result.Passed() {
		t.Error("expected fail when all specs fail")
	}
}

func TestAnyOf_CollectsAllReasonsOnFailure(t *testing.T) {
	s := spec.AnyOf("combo", isAdult, hasFunds)
	result := s.Evaluate(testCtx{Age: 16, Balance: 50})
	if len(result.FailureReasons()) != 2 {
		t.Errorf("FailureReasons() len = %d, want 2", len(result.FailureReasons()))
	}
}

func TestAnyOf_ShortCircuitsOnFirstPass(t *testing.T) {
	evalCount := 0
	counter := spec.New("Counter", func(c testCtx) bool {
		evalCount++
		return true
	}, blocked)
	never := spec.New("Never", func(c testCtx) bool {
		evalCount++
		return true
	}, blocked)

	s := spec.AnyOf("combo", counter, never)
	s.Evaluate(testCtx{})

	if evalCount != 1 {
		t.Errorf("evalCount = %d, want 1 (short-circuit)", evalCount)
	}
}

func TestAnyOf_Expression(t *testing.T) {
	s := spec.AnyOf("combo", isAdult, hasFunds)
	want := "(IsAdult OR HasFunds)"
	if s.Expression() != want {
		t.Errorf("Expression() = %q, want %q", s.Expression(), want)
	}
}

func TestAnyOf_SingleSpec(t *testing.T) {
	s := spec.AnyOf("single", isAdult)
	if s.Expression() != "IsAdult" {
		t.Errorf("Expression() = %q, want %q", s.Expression(), "IsAdult")
	}
}

func TestAnyOf_PanicsOnEmpty(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("expected panic for empty specs")
		}
	}()
	spec.AnyOf[testCtx, testReason]("empty")
}

// --- AnyOfAll ---

func TestAnyOfAll_EvaluatesAllSpecs(t *testing.T) {
	evalCount := 0
	s1 := spec.New("S1", func(c testCtx) bool { evalCount++; return true }, blocked)
	s2 := spec.New("S2", func(c testCtx) bool { evalCount++; return false }, blocked)

	s := spec.AnyOfAll("combo", s1, s2)
	result := s.Evaluate(testCtx{})

	if evalCount != 2 {
		t.Errorf("evalCount = %d, want 2 (no short-circuit)", evalCount)
	}
	if !result.Passed() {
		t.Error("expected pass because s1 passed")
	}
}

func TestAnyOfAll_FailsWhenAllFail(t *testing.T) {
	s := spec.AnyOfAll("combo", alwaysFail, alwaysFail)
	if s.Evaluate(testCtx{}).Passed() {
		t.Error("expected fail when all specs fail")
	}
}

func TestAnyOfAll_PanicsOnEmpty(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("expected panic for empty specs")
		}
	}()
	spec.AnyOfAll[testCtx, testReason]("empty")
}

func TestAnyOfAll_CollectsAllReasonsOnFailure(t *testing.T) {
	s := spec.AnyOfAll("combo", isAdult, hasFunds)
	result := s.Evaluate(testCtx{Age: 16, Balance: 50})
	if result.Passed() {
		t.Error("expected fail")
	}
	if len(result.FailureReasons()) != 2 {
		t.Errorf("FailureReasons() len = %d, want 2", len(result.FailureReasons()))
	}
}

// --- AllOf ---

func TestAllOf_PassesWhenAllSpecsPass(t *testing.T) {
	s := spec.AllOf("combo", isAdult, hasFunds)
	result := s.Evaluate(testCtx{Age: 20, Balance: 200})
	if !result.Passed() {
		t.Error("expected pass when all specs pass")
	}
}

func TestAllOf_FailsWhenFirstFails(t *testing.T) {
	s := spec.AllOf("combo", isAdult, hasFunds)
	result := s.Evaluate(testCtx{Age: 16, Balance: 200})
	if result.Passed() {
		t.Error("expected fail when first spec fails")
	}
}

func TestAllOf_CollectsAllFailureReasons(t *testing.T) {
	s := spec.AllOf("combo", isAdult, hasFunds)
	result := s.Evaluate(testCtx{Age: 16, Balance: 50})
	if len(result.FailureReasons()) != 2 {
		t.Errorf("FailureReasons() len = %d, want 2", len(result.FailureReasons()))
	}
}

func TestAllOf_EvaluatesAllSpecsEvenOnFailure(t *testing.T) {
	evalCount := 0
	s1 := spec.New("S1", func(c testCtx) bool { evalCount++; return false }, tooYoung)
	s2 := spec.New("S2", func(c testCtx) bool { evalCount++; return false }, insufficientFunds)

	spec.AllOf("combo", s1, s2).Evaluate(testCtx{})

	if evalCount != 2 {
		t.Errorf("evalCount = %d, want 2", evalCount)
	}
}

func TestAllOf_Expression(t *testing.T) {
	s := spec.AllOf("combo", isAdult, hasFunds)
	want := "(IsAdult AND HasFunds)"
	if s.Expression() != want {
		t.Errorf("Expression() = %q, want %q", s.Expression(), want)
	}
}

func TestAllOf_PanicsOnEmpty(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("expected panic for empty specs")
		}
	}()
	spec.AllOf[testCtx, testReason]("empty")
}

// --- Not ---

func TestNot_PassesWhenInnerFails(t *testing.T) {
	s := spec.Not("IsMinor", tooOld, isAdult)
	result := s.Evaluate(testCtx{Age: 16}) // isAdult fails → Not passes
	if !result.Passed() {
		t.Error("expected pass when inner spec fails")
	}
}

func TestNot_FailsWhenInnerPasses(t *testing.T) {
	s := spec.Not("IsMinor", tooOld, isAdult)
	result := s.Evaluate(testCtx{Age: 20}) // isAdult passes → Not fails
	if result.Passed() {
		t.Error("expected fail when inner spec passes")
	}
	reasons := result.FailureReasons()
	if len(reasons) != 1 || reasons[0] != tooOld {
		t.Errorf("FailureReasons() = %v, want [%v]", reasons, tooOld)
	}
}

func TestNot_Expression(t *testing.T) {
	s := spec.Not("IsMinor", tooOld, isAdult)
	want := "NOT IsAdult"
	if s.Expression() != want {
		t.Errorf("Expression() = %q, want %q", s.Expression(), want)
	}
}

// --- Nesting ---

func TestNested_AnyOfWithAllOf(t *testing.T) {
	financiallyQualified := spec.AnyOf("FinanciallyQualified",
		hasFunds,
		spec.AllOf("VerifiedAndAdult", isVerified, isAdult),
	)

	// passes because hasFunds passes
	result := financiallyQualified.Evaluate(testCtx{Balance: 200})
	if !result.Passed() {
		t.Error("expected pass via hasFunds")
	}

	// passes because isVerified AND isAdult pass
	result = financiallyQualified.Evaluate(testCtx{Age: 20, Verified: true})
	if !result.Passed() {
		t.Error("expected pass via VerifiedAndAdult")
	}

	// fails because neither branch passes
	result = financiallyQualified.Evaluate(testCtx{Age: 16, Balance: 50})
	if result.Passed() {
		t.Error("expected fail when neither branch passes")
	}
}

func TestNested_Expression(t *testing.T) {
	inner := spec.AllOf("Inner", isAdult, hasFunds)
	outer := spec.AnyOf("Outer", isVerified, inner)
	want := "(IsVerified OR (IsAdult AND HasFunds))"
	if outer.Expression() != want {
		t.Errorf("Expression() = %q, want %q", outer.Expression(), want)
	}
}
