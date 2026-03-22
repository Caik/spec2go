package spec_test

import (
	"testing"

	"github.com/caik/spec2go/pkg/spec"
)

func TestPass(t *testing.T) {
	r := spec.Pass[testReason]("MySpec")

	if r.Name() != "MySpec" {
		t.Errorf("Name() = %q, want %q", r.Name(), "MySpec")
	}
	if !r.Passed() {
		t.Error("Passed() = false, want true")
	}
	if r.FailureReasons() != nil {
		t.Errorf("FailureReasons() = %v, want nil", r.FailureReasons())
	}
}

func TestFail(t *testing.T) {
	r := spec.Fail("MySpec", tooYoung, insufficientFunds)

	if r.Name() != "MySpec" {
		t.Errorf("Name() = %q, want %q", r.Name(), "MySpec")
	}
	if r.Passed() {
		t.Error("Passed() = true, want false")
	}
	reasons := r.FailureReasons()
	if len(reasons) != 2 {
		t.Fatalf("len(FailureReasons()) = %d, want 2", len(reasons))
	}
	if reasons[0] != tooYoung {
		t.Errorf("FailureReasons()[0] = %v, want %v", reasons[0], tooYoung)
	}
	if reasons[1] != insufficientFunds {
		t.Errorf("FailureReasons()[1] = %v, want %v", reasons[1], insufficientFunds)
	}
}

func TestFailSingleReason(t *testing.T) {
	r := spec.Fail("AgeCheck", tooYoung)

	if r.Passed() {
		t.Error("Passed() = true, want false")
	}
	if len(r.FailureReasons()) != 1 {
		t.Errorf("len(FailureReasons()) = %d, want 1", len(r.FailureReasons()))
	}
}

func TestFailureReasonsReturnsCopy(t *testing.T) {
	r := spec.Fail("AgeCheck", tooYoung)
	reasons := r.FailureReasons()
	reasons[0] = notVerified // mutate the returned slice

	// original should be unaffected
	if r.FailureReasons()[0] != tooYoung {
		t.Error("FailureReasons() did not return a copy — mutation affected original")
	}
}
