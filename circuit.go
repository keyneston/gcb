package gcb

func GetCircuit(name string) *Circuit {
	return circuits.get(name)
}

type Circuit struct {
	name string
}

func NewCircuit(name string) *Circuit {
	return &Circuit{
		name: name,
	}
}
