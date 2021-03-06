// Package context implements golang.org/x/net/context.Context and provides
// marshallable context primitives to distribute information across event
// queues.
package context

import (
	nativecontext "context"
	"encoding/json"
	"sync"
	"time"
)

// Config represents the configuration used to create a new context.
type Config struct {
	// Settings.
	Context nativecontext.Context
}

// DefaultConfig provides a default configuration to create a new context by
// best effort.
func DefaultConfig() Config {
	newConfig := Config{
		// Settings.
		Context: nativecontext.Background(),
	}

	return newConfig
}

// New creates a new configured context object.
func New(config Config) (Context, error) {
	// Settings.
	if config.Context == nil {
		return nil, maskAnyf(invalidConfigError, "context must not be empty")
	}

	ctx, cancelFunc := nativecontext.WithCancel(config.Context)

	newContext := &context{
		// Internals.
		CancelFunc: cancelFunc,
		CancelOnce: sync.Once{},
		Context:    ctx,
		Storage:    map[string]interface{}{},
	}

	return newContext, nil
}

type context struct {
	// Internals.
	CancelFunc func()                 `json:"-"`
	CancelOnce sync.Once              `json:"-"`
	Context    nativecontext.Context  `json:"-"`
	Storage    map[string]interface{} `json:"storage"`
}

func (c *context) Cancel() {
	c.CancelOnce.Do(func() {
		c.CancelFunc()
	})
}

func (c *context) Clone() (Context, error) {
	newContext, err := New(DefaultConfig())
	if err != nil {
		return nil, maskAny(err)
	}

	for k, v := range c.Storage {
		newContext.(*context).Storage[k] = v
	}

	return newContext, nil
}

func (c *context) Create(key string, value interface{}) {
	c.Storage[key] = value
}

func (c *context) Deadline() (time.Time, bool) {
	return c.Context.Deadline()
}

func (c *context) Delete(key string) {
	delete(c.Storage, key)
}

func (c *context) Done() <-chan struct{} {
	return c.Context.Done()
}

func (c *context) Err() error {
	return c.Context.Err()
}

func (c *context) MarshalJSON() ([]byte, error) {
	type ContextClone context

	b, err := json.Marshal(&struct {
		*ContextClone
	}{
		ContextClone: (*ContextClone)(c),
	})
	if err != nil {
		return nil, maskAny(err)
	}

	return b, nil
}

func (c *context) UnmarshalJSON(b []byte) error {
	type ContextClone context

	aux := &struct {
		*ContextClone
	}{
		ContextClone: (*ContextClone)(c),
	}
	err := json.Unmarshal(b, &aux)
	if err != nil {
		return maskAny(err)
	}

	return nil
}

func (c *context) Search(key string) interface{} {
	v, ok := c.Storage[key]
	if ok {
		return v
	}

	return nil
}
