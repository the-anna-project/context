// Package destination stores and accesses the values defined in this package in
// and from a github.com/the-anna-project/context.Context.
package destination

import (
	"github.com/the-anna-project/context"
	"github.com/the-anna-project/gopkg"
)

// Value is the context value being managed by this package.
type Value struct {
	// ID represents the ID of the current destination.
	ID string `json:"id"`
	// Name represents the name of the current destination.
	Name string `json:"name"`
}

// Equals checks whether the properties of the current value equals the
// properties of the given value.
func (v Value) Equals(other Value) bool {
	if v.ID != other.ID {
		return false
	}
	if v.Name != other.Name {
		return false
	}

	return true
}

var (
	// valueKey is the key for context values in
	// github.com/the-anna-project/context.Context. Clients use
	// destination.NewContext and destination.FromContext instead of using this
	// key directly.
	valueKey = gopkg.String()

	// restoreKey is the key for restoring context values in
	// github.com/the-anna-project/context.Context. Clients use
	// destination.Disable and destination.Restore instead of using this key
	// directly.
	restoreKey = gopkg.String() + "/restore"
)

// Disable removes the context value being stored using valueKey and backs it up
// using restoreKey.
func Disable(ctx context.Context) context.Context {
	val, _ := FromContext(ctx)
	ctx.Create(restoreKey, val)
	ctx.Delete(valueKey)
	return ctx
}

// FromContext returns the context value stored in ctx, if any.
func FromContext(ctx context.Context) (Value, bool) {
	val, ok := ctx.Search(valueKey).(Value)
	return val, ok
}

// IsDisabled checks whether the given context has the context value removed and
// backed up.
func IsDisabled(ctx context.Context) bool {
	var ok bool

	_, ok = ctx.Search(valueKey).(Value)
	if ok {
		return false
	}
	_, ok = ctx.Search(restoreKey).(Value)
	if !ok {
		return false
	}

	return true
}

// NewContext returns a new github.com/the-anna-project/context.Context that
// carries the context value val.
func NewContext(ctx context.Context, val Value) context.Context {
	ctx.Create(valueKey, val)
	return ctx
}

// NewContextFromContexts sets the context value from the given list of contexts
// to the given single context. Therefore all context values transported by all
// contexts of the given list of contexts have to be equal.
func NewContextFromContexts(ctx context.Context, ctxs []context.Context) (context.Context, error) {
	var reference Value

	for i, c := range ctxs {
		value, _ := FromContext(c)
		if i == 0 {
			reference = value
		}
		if !value.Equals(reference) {
			return nil, maskAnyf(invalidExecutionError, "context values must be equal")
		}
	}

	ctx = NewContext(ctx, reference)

	return ctx, nil
}

// Restore sets the context value using the value being backed up by a previous
// call to Disable.
func Restore(ctx context.Context) context.Context {
	val, _ := ctx.Search(restoreKey).(Value)
	ctx.Create(valueKey, val)
	ctx.Delete(restoreKey)
	return ctx
}
