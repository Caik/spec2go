package spec_test

import "github.com/caik/spec2go/pkg/spec"

// testReason is the failure-reason type used across all tests.
type testReason string

const (
	tooYoung         testReason = "TOO_YOUNG"
	tooOld           testReason = "TOO_OLD"
	insufficientFunds testReason = "INSUFFICIENT_FUNDS"
	notVerified      testReason = "NOT_VERIFIED"
	blocked          testReason = "BLOCKED"
)

// testCtx is the evaluation context used across all tests.
type testCtx struct {
	Age      int
	Balance  float64
	Verified bool
	Blocked  bool
}

// Shared specifications used across test files.
var (
	isAdult    = spec.New("IsAdult", func(c testCtx) bool { return c.Age >= 18 }, tooYoung)
	isNotTooOld = spec.New("IsNotTooOld", func(c testCtx) bool { return c.Age <= 65 }, tooOld)
	hasFunds   = spec.New("HasFunds", func(c testCtx) bool { return c.Balance >= 100 }, insufficientFunds)
	isVerified = spec.New("IsVerified", func(c testCtx) bool { return c.Verified }, notVerified)
	isNotBlocked = spec.New("IsNotBlocked", func(c testCtx) bool { return !c.Blocked }, blocked)
	alwaysPass = spec.New("AlwaysPass", func(c testCtx) bool { return true }, blocked)
	alwaysFail = spec.New("AlwaysFail", func(c testCtx) bool { return false }, blocked)
)
