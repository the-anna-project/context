// Package context implements golang.org/x/net/context.Context and provides
// marshallable context primitives to distribute information across event
// queues.
package context

import (
	nativecontext "context"
	"encoding/json"
	"reflect"
	"time"

	clgtreeid "github.com/the-anna-project/context/clg/tree/id"
	currentbehaviourid "github.com/the-anna-project/context/current/behaviour/id"
	currentbehaviourinputtypes "github.com/the-anna-project/context/current/behaviour/input/types"
	currentbehaviourname "github.com/the-anna-project/context/current/behaviour/name"
	destinationid "github.com/the-anna-project/context/destination/id"
	destinationname "github.com/the-anna-project/context/destination/name"
	expectationcontext "github.com/the-anna-project/context/expectation"
	firstbehaviourid "github.com/the-anna-project/context/first/behaviour/id"
	firstbehaviourname "github.com/the-anna-project/context/first/behaviour/name"
	firstinformationid "github.com/the-anna-project/context/first/information/id"
	sessionid "github.com/the-anna-project/context/session/id"
	sourceids "github.com/the-anna-project/context/source/ids"
	sourcenames "github.com/the-anna-project/context/source/names"
	"github.com/the-anna-project/expectation"
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

	newContext := &context{
		// Internals.
		storage: map[interface{}]interface{}{},

		// Settings.
		context: config.Context,
	}

	return newContext, nil
}

