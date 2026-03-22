package main

import (
	"fmt"

	"github.com/caik/spec2go/pkg/spec"
)

// AccessDeniedReason represents why access to a feature was denied.
type AccessDeniedReason string

const (
	InsufficientRole  AccessDeniedReason = "INSUFFICIENT_ROLE"
	InsufficientPlan  AccessDeniedReason = "INSUFFICIENT_PLAN"
	RateLimitExceeded AccessDeniedReason = "RATE_LIMIT_EXCEEDED"
	RegionNotAllowed  AccessDeniedReason = "REGION_NOT_ALLOWED"
	AccountSuspended  AccessDeniedReason = "ACCOUNT_SUSPENDED"
)

// Role represents user permission levels (higher = more permissions).
type Role int

const (
	RoleGuest Role = iota
	RoleUser
	RoleModerator
	RoleAdmin
)

func (r Role) String() string {
	switch r {
	case RoleGuest:
		return "Guest"
	case RoleUser:
		return "User"
	case RoleModerator:
		return "Moderator"
	case RoleAdmin:
		return "Admin"
	default:
		return "Unknown"
	}
}

// Plan represents subscription tier levels (higher = more features).
type Plan int

const (
	PlanFree Plan = iota
	PlanBasic
	PlanPro
	PlanEnterprise
)

func (p Plan) String() string {
	switch p {
	case PlanFree:
		return "Free"
	case PlanBasic:
		return "Basic"
	case PlanPro:
		return "Pro"
	case PlanEnterprise:
		return "Enterprise"
	default:
		return "Unknown"
	}
}

// AccessContext holds all data needed to evaluate access control rules.
type AccessContext struct {
	Role           Role
	Plan           Plan
	RequestsPerMin int
	Region         string
	Suspended      bool
}

// --- Dynamic specification factories ---

func hasMinimumRole(minRole Role) spec.Specification[AccessContext, AccessDeniedReason] {
	name := fmt.Sprintf("HasMinimumRole(%s)", minRole)

	return spec.New(name, func(c AccessContext) bool {
		return c.Role >= minRole
	}, InsufficientRole)
}

func hasMinimumPlan(minPlan Plan) spec.Specification[AccessContext, AccessDeniedReason] {
	name := fmt.Sprintf("HasMinimumPlan(%s)", minPlan)

	return spec.New(name, func(c AccessContext) bool {
		return c.Plan >= minPlan
	}, InsufficientPlan)
}

func withinRateLimit(maxRPM int) spec.Specification[AccessContext, AccessDeniedReason] {
	name := fmt.Sprintf("WithinRateLimit(%d/min)", maxRPM)

	return spec.New(name, func(c AccessContext) bool {
		return c.RequestsPerMin <= maxRPM
	}, RateLimitExceeded)
}

func allowedRegion(regions ...string) spec.Specification[AccessContext, AccessDeniedReason] {
	allowed := make(map[string]bool, len(regions))

	for _, r := range regions {
		allowed[r] = true
	}

	return spec.New("AllowedRegion", func(c AccessContext) bool {
		return allowed[c.Region]
	}, RegionNotAllowed)
}

// --- Shared specifications ---

var notSuspended = spec.New("NotSuspended",
	func(c AccessContext) bool { return !c.Suspended },
	AccountSuspended,
)

// --- Policy factories ---

func basicAccessPolicy() *spec.Policy[AccessContext, AccessDeniedReason] {
	return spec.NewPolicy[AccessContext, AccessDeniedReason]().
		With(notSuspended).
		With(hasMinimumRole(RoleUser))
}

func apiAccessPolicy(plan Plan) *spec.Policy[AccessContext, AccessDeniedReason] {
	rpmLimit := map[Plan]int{
		PlanFree:       10,
		PlanBasic:      100,
		PlanPro:        1000,
		PlanEnterprise: 10000,
	}[plan]

	return basicAccessPolicy().
		With(hasMinimumPlan(plan)).
		With(withinRateLimit(rpmLimit))
}

func adminAccessPolicy() *spec.Policy[AccessContext, AccessDeniedReason] {
	return basicAccessPolicy().
		With(hasMinimumRole(RoleAdmin)).
		With(allowedRegion("US", "EU", "CA"))
}

func main() {
	users := []struct {
		name string
		ctx  AccessContext
	}{
		{"Alice (Pro user, US)", AccessContext{Role: RoleModerator, Plan: PlanPro, RequestsPerMin: 50, Region: "US", Suspended: false}},
		{"Bob (Free user, over limit)", AccessContext{Role: RoleUser, Plan: PlanFree, RequestsPerMin: 50, Region: "US", Suspended: false}},
		{"Carol (suspended admin)", AccessContext{Role: RoleAdmin, Plan: PlanEnterprise, RequestsPerMin: 5, Region: "EU", Suspended: true}},
		{"Dave (guest, restricted region)", AccessContext{Role: RoleGuest, Plan: PlanFree, RequestsPerMin: 1, Region: "CN", Suspended: false}},
	}

	policies := []struct {
		name   string
		policy *spec.Policy[AccessContext, AccessDeniedReason]
	}{
		{"Basic Access", basicAccessPolicy()},
		{"API Access (Pro plan)", apiAccessPolicy(PlanPro)},
		{"Admin Access", adminAccessPolicy()},
	}

	for _, pol := range policies {
		fmt.Printf("=== %s ===\n", pol.name)
		fmt.Printf("Policy: %s\n", pol.policy)

		for _, user := range users {
			result := pol.policy.EvaluateFailFast(user.ctx)
			status := "GRANTED"
			detail := ""

			if !result.AllPassed() {
				status = "DENIED"
				detail = fmt.Sprintf(" %v", result.FailureReasons())
			}

			fmt.Printf("  %-35s %s%s\n", user.name+":", status, detail)
		}

		fmt.Println()
	}
}
