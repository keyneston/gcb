package gcb

import (
	"testing"
	"time"

	"github.com/keyneston/gcb/printable"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTreeGet(t *testing.T) {
	r := &Tree{}

	r.set("foo", Settings{Timeout: printable.Dur(time.Second)})
	r.set("foo!bar!baz", Settings{FallbackTimeout: printable.Dur(5 * time.Second)})
	r.set("foo!bar", Settings{Timeout: printable.Dur(time.Second * 2)})

	needle := "foo!bar!baz"
	s := r.Get(needle)
	require.NotNil(t, s, "settings must not be nil")
	require.NotNil(t, s.Timeout, "Timeout must not be nil")
	require.NotNil(t, s.FallbackTimeout, "FallbackTimeout must not be nil")
	require.Equalf(t, time.Second*2, time.Duration(*s.Timeout),
		"s.Timeout = %v; want %v", s.Timeout, time.Second*3)
	require.Equalf(t, time.Second*5, time.Duration(*s.FallbackTimeout),
		"s.FallbackTimeout = %v; want %v", s.FallbackTimeout, time.Second*5)
}

func TestSplitName(t *testing.T) {
	type testCase struct {
		input    string
		expected []string
	}

	for _, c := range []testCase{
		{"foo", []string{"foo"}},
		{"foo!bar!baz", []string{"foo", "foo!bar", "foo!bar!baz"}},
	} {
		out := splitName(c.input)

		assert.Equal(t, c.expected, out, "Case: %#q", c.input)
	}
}
