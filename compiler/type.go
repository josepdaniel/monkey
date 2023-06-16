package compiler

type Tipe struct {
	Name string
	Size int
}

func (t Tipe) IsEqualTo(other Tipe) bool {
	return t.Name == other.Name && t.Size == other.Size
}

func T_NEVER(size int) Tipe {
	return Tipe{
		Name: "never",
		Size: size,
	}
}

var T_INT = Tipe{
	Name: "int",
	Size: 8,
}

var T_BOOL = Tipe{
	Name: "bool",
	Size: 8,
}

// An arrow type represents the address of a function.
var T_ARROW = Tipe{
	Name: "arrow",
	Size: 8,
}
