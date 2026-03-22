package main

import (
	"fmt"

	"github.com/caik/spec2go/pkg/spec"
)

// LoanIneligibilityReason is the set of reasons a loan application may be rejected.
type LoanIneligibilityReason string

const (
	AgeTooYoung        LoanIneligibilityReason = "AGE_TOO_YOUNG"
	InsufficientIncome LoanIneligibilityReason = "INSUFFICIENT_INCOME"
	PoorCreditScore    LoanIneligibilityReason = "POOR_CREDIT_SCORE"
	NotEmployed        LoanIneligibilityReason = "NOT_EMPLOYED"
)

// LoanApplicationContext holds the data evaluated by loan specifications.
type LoanApplicationContext struct {
	Age          int
	AnnualIncome float64
	CreditScore  int
	Employed     bool
}

func main() {
	// --- Define atomic specifications ---
	ageCheck := spec.New("AgeMinimum",
		func(c LoanApplicationContext) bool { return c.Age >= 18 },
		AgeTooYoung,
	)

	incomeCheck := spec.New("SufficientIncome",
		func(c LoanApplicationContext) bool { return c.AnnualIncome >= 30_000 },
		InsufficientIncome,
	)

	creditCheck := spec.New("GoodCreditScore",
		func(c LoanApplicationContext) bool { return c.CreditScore >= 650 },
		PoorCreditScore,
	)

	employmentCheck := spec.New("IsEmployed",
		func(c LoanApplicationContext) bool { return c.Employed },
		NotEmployed,
	)

	// --- Build composite specification ---
	// Applicant is financially qualified if they have good credit OR (sufficient income AND is employed)
	financiallyQualified := spec.AnyOf("FinanciallyQualified",
		creditCheck,
		spec.AllOf("IncomeAndEmployment", incomeCheck, employmentCheck),
	)

	// --- Build policy ---
	loanPolicy := spec.NewPolicy[LoanApplicationContext, LoanIneligibilityReason]().
		With(ageCheck).
		With(financiallyQualified)

	fmt.Printf("Policy structure: %s\n\n", loanPolicy)

	// --- Evaluate applications ---
	applications := []struct {
		name string
		ctx  LoanApplicationContext
	}{
		{"Alice (qualified)", LoanApplicationContext{Age: 30, AnnualIncome: 50_000, CreditScore: 700, Employed: true}},
		{"Bob (too young)", LoanApplicationContext{Age: 17, AnnualIncome: 50_000, CreditScore: 700, Employed: true}},
		{"Carol (poor credit, low income)", LoanApplicationContext{Age: 25, AnnualIncome: 20_000, CreditScore: 500, Employed: false}},
		{"Dave (good credit, unemployed)", LoanApplicationContext{Age: 28, AnnualIncome: 0, CreditScore: 720, Employed: false}},
	}

	for _, app := range applications {
		fmt.Printf("Applicant: %s\n", app.name)

		// Fail-fast: stop at first failure
		result := loanPolicy.EvaluateFailFast(app.ctx)

		if result.AllPassed() {
			fmt.Println("  Decision: APPROVED")
		} else {
			fmt.Printf("  Decision: DENIED (fail-fast reason: %v)\n", result.FailureReasons())
		}

		// Evaluate all: collect every failure
		allResult := loanPolicy.EvaluateAll(app.ctx)

		if !allResult.AllPassed() {
			fmt.Printf("  All failures: %v\n", allResult.FailureReasons())
		}

		fmt.Println()
	}
}
