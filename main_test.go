package main

import (
	"bytes"
	"testing"
)

func TestHandleOptions(t *testing.T) {
	buffer := bytes.Buffer{}

	HandleOptions(&buffer, true, "", "")

	got := buffer.String()
	want := helpText

	if got != want {
		t.Errorf("got %q want %q", got, want)
	}
}
