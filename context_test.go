package context

import (
	"bytes"
	"encoding/json"
	"testing"
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
