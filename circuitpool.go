package gcb

import "sync"

type circuitPool struct {
	circuits map[string]*Circuit
	rw       *sync.RWMutex
}

func (p *circuitPool) get(name string) *Circuit {
	// Attempt a read-only lock path.
	if circuit := p.attemptGet(name); circuit != nil {
		return circuit
	}

	// if that fails, then go through the more expensive write lock path
	return p.insert(name)
}

// attemptGet attempts to get the circuit in question, but will return nil if
// the circuit doesn't exist. This is in a sub-function for locking reasons.
func (p *circuitPool) attemptGet(name string) *Circuit {
	p.rw.RLock()
	defer p.rw.RUnlock()

	return p.circuits[name]
}

// insert acquires a write lock and then attempts an additional read from the
// map. If the circuit still doesn't exist, then it creates one. If it was
// created in between the call attemptGet, and insert, then it will simply
// return what as acquired.
func (p *circuitPool) insert(name string) *Circuit {
	p.rw.Lock()
	defer p.rw.Unlock()

	if circuit, ok := p.circuits[name]; ok {
		return circuit
	}

	circuit := NewCircuit(name)
	p.circuits[name] = circuit
	return circuit
}

// circuits is the shared pool of circuits. This is a global variable
// as a consequence of how go generics work, and the inability to use generics
// on methods.

var circuits = circuitPool{
	circuits: map[string]*Circuit{},
	rw:       &sync.RWMutex{},
}
