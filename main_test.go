package main

import (
	"bytes"
	"testing"
)

func TestHandleOptions(t *testing.T) {

	tests := []struct {
		desc    string
		help    bool
		region  string
		profile string
		want    AppOptions
	}{
		{
			desc: "help flag supplied",
			help: true,
			want: AppOptions{Help: true},
		},
		{
			desc:   "region supplied without profile",
			region: "us-west-2",
			want:   AppOptions{Help: true},
		},
		{
			desc:    "profile supplied",
			profile: "devops-profile",
			want:    AppOptions{Profile: "devops-profile"},
		},
	}

	for _, test := range tests {

		buffer := bytes.Buffer{}
		got := HandleOptions(&buffer, test.help, test.region, test.profile)

		if test.help && buffer.String() != helpText {
			t.Errorf("%s: help text want %q but got %q", test.desc, helpText, buffer.String())
		}

		if test.profile == "" && !test.help && buffer.String() != missingProfileWarning+helpText {
			t.Errorf("%s: want %q but got %q", test.desc, missingProfileWarning+helpText, buffer.String())
		}

		if test.want != got {
			t.Errorf("%s: wanted options %+v but got %+v", test.desc, test.want, got)
		}
	}
}
