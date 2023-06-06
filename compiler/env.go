package compiler

import (
	"errors"
	"fmt"
)

type Env struct {
	locals []string
}

func (env *Env) addLocal(s string) *Env {
	if s == NEVER_LOCAL {
		panic(fmt.Sprint(s, " is not a valid name for assignment. If you are seeing this error, something has gone terribly wrong."))
	}
	return &Env{
		locals: append(env.locals, s),
	}
}

/*
Used when an element has been pushed onto the stack without calling 'addLocal',
and lexical address resolution still needs to work. The created symbol is gauranteed
to never match a local bound by the program author.
*/
func (env *Env) addNever() *Env {
	return &Env{
		locals: append(env.locals, NEVER_LOCAL),
	}
}

func (env *Env) lexicalAddress(s string) (int, error) {
	for i := len(env.locals) - 1; i >= 0; i-- {
		if s == env.locals[i] {
			return len(env.locals) - i - 1, nil
		}
	}
	return 0, errors.New(fmt.Sprint("unbound variable ", s))
}

const (
	NEVER_LOCAL = "-"
)
