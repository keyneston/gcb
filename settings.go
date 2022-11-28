package gcb

import (
	"strings"
	"sync/atomic"
	"time"

	"github.com/keyneston/gcb/printable"
)

var defaultSettings *atomic.Value = &atomic.Value{}

// Separator is used to separate sub-keys for hierarchical settings.
var Separator = "!"

func init() {
	timeout := printable.Duration(time.Second)
	SetDefaultSettings(Settings{
		Timeout:         &timeout,
		FallbackTimeout: &timeout,
	})
}

func GetDefaultSettings() Settings {
	return *defaultSettings.Load().(*Settings)
}

func SetDefaultSettings(s Settings) {
	defaultSettings.Store(&s)
}

type Tree struct {
	*node
}

func (r *Tree) Get(key string) *node {
	n := r.node
	for n != nil {
		switch strings.Compare(key, n.Name) {
		case 0:
			return n
		case -1:
			n = n.Left
		case 1:
			n = n.Right
		}
	}

	return nil
}

func (r *Tree) set(name string, settings Settings) error {
	segments := strings.Split(name, Separator)

	// TODO: what is the correct behaviour if a name is already set.
	// TODO: some type of rebalancing
	for i := range segments {
		name := strings.Join(segments[0:i], Separator)
		if r.node == nil {
			r.node = &node{Name: name}
		}

		n := &r.node
		for *n != nil {
			if strings.Compare(name, (*n).Name) < 0 {
				n = &(*n).Left
			} else {
				n = &(*n).Right
			}
		}

		if *n == nil {
			*n = &node{Name: name}
			// TODO: settings should be nil????
			(*n).setSettings(settings)
		}
	}

	return nil
}

// splitName takes a name, splits it into hierarchical segments.
func splitName(name string) []string {
	segments := strings.Split(name, Separator)

	out := make([]string, 0, len(segments))
	for n := range segments {
		out = append(out, strings.Join(segments[0:n+1], Separator))
	}

	return out
}

type node struct {
	Name  string `json:"name,omitempty"`
	Left  *node  `json:"left,omitempty"`
	Right *node  `json:"right,omitempty"`

	Settings *printable.Atomic `json:"settings"`
	circuit  *printable.Atomic
}

func (n *node) GetSettings() Settings {
	return *n.Settings.Load().(*Settings)
}

func (n *node) setSettings(settings Settings) *node {
	if n.Settings == nil {
		// TODO: race condition here
		n.Settings = &printable.Atomic{}
	}
	n.Settings.Store(&settings)
	return n
}

func (n *node) Circuit() *Circuit {
	return n.circuit.Load().(*Circuit)
}

// Everything is a nil-pointer to allow optional settings.
type Settings struct {
	Timeout         *printable.Duration `json:"timeout"`
	FallbackTimeout *printable.Duration `json:"fallback_timeout"`

	Percent *float32 `json:"percent"`
}
