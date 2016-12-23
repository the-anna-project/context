// Package ids stores and accesses a source IDs in and from a
// github.com/the-anna-project/context.Context.
package ids

import (
	"github.com/the-anna-project/context"
)

// key is an unexported type for keys defined in this package. This prevents
// collisions with keys defined in other packages.
type key string

// idsKey is the key for source IDs values in
// github.com/the-anna-project/context.Context. Clients use ids.NewContext and
// ids.FromContext instead of using this key directly.
var idsKey key = "source-ids"

// NewContext returns a new github.com/the-anna-project/context.Context that
// carries value v.
func NewContext(ctx context.Context, v []string) context.Context {
	if v == nil {
		// In case the given value is empty we do not add it, but only return the
		// given context as it is. That way the existence check when reading the
		// context works as expected when no value or an empty value was tried to be
		// added.
		return ctx
	}

	return context.WithValue(ctx, idsKey, v)
}

// FromContext returns the source IDs value stored in ctx, if any.
func FromContext(ctx context.Context) ([]string, bool) {
	v, ok := ctx.Value(idsKey).([]string)
	return v, ok
}
