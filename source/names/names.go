// Package names stores and accesses a source names in and from a
// github.com/the-anna-project/context.Context.
package names

import (
	"github.com/the-anna-project/context"
)

// key is an unexported type for keys defined in this package. This prevents
// collisions with keys defined in other packages.
type key string

// namesKey is the key for source names values in
// github.com/the-anna-project/context.Context. Clients use names.NewContext and
// names.FromContext instead of using this key directly.
var namesKey key = "source-names"

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

	return context.WithValue(ctx, namesKey, v)
}

// FromContext returns the source names value stored in ctx, if any.
func FromContext(ctx context.Context) ([]string, bool) {
	v, ok := ctx.Value(namesKey).([]string)
	return v, ok
}