// NewFromContexts creates a new context from the given list of contexts.
// Therefore all information transported by all the gigen contexts have to be
// equal, but the following one.
//
//     source/ids
//     source/names
//
func NewFromContexts(ctxs []Context) (Context, error) {
	var clgTreeIDRef string
	for i, ctx := range ctxs {
		clgTreeID, ok := clgtreeid.FromContext(ctx)
		if !ok {
			return nil, maskAnyf(invalidContextError, "clg tree id must not be empty")
		}
		if i == 0 {
			clgTreeIDRef = clgTreeID
		}
		if clgTreeID != clgTreeIDRef {
			return nil, maskAnyf(invalidContextError, "clg tree ids must be equal")
		}
	}

	var currentBehaviourIDRef string
	for i, ctx := range ctxs {
		currentBehaviourID, ok := currentbehaviourid.FromContext(ctx)
		if !ok {
			return nil, maskAnyf(invalidContextError, "current behaviour id must not be empty")
		}
		if i == 0 {
			currentBehaviourIDRef = currentBehaviourID
		}
		if currentBehaviourID != currentBehaviourIDRef {
			return nil, maskAnyf(invalidContextError, "current behaviour ids must be equal")
		}
	}

	var currentBehaviourInputTypesRef []string
	for i, ctx := range ctxs {
		currentBehaviourInputTypes, ok := currentbehaviourinputtypes.FromContext(ctx)
		if !ok {
			return nil, maskAnyf(invalidContextError, "current behaviour input types must not be empty")
		}
		if i == 0 {
			currentBehaviourInputTypesRef = currentBehaviourInputTypes
		}
		if !reflect.DeepEqual(currentBehaviourInputTypes, currentBehaviourInputTypesRef) {
			return nil, maskAnyf(invalidContextError, "current behaviour input types must be equal")
		}
	}

	var currentBehaviourNameRef string
	for i, ctx := range ctxs {
		currentBehaviourName, ok := currentbehaviourname.FromContext(ctx)
		if !ok {
			return nil, maskAnyf(invalidContextError, "current behaviour name must not be empty")
		}
		if i == 0 {
			currentBehaviourNameRef = currentBehaviourName
		}
		if currentBehaviourName != currentBehaviourNameRef {
			return nil, maskAnyf(invalidContextError, "current behaviour names must be equal")
		}
	}

	var destinationIDRef string
	for i, ctx := range ctxs {
		destinationID, ok := destinationid.FromContext(ctx)
		if !ok {
			return nil, maskAnyf(invalidContextError, "destination id must not be empty")
		}
		if i == 0 {
			destinationIDRef = destinationID
		}
		if destinationID != destinationIDRef {
			return nil, maskAnyf(invalidContextError, "destination ids must be equal")
		}
	}

	var destinationNameRef string
	for i, ctx := range ctxs {
		destinationName, ok := destinationname.FromContext(ctx)
		if !ok {
			return nil, maskAnyf(invalidContextError, "destination name must not be empty")
		}
		if i == 0 {
			destinationNameRef = destinationName
		}
		if destinationName != destinationNameRef {
			return nil, maskAnyf(invalidContextError, "destination names must be equal")
		}
	}

	var expectationRef expectation.Expectation
	for i, ctx := range ctxs {
		e, _ := expectationcontext.FromContext(ctx)
		// We do not throw an error in case there is no expectation, because there
		// does not have to be any expectation. The rule is only that if there is
		// any expectation, all have to be the same. That also means that if there
		// is no expectation, all contexts must not have any expectation.
		if i == 0 {
			expectationRef = e
		}
		if !e.Equals(expectationRef) {
			return nil, maskAnyf(invalidContextError, "expectations must be equal")
		}
	}

	var firstBehaviourIDRef string
	for i, ctx := range ctxs {
		firstBehaviourID, ok := firstbehaviourid.FromContext(ctx)
		if !ok {
			return nil, maskAnyf(invalidContextError, "first behaviour id must not be empty")
		}
		if i == 0 {
			firstBehaviourIDRef = firstBehaviourID
		}
		if firstBehaviourID != firstBehaviourIDRef {
			return nil, maskAnyf(invalidContextError, "first behaviour ids must be equal")
		}
	}

	var firstBehaviourNameRef string
	for i, ctx := range ctxs {
		firstBehaviourName, ok := firstbehaviourname.FromContext(ctx)
		if !ok {
			return nil, maskAnyf(invalidContextError, "first behaviour name must not be empty")
		}
		if i == 0 {
			firstBehaviourNameRef = firstBehaviourName
		}
		if firstBehaviourName != firstBehaviourNameRef {
			return nil, maskAnyf(invalidContextError, "first behaviour names must be equal")
		}
	}

	var firstInformationIDRef string
	for i, ctx := range ctxs {
		firstInformationID, ok := firstinformationid.FromContext(ctx)
		if !ok {
			return nil, maskAnyf(invalidContextError, "first information id must not be empty")
		}
		if i == 0 {
			firstInformationIDRef = firstInformationID
		}
		if firstInformationID != firstInformationIDRef {
			return nil, maskAnyf(invalidContextError, "first information ids must be equal")
		}
	}

	var sessionIDRef string
	for i, ctx := range ctxs {
		sessionID, ok := sessionid.FromContext(ctx)
		if !ok {
			return nil, maskAnyf(invalidContextError, "source ids must not be empty")
		}
		if i == 0 {
			sessionIDRef = sessionID
		}
		if sessionID != sessionIDRef {
			return nil, maskAnyf(invalidContextError, "session ids must be equal")
		}
	}

	var allSourceIDs []string
	for _, ctx := range ctxs {
		sourceIDs, ok := sourceids.FromContext(ctx)
		if !ok {
			return nil, maskAnyf(invalidContextError, "source ids must not be empty")
		}
		allSourceIDs = append(allSourceIDs, sourceIDs...)
	}

	var allSourceNames []string
	for _, ctx := range ctxs {
		sourceNames, ok := sourceids.FromContext(ctx)
		if !ok {
			return nil, maskAnyf(invalidContextError, "source names must not be empty")
		}
		allSourceNames = append(allSourceNames, sourceNames...)
	}

	config := DefaultConfig()
	newCtx, err := New(config)
	if err != nil {
		return nil, maskAny(err)
	}

	newCtx = clgtreeid.NewContext(newCtx, clgTreeIDRef)
	newCtx = currentbehaviourid.NewContext(newCtx, currentBehaviourIDRef)
	newCtx = currentbehaviourinputtypes.NewContext(newCtx, currentBehaviourInputTypesRef)
	newCtx = currentbehaviourname.NewContext(newCtx, currentBehaviourNameRef)
	newCtx = destinationid.NewContext(newCtx, destinationIDRef)
	newCtx = destinationname.NewContext(newCtx, destinationNameRef)
	newCtx = expectationcontext.NewContext(newCtx, expectationRef)
	newCtx = firstbehaviourid.NewContext(newCtx, firstBehaviourIDRef)
	newCtx = firstbehaviourname.NewContext(newCtx, firstBehaviourNameRef)
	newCtx = firstinformationid.NewContext(newCtx, firstInformationIDRef)
	newCtx = sessionid.NewContext(newCtx, sessionIDRef)
	newCtx = sourceids.NewContext(newCtx, allSourceIDs)
	newCtx = sourcenames.NewContext(newCtx, allSourceNames)

	return newCtx, nil
}

type context struct {
	// Internals.
	storage map[interface{}]interface{}

	// Settings.
	context nativecontext.Context
}

func (c *context) Deadline() (time.Time, bool) {
	return c.context.Deadline()
}

func (c *context) Done() <-chan struct{} {
	return c.context.Done()
}

func (c *context) Err() error {
	return c.context.Err()
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

func (c *context) SetValue(key, value interface{}) {
	c.storage[key] = value
}

func (c *context) Value(key interface{}) interface{} {
	v, ok := c.storage[key]
	if ok {
		return v
	}

	return nil
}
