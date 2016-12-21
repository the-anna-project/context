// Package expectation stores and accesses a expectation in and from a
// github.com/the-anna-project/context.Context.
package expectation

import (
	"github.com/the-anna-project/expectation"

	"github.com/the-anna-project/context"
)

// key is an unexported type for keys defined in this package. This prevents
// collisions with keys defined in other packages.
type key string

// expectationKey is the key for expectation values in
// github.com/the-anna-project/context.Context. Clients use
// expectation.NewContext and expectation.FromContext instead of using this key
// directly.
var expectationKey key = "expectation"

// NewContext returns a new github.com/the-anna-project/context.Context that
// carries value v.
func NewContext(ctx context.Context, v expectation.Expectation) context.Context {
	if v == nil {
		// In case the given value is empty we do not add it, but only return the
		// given context as it is. That way the existence check when reading the
		// context works as expected when no value or an empty value was tried to be
		// added.
		return ctx
	}

	return context.WithValue(ctx, expectationKey, v)
}

// FromContext returns the expectation value stored in ctx, if any.
func FromContext(ctx context.Context) (expectation.Expectation, bool) {
	v, ok := ctx.Value(expectationKey).(expectation.Expectation)
	return v, ok
}
