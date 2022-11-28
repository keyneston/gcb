package printable

import (
	"encoding/json"
	"sync/atomic"
	"time"
)

type Atomic struct {
	atomic.Value
}

func (p *Atomic) MarshalJSON() ([]byte, error) {
	return json.Marshal(p.Load())
}

type Duration time.Duration

func (d Duration) MarshalJSON() ([]byte, error) {
	return json.Marshal(time.Duration(d).String())
}

func (d Duration) D() time.Duration {
	return time.Duration(d)
}
