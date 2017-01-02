// Package expectation stores and accesses a expectation in and from a
// github.com/the-anna-project/context/spec.Context.
package expectation

import (
	"github.com/the-anna-project/context/spec"
	"github.com/the-anna-project/expectation"
)

// key is an unexported type for keys defined in this package. This prevents
// collisions with keys defined in other packages.
type key string

// expectationKey is the key for expectation values in
// github.com/the-anna-project/context/spec.Context. Clients use
// expectation.NewContext and expectation.FromContext instead of using this key
// directly.
var expectationKey key = "expectation"

// NewContext returns a new github.com/the-anna-project/context/spec.Context that
// carries value v.
func NewContext(ctx spec.Context, v expectation.Expectation) spec.Context {
	if v == nil {
		// In case the given value is empty we do not add it, but only return the
		// given context as it is. That way the existence check when reading the
		// context works as expected when no value or an empty value was tried to be
		// added.
		return ctx
	}

	ctx.SetValue(expectationKey, v)

	return ctx
}

// FromContext returns the expectation value stored in ctx, if any.
func FromContext(ctx spec.Context) (expectation.Expectation, bool) {
	v, ok := ctx.Value(expectationKey).(expectation.Expectation)
	return v, ok
}
