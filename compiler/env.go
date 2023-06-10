package compiler

import (
	"errors"
	"fmt"
)

type Tipe struct {
	Name string
	Size int
}

type Binding struct {
	Name string
	Tipe Tipe
}

type Env struct {
	globals []Binding
	tipes   map[string]Tipe
}

func NewEnv() *Env {
	return &Env{
		globals: []Binding{},
		tipes: map[string]Tipe{
			"int":  T_INT,
			"bool": T_BOOL,
		},
	}
}

func (env *Env) addBinding(s string, tipe string) (*Env, error) {
	if s == NEVER {
		panic(fmt.Sprint(s, " is not a valid name for assignment. If you are seeing this error, something has gone terribly wrong."))
	}

	t_var, ok := env.tipes[tipe]
	if !ok {
		err := fmt.Sprint("unknown type ", tipe)
		return env, errors.New(err)
	}

	return &Env{
		globals: append(env.globals, Binding{Name: s, Tipe: t_var}),
		tipes:   env.tipes,
	}, nil
}

/*
Used when an element has been pushed onto the stack without calling 'declare',
and lexical address resolution still needs to work. The created symbol is gauranteed
to never match a local bound by the program author.
*/
func (env *Env) addNever(size int) *Env {
	return &Env{
		globals: append(env.globals, Binding{Name: NEVER, Tipe: T_NEVER(size)}),
	}
}

func (env *Env) lexicalAddress(s string) (int, Tipe, error) {
	jump := 0
	for i := len(env.globals) - 1; i >= 0; i-- {
		if s == env.globals[i].Name {
			return jump, env.globals[i].Tipe, nil
		}
		jump += env.globals[i].Tipe.Size
	}
	return 0, T_NEVER(0), errors.New(fmt.Sprint("unbound variable ", s))
}

func (env *Env) lookupTipe(tipe string) (Tipe, bool) {
	res, ok := env.tipes[tipe]
	return res, ok
}

const (
	NEVER = "-"
)

var T_INT = Tipe{
	Name: "int",
	Size: 8,
}

func T_NEVER(size int) Tipe {
	return Tipe{
		Name: "never",
		Size: size,
	}
}

var T_BOOL = Tipe{
	Name: "bool",
	Size: 8,
}
