package merge

import (
	"github.com/the-anna-project/context"
	currentbehaviour "github.com/the-anna-project/context/current/behaviour"
	currentclgtree "github.com/the-anna-project/context/current/clg/tree"
	currentdestination "github.com/the-anna-project/context/current/destination"
	currentexpectation "github.com/the-anna-project/context/current/expectation"
	currentsession "github.com/the-anna-project/context/current/session"
	currentsource "github.com/the-anna-project/context/current/source"
	currentstage "github.com/the-anna-project/context/current/stage"
	firstbehaviour "github.com/the-anna-project/context/first/behaviour"
	firstinformation "github.com/the-anna-project/context/first/information"
)

// NewContextFromContexts creates a new context from the given list of contexts.
// Therefore all information transported by all the given contexts have to be
// equal, but the following ones.
//
//     current/source
//
func NewContextFromContexts(ctx context.Context, ctxs []context.Context) (context.Context, error) {
	var err error

	modifiers := []func(ctx context.Context, ctxs []context.Context) (context.Context, error){
		currentbehaviour.NewContextFromContexts,
		currentclgtree.NewContextFromContexts,
		currentdestination.NewContextFromContexts,
		currentexpectation.NewContextFromContexts,
		currentsession.NewContextFromContexts,
		currentsource.NewContextFromContexts,
		currentstage.NewContextFromContexts,
		firstbehaviour.NewContextFromContexts,
		firstinformation.NewContextFromContexts,
	}

	for _, m := range modifiers {
		ctx, err = m(ctx, ctxs)
		if err != nil {
			return nil, maskAny(err)
		}
	}

	return ctx, nil
}
