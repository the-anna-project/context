package context

import (
	"bytes"
	nativecontext "context"
	"encoding/json"
	"testing"
	"time"
)

func Test_JSON(t *testing.T) {
	// Create new context and attach some information to it.
	ctx := testNewContext(t)

	// Marshal and unmarshal the context.
	b, err := json.Marshal(ctx)
	if err != nil {
		t.Fatal("expected", nil, "got", err)
	}
	other, err := New(DefaultConfig())
	if err != nil {
		t.Fatal("expected", nil, "got", err)
	}
	err = json.Unmarshal(b, other)
	if err != nil {
		t.Fatal("expected", nil, "got", err)
	}

	// Check the contexts are equal.
	b1, err := json.Marshal(ctx)
	if err != nil {
		panic(err)
	}
	b2, err := json.Marshal(other)
	if err != nil {
		panic(err)
	}
	r := bytes.Compare(b1, b2)
	if r != 0 {
		t.Fatal("expected", 0, "got", r)
	}

	// Verify the values are actually the same.
	v, ok := ctx.(*context).Storage["foo"]
	if !ok {
		t.Fatal("expected", true, "got", false)
	}
	if v != "bar" {
		t.Fatal("expected", "bar", "got", v)
	}
	v, ok = other.(*context).Storage["foo"]
	if !ok {
		t.Fatal("expected", true, "got", false)
	}
	if v != "bar" {
		t.Fatal("expected", "bar", "got", v)
	}
}

func Test_Cancel(t *testing.T) {
	ctx := testNewContext(t)

	select {
	case <-time.After(5 * time.Millisecond):
		// Here a timeout should happen because the context was not yet canceled.
	case <-ctx.Done():
		t.Fatal("expected", "timeout", "got", "cancel")
	}

	// Canceling the context should cause the created context to be canceled.
	ctx.Cancel()

	select {
	case <-time.After(5 * time.Millisecond):
		t.Fatal("expected", "cancel", "got", "timeout")
	case <-ctx.Done():
		// The test was successful when the created context was canceled on
		// cancelation.
	}
}

func Test_Cancel_Underlying(t *testing.T) {
	nativeCtx, cancelFunc := nativecontext.WithCancel(nativecontext.Background())

	config := DefaultConfig()
	config.Context = nativeCtx
	ctx, err := New(config)
	if err != nil {
		t.Fatal("expected", nil, "got", err)
	}

	select {
	case <-time.After(5 * time.Millisecond):
		// Here a timeout should happen because the underlying context was not yet
		// canceled.
	case <-ctx.Done():
		t.Fatal("expected", "timeout", "got", "cancel")
	}

	// Canceling the native context should cause the created context to be
	// canceled, because the native context is configured as underlying context.
	cancelFunc()

	select {
	case <-time.After(5 * time.Millisecond):
		t.Fatal("expected", "cancel", "got", "timeout")
	case <-ctx.Done():
		// The test was successful when the created context was canceled on
		// cancelation of the configured native context.
	}
}

func testNewContext(t *testing.T) Context {
	var ctx Context
	{
		var err error
		ctx, err = New(DefaultConfig())
		if err != nil {
			t.Fatal("expected", nil, "got", err)
		}
		ctx.(*context).Storage["foo"] = "bar"
		ctx.(*context).Storage["other"] = 45
	}

	return ctx
}
