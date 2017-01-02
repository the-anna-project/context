// Package id stores and accesses a destination ID in and from a
// github.com/the-anna-project/context/spec.Context.
package id

import (
	"github.com/the-anna-project/context/spec"
)

// key is an unexported type for keys defined in this package. This prevents
// collisions with keys defined in other packages.
type key string

// idKey is the key for destination ID values in
// github.com/the-anna-project/context/spec.Context. Clients use id.NewContext and
// id.FromContext instead of using this key directly.
var idKey key = "destination-id"

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

	ctx.SetValue(idKey, v)

	return ctx
}

// FromContext returns the destination ID value stored in ctx, if any.
func FromContext(ctx spec.Context) (string, bool) {
	v, ok := ctx.Value(idKey).(string)
	return v, ok
}
