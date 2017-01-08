package context

import (
	"encoding/json"
	"time"
)

// Context is a marshallable container used to transport information across
// processes. That is why the interface orients on the native golang context,
// but does not fully replicate it. Reason for this is that a context maps keys
// to values. JSON requires keys of maps to be strings. The native golang
// context defines keys as interface{} to leverage go's type system for unique
// keys. To prevent key collisions from interfering packages, each package
// should prefix its key with its own package path like it would be imported.
// This would be a good context package key, even though it is really long.
//
//     github.com/the-anna-project/context/current/behaviour
//
// This package helps to automatically use the package base path.
//
//     github.com/the-anna-project/gopkg
//
type Context interface {
	Deadline() (time.Time, bool)
	DeleteValue(key string)
	Done() <-chan struct{}
	Err() error
	json.Marshaler
	json.Unmarshaler
	// SetValue stores the given key/value pair within the current context. In
	// case a key is provided that already exists, this key's value will be
	// overwritten with the given one.
	SetValue(key string, value interface{})
	Value(key string) interface{}
}
