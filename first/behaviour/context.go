// Package behaviour stores and accesses the values defined in this package in
// and from a github.com/the-anna-project/context.Context.
package behaviour

import (
	"reflect"

	"github.com/the-anna-project/context"
)

// key is an unexported type for keys defined in this package. This prevents
// collisions with keys defined in other packages.
type key string

// Value is the context value being managed by this package.
type Value struct {
	// ID represents the ID of the first behaviour.
	ID string `json:"id"`
	// Input reprepresents an extension of the context value.
	Input ValueInput `json:"input"`
	// Name represents the name of the first behaviour.
	Name string `json:"name"`
}

// ValueInput extends the context value.
type ValueInput struct {
	// Types represents the input types of the first behaviour.
	Types []string `json:"types"`
}

// Equals checks whether the properties of the current value equals the
// properties of the given value.
func (v Value) Equals(other Value) bool {
	if v.ID != other.ID {
		return false
	}
	if !reflect.DeepEqual(v.Input.Types, other.Input.Types) {
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
	// behaviour.NewContext and behaviour.FromContext instead of using this key
	// directly.
	valueKey key = "first-behaviour"

	// restoreKey is the key for restoring context values in
	// github.com/the-anna-project/context.Context. Clients use behaviour.Disable
	// and behaviour.Restore instead of using this key directly.
	restoreKey key = "restore-first-behaviour"
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
func FromContext(ctx context.Context) (Value, bool) {
	val, ok := ctx.Value(valueKey).(Value)
	return val, ok
}

// IsDisabled checks whether the given context has the context value removed and
// backed up.
func IsDisabled(ctx context.Context) bool {
	var ok bool

	_, ok = ctx.Value(valueKey).(Value)
	if ok {
		return false
	}
	_, ok = ctx.Value(restoreKey).(Value)
	if !ok {
		return false
	}

	return true
}

// NewContext returns a new github.com/the-anna-project/context.Context that
// carries the context value val.
func NewContext(ctx context.Context, val Value) context.Context {
	ctx.SetValue(valueKey, val)
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
	val, _ := ctx.Value(restoreKey).(Value)
	ctx.SetValue(valueKey, val)
	ctx.DeleteValue(restoreKey)
	return ctx
}
