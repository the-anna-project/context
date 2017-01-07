// Package source stores and accesses the values defined in this package in and
// from a github.com/the-anna-project/context.Context.
package source

import (
	"reflect"

	"github.com/the-anna-project/context"
)

// key is an unexported type for keys defined in this package. This prevents
// collisions with keys defined in other packages.
type key string

// Value is the context value being managed by this package.
type Value struct {
	// IDs represents the IDs of the current source.
	IDs []string `json:"ids"`
	// Names represents the names of the current source.
	Names []string `json:"names"`
}

// Equals checks whether the properties of the current value equals the
// properties of the given value.
func (v Value) Equals(other Value) bool {
	if !reflect.DeepEqual(v.IDs, other.IDs) {
		return false
	}
	if !reflect.DeepEqual(v.Names, other.Names) {
		return false
	}

	return true
}

var (
	// valueKey is the key for context values in
	// github.com/the-anna-project/context.Context. Clients use source.NewContext
	// and source.FromContext instead of using this key directly.
	valueKey key = "current-source"

	// restoreKey is the key for restoring context values in
	// github.com/the-anna-project/context.Context. Clients use source.Disable and
	// source.Restore instead of using this key directly.
	restoreKey key = "restore-current-source"
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
	var allIDs []string
	for _, c := range ctxs {
		v, _ := FromContext(c)
		allIDs = append(allIDs, v.IDs...)
	}

	var allNames []string
	for _, c := range ctxs {
		v, _ := FromContext(c)
		allNames = append(allNames, v.Names...)
	}

	val := Value{
		IDs:   allIDs,
		Names: allNames,
	}

	ctx = NewContext(ctx, val)

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
