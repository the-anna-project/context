// Package expectation stores and accesses the values defined in this package in
// and from a github.com/the-anna-project/context.Context.
package expectation

import (
	"github.com/the-anna-project/context"
	"github.com/the-anna-project/expectation"
)

// key is an unexported type for keys defined in this package. This prevents
// collisions with keys defined in other packages.
type key string

var (
	// valueKey is the key for context values in
	// github.com/the-anna-project/context.Context. Clients use
	// expectation.NewContext and expectation.FromContext instead of using this
	// key directly.
	valueKey key = "current-expectation"

	// restoreKey is the key for restoring context values in
	// github.com/the-anna-project/context.Context. Clients use
	// expectation.Disable and expectation.Restore instead of using this key
	// directly.
	restoreKey key = "restore-current-expectation"
)

// Disable removes the context value being stored using valueKey and backs it up
// using restoreKey.
func Disable(ctx context.Context) context.Context {
	val, _ := FromContext(ctx)
	ctx.SetValue(restoreKey, val)
	ctx.DeleteValue(valueKey)
	return ctx
}

// FromContext returns the context value stored in ctx, if any.
func FromContext(ctx context.Context) (expectation.Expectation, bool) {
	val, ok := ctx.Value(valueKey).(expectation.Expectation)
	return val, ok
}

// IsDisabled checks whether the given context has the context value removed and
// backed up.
func IsDisabled(ctx context.Context) bool {
	var ok bool

	_, ok = ctx.Value(valueKey).(expectation.Expectation)
	if ok {
		return false
	}
	_, ok = ctx.Value(restoreKey).(expectation.Expectation)
	if !ok {
		return false
	}

	return true
}

// NewContext returns a new github.com/the-anna-project/context.Context that
// carries the context value val.
func NewContext(ctx context.Context, val expectation.Expectation) context.Context {
	ctx.SetValue(valueKey, val)
	return ctx
}

// NewContextFromContexts sets the context value from the given list of contexts
// to the given single context. Therefore all context values transported by all
// contexts of the given list of contexts have to be equal.
func NewContextFromContexts(ctx context.Context, ctxs []context.Context) (context.Context, error) {
	var reference expectation.Expectation

	for i, c := range ctxs {
		value, _ := FromContext(c)
		if i == 0 {
			reference = value
		}
		if value == nil && reference != nil {
			return nil, maskAnyf(invalidExecutionError, "context values must be equal")
		}
		if value != nil && !value.Equals(reference) {
			return nil, maskAnyf(invalidExecutionError, "context values must be equal")
		}
	}

	ctx = NewContext(ctx, reference)

	return ctx, nil
}

// Restore sets the context value using the value being backed up by a previous
// call to Disable.
func Restore(ctx context.Context) context.Context {
	val, _ := ctx.Value(restoreKey).(expectation.Expectation)
	ctx.SetValue(valueKey, val)
	ctx.DeleteValue(restoreKey)
	return ctx
}
