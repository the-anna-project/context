package trial

import (
	"testing"

	"github.com/the-anna-project/context"
)

func Test_Disable_Restore(t *testing.T) {
	var val Value
	var ok bool

	ctx := testNewContext(t)
	expected := testNewValue(t)

	// There should be no value in the default context.
	_, ok = FromContext(ctx)
	if ok {
		t.Fatal("expected", false, "got", true)
	}
	ok = IsDisabled(ctx)
	if ok {
		t.Fatal("expected", false, "got", true)
	}

	// Setting the context value should not disable it, but set the context value.
	ctx = NewContext(ctx, expected)
	ok = IsDisabled(ctx)
	if ok {
		t.Fatal("expected", false, "got", true)
	}
	val, ok = FromContext(ctx)
	if !val.Equals(expected) {
		t.Fatal("expected", true, "got", false)
	}

	// Disable the context value should remove it.
	ctx = Disable(ctx)
	ok = IsDisabled(ctx)
	if !ok {
		t.Fatal("expected", true, "got", false)
	}
	_, ok = FromContext(ctx)
	if ok {
		t.Fatal("expected", false, "got", true)
	}

	// Restore the context value should bring it back.
	ctx = Restore(ctx)
	ok = IsDisabled(ctx)
	if ok {
		t.Fatal("expected", false, "got", true)
	}
	val, ok = FromContext(ctx)
	if !val.Equals(expected) {
		t.Fatal("expected", true, "got", false)
	}
}

func Test_NewContextFromContexts(t *testing.T) {
	testCases := []struct {
		Ctx          context.Context
		Ctxs         []context.Context
		ErrorMatcher func(err error) bool
		Expected     Value
	}{
		// Everything is default. Everything should have zero values.
		{
			Ctx:          testNewContext(t),
			Ctxs:         testNewContexts(t),
			ErrorMatcher: nil,
			Expected:     Value{},
		},
		// Given contexts have zero values. Merging should overwrite context value
		// to zero value.
		{
			Ctx: func() context.Context {
				ctx := testNewContext(t)
				ctx = NewContext(ctx, testNewValue(t))
				return ctx
			}(),
			Ctxs:         testNewContexts(t),
			ErrorMatcher: nil,
			Expected:     Value{},
		},
		// Overwriting the zero value of the context with some value should set the
		// context value to this value.
		{
			Ctx: testNewContext(t),
			Ctxs: func() []context.Context {
				ctxs := testNewContexts(t)
				ctxs = testAllContextsWithValue(ctxs, testNewValue(t))
				return ctxs
			}(),
			ErrorMatcher: nil,
			Expected:     testNewValue(t),
		},
		// Overwriting the context value with the context value should not change
		// the context value.
		{
			Ctx: func() context.Context {
				ctx := testNewContext(t)
				ctx = NewContext(ctx, testNewValue(t))
				return ctx
			}(),
			Ctxs: func() []context.Context {
				ctxs := testNewContexts(t)
				ctxs = testAllContextsWithValue(ctxs, testNewValue(t))
				return ctxs
			}(),
			ErrorMatcher: nil,
			Expected:     testNewValue(t),
		},
		// Providing a list of contexts with different values causes the merge to
		// fail.
		{
			Ctx: testNewContext(t),
			Ctxs: func() []context.Context {
				ctxs := testNewContexts(t)
				ctxs = testOneContextWithValue(ctxs, testNewValue(t))
				return ctxs
			}(),
			ErrorMatcher: IsInvalidExecution,
			Expected:     Value{},
		},
	}

	for i, testCase := range testCases {
		ctx, err := NewContextFromContexts(testCase.Ctx, testCase.Ctxs)
		if (err != nil && testCase.ErrorMatcher == nil) || (testCase.ErrorMatcher != nil && !testCase.ErrorMatcher(err)) {
			t.Fatal("case", i+1, "expected", true, "got", false)
		}
		if testCase.ErrorMatcher == nil {
			val, ok := FromContext(ctx)
			if !val.Equals(testCase.Expected) {
				t.Fatal("case", i+1, "expected", testCase.Expected, "got", val)
			}
			if !ok {
				t.Fatal("case", i+1, "expected", true, "got", false)
			}
		}
	}
}

func testAllContextsWithValue(ctxs []context.Context, val Value) []context.Context {
	for i, c := range ctxs {
		ctxs[i] = NewContext(c, val)
	}

	return ctxs
}

func testOneContextWithValue(ctxs []context.Context, val Value) []context.Context {
	for i, c := range ctxs {
		ctxs[i] = NewContext(c, val)
		break
	}

	return ctxs
}

func testNewContext(t *testing.T) context.Context {
	var ctx context.Context
	{
		var err error
		ctx, err = context.New(context.DefaultConfig())
		if err != nil {
			t.Fatal("expected", nil, "got", err)
		}
	}

	return ctx
}

func testNewContexts(t *testing.T) []context.Context {
	var ctxs []context.Context
	{
		for i := 0; i < 3; i++ {
			ctx := testNewContext(t)
			ctxs = append(ctxs, ctx)
		}
	}

	return ctxs
}

func testNewValue(t *testing.T) Value {
	return Value{
		Scope: "scope",
	}
}
