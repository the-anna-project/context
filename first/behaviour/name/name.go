// Package name stores and accesses the first behaviour name in and from a
// github.com/the-anna-project/context/spec.Context.
package name

import (
	"github.com/the-anna-project/context/spec"
)

// key is an unexported type for keys defined in this package. This prevents
// collisions with keys defined in other packages.
type key string

// nameKey is the key for the first behaviour name values in
// github.com/the-anna-project/context/spec.Context. Clients use name.NewContext and
// name.FromContext instead of using this key directly.
var nameKey key = "first-behaviour-name"

// NewContext returns a new github.com/the-anna-project/context/spec.Context that
// carries value v.
func NewContext(ctx spec.Context, v string) spec.Context {
	if v == "" {
		// In case the given value is empty we do not add it, but only return the
		// given context as it is. That way the existence check when reading the
		// context works as expected when no value or an empty value was tried to be
		// added.
		return ctx
	}

	ctx.SetValue(nameKey, v)

	return ctx
}

// FromContext returns the first behaviour name value stored in ctx, if any.
func FromContext(ctx spec.Context) (string, bool) {
	v, ok := ctx.Value(nameKey).(string)
	return v, ok
}
