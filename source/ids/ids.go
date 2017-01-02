// Package ids stores and accesses a source IDs in and from a
// github.com/the-anna-project/context/spec.Context.
package ids

import (
	"github.com/the-anna-project/context/spec"
)

// key is an unexported type for keys defined in this package. This prevents
// collisions with keys defined in other packages.
type key string

// idsKey is the key for source IDs values in
// github.com/the-anna-project/context/spec.Context. Clients use ids.NewContext and
// ids.FromContext instead of using this key directly.
var idsKey key = "source-ids"

// NewContext returns a new github.com/the-anna-project/context/spec.Context that
// carries value v.
func NewContext(ctx spec.Context, v []string) spec.Context {
	if len(v) == 0 {
		// In case the given value is empty we do not add it, but only return the
		// given context as it is. That way the existence check when reading the
		// context works as expected when no value or an empty value was tried to be
		// added.
		return ctx
	}

	ctx.SetValue(idsKey, v)

	return ctx
}

// FromContext returns the source IDs value stored in ctx, if any.
func FromContext(ctx spec.Context) ([]string, bool) {
	v, ok := ctx.Value(idsKey).([]string)
	return v, ok
}
