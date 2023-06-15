package sql

import "fmt"

type Placeholder struct {
	idx  int
	Args []interface{}
}

func BuildPlaceholder(cap int) Placeholder {
	return Placeholder{
		idx:  0,
		Args: make([]interface{}, 0, cap),
	}
}

func (pl *Placeholder) Push(arg interface{}) string {
	pl.Args = append(pl.Args, arg)
	pl.idx += 1
	return fmt.Sprintf("$%d", pl.idx)
}
