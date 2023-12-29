package main

import (
	"bytes"
	"testing"
)

func TestHandleOptions(t *testing.T) {

	tests := []struct {
		desc     string
		showHelp bool
		region   string
		profile  string
		want     AppOptions
	}{
		{
			showHelp: true,
			want:     AppOptions{ShowHelp: true},
		},
	}

	for _, test := range tests {

		buffer := bytes.Buffer{}
		got := HandleOptions(&buffer, true, "", "")

		if test.showHelp && buffer.String() != helpText {
			t.Errorf("%s: help text want: %q but got %q", test.desc, helpText, buffer.String())
		}

		if test.want != got {
			t.Errorf("%s: wanted options %+v but got %+v", test.desc, test.want, got)
		}
	}
}
