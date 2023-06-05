package compiler

import (
	"errors"
	"fmt"
)

type Env struct {
	locals []string
}

func (env *Env) addLocal(s string) *Env {
	return &Env{
		locals: append(env.locals, s),
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
