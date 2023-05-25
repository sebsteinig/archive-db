package utils

import "fmt"

type Placholder struct {
	idx  int
	Args []interface{}
}

func (pl *Placholder) Build(length int, cap int) {
	*pl = Placholder{
		idx:  0,
		Args: make([]interface{}, length, cap),
	}
}

func (pl *Placholder) Get(arg interface{}) string {
	pl.Args = append(pl.Args, arg)
	pl.idx += 1
	return fmt.Sprintf("$%d", pl.idx)
}

func (pl *Placholder) Wrap(arg interface{}) string {
	pl.Args = append(pl.Args, fmt.Sprintf("'%v'", arg))
	pl.idx += 1
	return fmt.Sprintf("$%d", pl.idx)
}
