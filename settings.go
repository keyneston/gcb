package gcb

import (
	"encoding/json"
	"log"
	"sort"
	"strings"
	"sync/atomic"
	"time"

	"github.com/keyneston/gcb/printable"
)

var defaultSettings *atomic.Value = &atomic.Value{}

// Separator is used to separate sub-keys for hierarchical settings.
var Separator = "!"

func init() {
	timeout := printable.Dur(time.Second)
	SetDefaultSettings(Settings{
		Timeout:         timeout,
		FallbackTimeout: timeout,
	})
}

func mustJSON(v any) string {
	out, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		log.Fatalf("JSON shouldn't fail to encode: %v", err)
	}

	return string(out)
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

func (r *Tree) String() string {
	return mustJSON(r)
}

func (r *Tree) Get(key string) Settings {
	split := splitName(key)
	results := make([]namedSetting, 0, len(split))

	n := r.node
loop:
	for n != nil {
		for _, needle := range split {
			if n.Name == needle {
				results = append(results, namedSetting{n.Name, n.GetSettings()})
			}
		}

		switch strings.Compare(key, n.Name) {
		case 0:
			break loop
		case -1:
			n = n.Left
		case 1:
			n = n.Right
		}
	}

	sort.Sort(namedSettings(results))

	final := &Settings{}
	for _, r := range results {
		final.Merge(&r.Settings)
	}

	return *final
}

func (r *Tree) set(name string, settings Settings) error {
	segments := splitName(name)

	// TODO: what is the correct behaviour if a name is already set.
	// TODO: some type of rebalancing
	for _, name := range segments {
		//name := strings.Join(segments[0:i], Separator)
		if r.node == nil {
			r.node = &node{Name: name}
		}

		n := &r.node
	tree_descent:
		for *n != nil {
			switch strings.Compare(name, (*n).Name) {
			case 0:
				(*n).setSettings(settings)
				break tree_descent
			case -1:
				n = &(*n).Left
			case 1:
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
	PrimaryTimeout  *printable.Duration `json:"primary_timeout"`
	FallbackTimeout *printable.Duration `json:"fallback_timeout"`

	BreakPercent *float32 `json:"break_percent"`
}

func (s Settings) String() string {
	return mustJSON(s)
}

func (s *Settings) Merge(s2 *Settings) {
	if s2 == nil {
		return
	}

	if s2.Timeout != nil {
		s.Timeout = s2.Timeout
	}
	if s2.PrimaryTimeout != nil {
		s.PrimaryTimeout = s2.PrimaryTimeout
	}
	if s2.FallbackTimeout != nil {
		s.FallbackTimeout = s2.FallbackTimeout
	}
	if s2.BreakPercent != nil {
		s.BreakPercent = s2.BreakPercent
	}
}

type namedSetting struct {
	Name     string
	Settings Settings
}

func (n namedSetting) String() string {
	return mustJSON(n)
}

type namedSettings []namedSetting

func (s namedSettings) Len() int {
	return len(s)
}
func (s namedSettings) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}
func (s namedSettings) Less(i, j int) bool {
	return len(s[i].Name) < len(s[j].Name)
}
