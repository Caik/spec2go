package spec_test

import (
	"testing"

	"github.com/caik/spec2go/pkg/spec"
)

// --- AnyOf ---

func TestAnyOf_PassesWhenFirstSpecPasses(t *testing.T) {
	s, err := spec.AnyOf("combo", isAdult, hasFunds)
	if err != nil {
		t.Fatal(err)
	}
	result := s.Evaluate(testCtx{Age: 20, Balance: 0})
	if !result.Passed() {
		t.Error("expected pass when first spec passes")
	}
}

func TestAnyOf_PassesWhenSecondSpecPasses(t *testing.T) {
	s, err := spec.AnyOf("combo", isAdult, hasFunds)
	if err != nil {
		t.Fatal(err)
	}
	result := s.Evaluate(testCtx{Age: 16, Balance: 200})
	if !result.Passed() {
		t.Error("expected pass when second spec passes")
	}
}

func TestAnyOf_FailsWhenAllSpecsFail(t *testing.T) {
	s, err := spec.AnyOf("combo", isAdult, hasFunds)
	if err != nil {
		t.Fatal(err)
	}
	result := s.Evaluate(testCtx{Age: 16, Balance: 50})
	if result.Passed() {
		t.Error("expected fail when all specs fail")
	}
}

func TestAnyOf_CollectsAllReasonsOnFailure(t *testing.T) {
	s, err := spec.AnyOf("combo", isAdult, hasFunds)
	if err != nil {
		t.Fatal(err)
	}
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

	s, err := spec.AnyOf("combo", counter, never)
	if err != nil {
		t.Fatal(err)
	}
	s.Evaluate(testCtx{})

	if evalCount != 1 {
		t.Errorf("evalCount = %d, want 1 (short-circuit)", evalCount)
	}
}

func TestAnyOf_Expression(t *testing.T) {
	s, err := spec.AnyOf("combo", isAdult, hasFunds)
	if err != nil {
		t.Fatal(err)
	}
	want := "(IsAdult OR HasFunds)"
	if s.Expression() != want {
		t.Errorf("Expression() = %q, want %q", s.Expression(), want)
	}
}

func TestAnyOf_SingleSpec(t *testing.T) {
	s, err := spec.AnyOf("single", isAdult)
	if err != nil {
		t.Fatal(err)
	}
	if s.Expression() != "IsAdult" {
		t.Errorf("Expression() = %q, want %q", s.Expression(), "IsAdult")
	}
}

func TestAnyOf_ErrorOnEmpty(t *testing.T) {
	_, err := spec.AnyOf[testCtx, testReason]("empty")
	if err == nil {
		t.Error("expected error for empty specs")
	}
}

// --- AnyOfAll ---

func TestAnyOfAll_EvaluatesAllSpecs(t *testing.T) {
	evalCount := 0
	s1 := spec.New("S1", func(c testCtx) bool { evalCount++; return true }, blocked)
	s2 := spec.New("S2", func(c testCtx) bool { evalCount++; return false }, blocked)

	s, err := spec.AnyOfAll("combo", s1, s2)
	if err != nil {
		t.Fatal(err)
	}
	result := s.Evaluate(testCtx{})

	if evalCount != 2 {
		t.Errorf("evalCount = %d, want 2 (no short-circuit)", evalCount)
	}
	if !result.Passed() {
		t.Error("expected pass because s1 passed")
	}
}

func TestAnyOfAll_FailsWhenAllFail(t *testing.T) {
	s, err := spec.AnyOfAll("combo", alwaysFail, alwaysFail)
	if err != nil {
		t.Fatal(err)
	}
	if s.Evaluate(testCtx{}).Passed() {
		t.Error("expected fail when all specs fail")
	}
}

func TestAnyOfAll_ErrorOnEmpty(t *testing.T) {
	_, err := spec.AnyOfAll[testCtx, testReason]("empty")
	if err == nil {
		t.Error("expected error for empty specs")
	}
}

func TestAnyOfAll_CollectsAllReasonsOnFailure(t *testing.T) {
	s, err := spec.AnyOfAll("combo", isAdult, hasFunds)
	if err != nil {
		t.Fatal(err)
	}
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
	s, err := spec.AllOf("combo", isAdult, hasFunds)
	if err != nil {
		t.Fatal(err)
	}
	result := s.Evaluate(testCtx{Age: 20, Balance: 200})
	if !result.Passed() {
		t.Error("expected pass when all specs pass")
	}
}

func TestAllOf_FailsWhenFirstFails(t *testing.T) {
	s, err := spec.AllOf("combo", isAdult, hasFunds)
	if err != nil {
		t.Fatal(err)
	}
	result := s.Evaluate(testCtx{Age: 16, Balance: 200})
	if result.Passed() {
		t.Error("expected fail when first spec fails")
	}
}

func TestAllOf_CollectsAllFailureReasons(t *testing.T) {
	s, err := spec.AllOf("combo", isAdult, hasFunds)
	if err != nil {
		t.Fatal(err)
	}
	result := s.Evaluate(testCtx{Age: 16, Balance: 50})
	if len(result.FailureReasons()) != 2 {
		t.Errorf("FailureReasons() len = %d, want 2", len(result.FailureReasons()))
	}
}

func TestAllOf_EvaluatesAllSpecsEvenOnFailure(t *testing.T) {
	evalCount := 0
	s1 := spec.New("S1", func(c testCtx) bool { evalCount++; return false }, tooYoung)
	s2 := spec.New("S2", func(c testCtx) bool { evalCount++; return false }, insufficientFunds)

	s, err := spec.AllOf("combo", s1, s2)
	if err != nil {
		t.Fatal(err)
	}
	s.Evaluate(testCtx{})

	if evalCount != 2 {
		t.Errorf("evalCount = %d, want 2", evalCount)
	}
}

func TestAllOf_Expression(t *testing.T) {
	s, err := spec.AllOf("combo", isAdult, hasFunds)
	if err != nil {
		t.Fatal(err)
	}
	want := "(IsAdult AND HasFunds)"
	if s.Expression() != want {
		t.Errorf("Expression() = %q, want %q", s.Expression(), want)
	}
}

func TestAllOf_ErrorOnEmpty(t *testing.T) {
	_, err := spec.AllOf[testCtx, testReason]("empty")
	if err == nil {
		t.Error("expected error for empty specs")
	}
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
	verifiedAndAdult, err := spec.AllOf("VerifiedAndAdult", isVerified, isAdult)
	if err != nil {
		t.Fatal(err)
	}
	financiallyQualified, err := spec.AnyOf("FinanciallyQualified", hasFunds, verifiedAndAdult)
	if err != nil {
		t.Fatal(err)
	}

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
	inner, err := spec.AllOf("Inner", isAdult, hasFunds)
	if err != nil {
		t.Fatal(err)
	}
	outer, err := spec.AnyOf("Outer", isVerified, inner)
	if err != nil {
		t.Fatal(err)
	}
	want := "(IsVerified OR (IsAdult AND HasFunds))"
	if outer.Expression() != want {
		t.Errorf("Expression() = %q, want %q", outer.Expression(), want)
	}
}
