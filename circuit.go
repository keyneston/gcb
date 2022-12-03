package gcb

func GetCircuit(name string) *Circuit {
	return circuits.get(name)
}

type Circuit struct {
	name    string
	percent int32
}

func NewCircuit(name string) *Circuit {
	return &Circuit{
		name: name,
	}
}
