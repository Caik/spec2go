package spec

// Specification is a single evaluable condition.
//
// T is the context type being evaluated.
// R is the comparable failure-reason type (e.g. a typed string constant or iota int).
type Specification[T any, R comparable] interface {
	// Evaluate runs the specification against the given context and returns its result.
	Evaluate(ctx T) SpecificationResult[R]
	// Name returns the human-readable name of this specification.
	Name() string
	// Expression returns a SQL-like representation of this specification (used in policy String()).
	Expression() string
}

// Func is a concrete [Specification] backed by a function.
// Use [New] to create simple specifications, or build a Func directly for more control.
type Func[T any, R comparable] struct {
	name       string
	expression string
	fn         func(T) SpecificationResult[R]
}

// Evaluate implements [Specification].
func (f Func[T, R]) Evaluate(ctx T) SpecificationResult[R] { return f.fn(ctx) }

// Name implements [Specification].
func (f Func[T, R]) Name() string { return f.name }

// Expression implements [Specification]. Falls back to Name if no expression was set.
func (f Func[T, R]) Expression() string {
	if f.expression != "" {
		return f.expression
	}
	return f.name
}

// New creates a [Specification] from a name, a predicate, and a single failure reason.
// The specification passes when predicate returns true, and fails with the given reason otherwise.
func New[T any, R comparable](name string, predicate func(T) bool, reason R) Specification[T, R] {
	return Func[T, R]{
		name: name,
		fn: func(ctx T) SpecificationResult[R] {
			if predicate(ctx) {
				return Pass[R](name)
			}
			return Fail(name, reason)
		},
	}
}

// NamedSpec is an embed helper that satisfies the Name and Expression methods of [Specification].
// Embed it in custom structs to avoid boilerplate — the struct only needs to implement Evaluate.
//
// Example:
//
//	type MySpec struct {
//	    spec.NamedSpec[MyCtx, MyReason]
//	}
//
//	func (s MySpec) Evaluate(ctx MyCtx) spec.SpecificationResult[MyReason] { ... }
//
//	s := MySpec{spec.NamedSpec[MyCtx, MyReason]{N: "MySpec"}}
type NamedSpec[T any, R comparable] struct {
	N string
}

// Name implements [Specification].
func (n NamedSpec[T, R]) Name() string { return n.N }

// Expression implements [Specification].
func (n NamedSpec[T, R]) Expression() string { return n.N }
