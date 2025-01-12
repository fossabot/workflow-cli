package parser

import (
	"bytes"
	"errors"
	"testing"

	"github.com/drycc/workflow-cli/pkg/testutil"
	"github.com/stretchr/testify/assert"
)

// Create fake implementations of each method that return the argument
// we expect to have called the function (as an error to satisfy the interface).

func (d FakeDryccCmd) TLSInfo(string) error {
	return errors.New("tls:info")
}

func (d FakeDryccCmd) TLSForceEnable(string) error {
	return errors.New("tls:force:enable")
}

func (d FakeDryccCmd) TLSForceDisable(string) error {
	return errors.New("tls:force:disable")
}

func (d FakeDryccCmd) TLSAutoEnable(string) error {
	return errors.New("tls:auto:disable")
}

func (d FakeDryccCmd) TLSAutoDisable(string) error {
	return errors.New("tls:auto:disable")
}

func TestTLS(t *testing.T) {
	t.Parallel()

	cf, server, err := testutil.NewTestServerAndClient()
	if err != nil {
		t.Fatal(err)
	}
	defer server.Close()
	var b bytes.Buffer
	cmdr := FakeDryccCmd{WOut: &b, ConfigFile: cf}

	// cases defines the arguments and expected return of the call.
	// if expected is "", it defaults to args[0].
	cases := []struct {
		args     []string
		expected string
	}{
		{
			args:     []string{"tls:info"},
			expected: "",
		},
		{
			args:     []string{"tls:force:enable"},
			expected: "",
		},
		{
			args:     []string{"tls:force:disable"},
			expected: "",
		},
		{
			args:     []string{"tls"},
			expected: "tls:info",
		},
	}

	// For each case, check that calling the route with the arguments
	// returns the expected error, which is args[0] if not provided.
	for _, c := range cases {
		var expected string
		if c.expected == "" {
			expected = c.args[0]
		} else {
			expected = c.expected
		}
		err = TLS(c.args, cmdr)
		assert.Error(t, errors.New(expected), err)
	}
}
