package inspection

type Inspection struct {
	Type  Type
	Value any
}

type Type int

// inspection types
const (
	DefaultType Type = iota + 1
	BoolType
)

const (
	Dead = iota + 1
	Less50
	More50
)

type Default int

type Bool bool
