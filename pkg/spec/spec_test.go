package spec_test

import (
	"testing"

	"github.com/caik/spec2go/pkg/spec"
)

func TestNew_PassesWhenPredicateTrue(t *testing.T) {
	ctx := testCtx{Age: 20}
	result := isAdult.Evaluate(ctx)

	if !result.Passed() {
		t.Errorf("expected pass for age %d", ctx.Age)
	}
}

func TestNew_FailsWhenPredicateFalse(t *testing.T) {
	ctx := testCtx{Age: 16}
	result := isAdult.Evaluate(ctx)

	if result.Passed() {
		t.Errorf("expected fail for age %d", ctx.Age)
	}
	reasons := result.FailureReasons()
	if len(reasons) != 1 || reasons[0] != tooYoung {
		t.Errorf("FailureReasons() = %v, want [%v]", reasons, tooYoung)
	}
}

func TestNew_ResultName(t *testing.T) {
	result := isAdult.Evaluate(testCtx{Age: 20})
	if result.Name() != "IsAdult" {
		t.Errorf("result.Name() = %q, want %q", result.Name(), "IsAdult")
	}
}

func TestNew_SpecName(t *testing.T) {
	if isAdult.Name() != "IsAdult" {
		t.Errorf("spec.Name() = %q, want %q", isAdult.Name(), "IsAdult")
	}
}

func TestNew_Expression(t *testing.T) {
	if isAdult.Expression() != "IsAdult" {
		t.Errorf("spec.Expression() = %q, want %q", isAdult.Expression(), "IsAdult")
	}
}

func TestFunc_CustomExpression(t *testing.T) {
	s := spec.Func[testCtx, testReason]{} // zero value — expression falls back to name
	_ = s                                  // just verifying Func is an exported type

	// Func with explicit expression
	s2 := spec.New("AgeCheck", func(c testCtx) bool { return c.Age >= 18 }, tooYoung)
	if s2.Expression() != "AgeCheck" {
		t.Errorf("Expression() = %q, want %q", s2.Expression(), "AgeCheck")
	}
}

func TestNamedSpec_ImplementsSpecification(t *testing.T) {
	type mySpec struct {
		spec.NamedSpec[testCtx, testReason]
	}
	// NamedSpec must satisfy Name() and Expression() so a custom struct only needs Evaluate.
	n := spec.NamedSpec[testCtx, testReason]{N: "CustomSpec"}
	if n.Name() != "CustomSpec" {
		t.Errorf("Name() = %q, want %q", n.Name(), "CustomSpec")
	}
	if n.Expression() != "CustomSpec" {
		t.Errorf("Expression() = %q, want %q", n.Expression(), "CustomSpec")
	}
}

func TestNew_BoundaryAge(t *testing.T) {
	cases := []struct {
		age    int
		passes bool
	}{
		{17, false},
		{18, true},
		{19, true},
	}
	for _, tc := range cases {
		result := isAdult.Evaluate(testCtx{Age: tc.age})
		if result.Passed() != tc.passes {
			t.Errorf("age %d: Passed() = %v, want %v", tc.age, result.Passed(), tc.passes)
		}
	}
}
