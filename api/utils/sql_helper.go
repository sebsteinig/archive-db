package utils

import "fmt"

type Placeholder struct {
	idx  int
	Args []interface{}
}

func (pl *Placeholder) Build(length int, cap int) {
	*pl = Placeholder{
		idx:  0,
		Args: make([]interface{}, length, cap),
	}
}

func (pl *Placeholder) Get(arg interface{}) string {
	pl.Args = append(pl.Args, arg)
	pl.idx += 1
	return fmt.Sprintf("$%d", pl.idx)
}
