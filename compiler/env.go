package compiler

type Env struct {
	locals []string
}

func (env *Env) addLocal(s string) *Env {
	return &Env{
		locals: append(env.locals, s),
	}
}
