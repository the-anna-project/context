package context

import (
	"github.com/the-anna-project/context/spec"
)

// Context is a simple redirect to the context specification. That aims to
// provide a simpler interface for clients using these context packages.
type Context spec.Context
