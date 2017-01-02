// Package types stores and accesses the current behaviour's input types in and
// from a github.com/the-anna-project/context/spec.Context.
package types

import (
	"github.com/the-anna-project/context/spec"
)

// key is an unexported type for keys defined in this package. This prevents
// collisions with keys defined in other packages.
type key string

// typesKey is the key for the current behaviour's input types values in
// github.com/the-anna-project/context/spec.Context. Clients use
// types.NewContext and types.FromContext instead of using this key directly.
var typesKey key = "current-behaviour-input-types"

// NewContext returns a new github.com/the-anna-project/context/spec.Context
// that carries value v.
func NewContext(ctx spec.Context, v []string) spec.Context {
	if len(v) == 0 {
		// In case the given value is empty we do not add it, but only return the
		// given context as it is. That way the existence check when reading the
		// context works as expected when no value or an empty value was tried to be
		// added.
		return ctx
	}

	ctx.SetValue(typesKey, v)

	return ctx
}

// FromContext returns the current behaviour types value stored in ctx, if any.
func FromContext(ctx spec.Context) ([]string, bool) {
	v, ok := ctx.Value(typesKey).([]string)
	return v, ok
}
