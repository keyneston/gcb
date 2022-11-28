package gcb

import (
	"encoding/json"
	"log"
	"os"
	"testing"
	"time"

	"github.com/keyneston/gcb/printable"
	"github.com/stretchr/testify/assert"
)

func TestTreeGet(t *testing.T) {
	t.Skip()

	r := &Tree{}

	timeout := printable.Duration(time.Second)
	s := Settings{Timeout: &timeout}

	r.set("foo", s)
	r.set("foo!bar!baz", s)
	r.set("foo!bar", s)

	t.Errorf("Tree: %#v", r)

	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "    ")

	enc.Encode(r)

	n := r.Get("foo!bar!baz")
	log.Printf("Got n: %#v", n)
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
