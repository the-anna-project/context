// Package name stores and accesses the current behaviour name in and from a
// github.com/the-anna-project/context.Context.
package name

import (
	"github.com/the-anna-project/context"
)

// key is an unexported type for keys defined in this package. This prevents
// collisions with keys defined in other packages.
type key string

// nameKey is the key for the current behaviour name values in
// github.com/the-anna-project/context.Context. Clients use name.NewContext and
// name.FromContext instead of using this key directly.
var nameKey key = "current-behaviour-name"

// NewContext returns a new github.com/the-anna-project/context.Context that
// carries value v.
func NewContext(ctx context.Context, v string) context.Context {
	if v == "" {
		// In case the given value is empty we do not add it, but only return the
		// given context as it is. That way the existence check when reading the
		// context works as expected when no value or an empty value was tried to be
		// added.
		return ctx
	}

	return context.WithValue(ctx, nameKey, v)
}

// FromContext returns the current behaviour name value stored in ctx, if any.
func FromContext(ctx context.Context) (string, bool) {
	v, ok := ctx.Value(nameKey).(string)
	return v, ok
}
