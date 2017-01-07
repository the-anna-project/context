package context

import (
	nativecontext "context"
	"encoding/json"
)

type Context interface {
	DeleteValue(key interface{})
	nativecontext.Context
	json.Marshaler
	json.Unmarshaler
	SetValue(key, value interface{})
}
