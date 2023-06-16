package compiler

import (
	"errors"
	"fmt"
	"monkey/parser"
)

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

func (env *Env) addBinding(s string, tipe parser.TypeExpression) (*Env, error) {
	if s == NEVER {
		panic(fmt.Sprint(s, " is not a valid name for assignment. If you are seeing this error, something has gone terribly wrong."))
	}

	t_var, ok := env.lookupTipe(tipe)
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
Used when an element has been pushed onto the stack without calling 'addBinding',
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

func (env *Env) lookupTipe(tipe parser.TypeExpression) (Tipe, bool) {
	switch tipe := tipe.(type) {
	case *parser.LiteralType:
		name := tipe.Name
		res, ok := env.tipes[name]
		return res, ok
	case *parser.ArrowType:
		// If it is a function, we need to check all parameters
		params := tipe.Parameters
		for _, param := range params {
			if _, ok := env.lookupTipe(param); !ok {
				return T_NEVER(0), false
			}
		}
		// We also need to check the return type
		if _, ok := env.lookupTipe(tipe.Returns); !ok {
			return T_NEVER(0), false
		}
		return T_ARROW, true
	default:
		return T_NEVER(0), false
	}
}

const (
	NEVER = "-"
)
