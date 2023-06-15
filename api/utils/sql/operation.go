package sql

import (
	"fmt"
	"reflect"
	"strings"
)

type SqlBuilder interface {
	Build(pl *Placeholder) string
}
type EqualBuilder struct {
	Key          string
	Value        any
	And_Prefix   bool
	Or_Prefix    bool
	Where_Prefix bool
}

func (builder EqualBuilder) Build(pl *Placeholder) string {
	if reflect.ValueOf(builder.Value).IsZero() {
		return ""
	}
	prefix := ""
	if builder.Where_Prefix {
		prefix += " WHERE "
	} else if builder.And_Prefix {
		prefix += " AND "
	} else if builder.Or_Prefix {
		prefix += " OR "
	}
	return prefix + fmt.Sprintf("%s = %s", builder.Key, pl.Push(builder.Value))
}

type LikeBuilder struct {
	Key          string
	Value        any
	And_Prefix   bool
	Or_Prefix    bool
	Where_Prefix bool
}

func (builder LikeBuilder) Build(pl *Placeholder) string {
	if reflect.ValueOf(builder.Value).IsZero() {
		return ""
	}
	prefix := ""
	if builder.Where_Prefix {
		prefix += " WHERE "
	} else if builder.And_Prefix {
		prefix += " AND "
	} else if builder.Or_Prefix {
		prefix += " OR "
	}
	return prefix + fmt.Sprintf("%s LIKE %s || '%%'", builder.Key, pl.Push(builder.Value))
}

type ILikeBuilder struct {
	Key          string
	Value        any
	And_Prefix   bool
	Or_Prefix    bool
	Where_Prefix bool
}

func (builder ILikeBuilder) Build(pl *Placeholder) string {
	if reflect.ValueOf(builder.Value).IsZero() {
		return ""
	}
	prefix := ""
	if builder.Where_Prefix {
		prefix += " WHERE "
	} else if builder.And_Prefix {
		prefix += " AND "
	} else if builder.Or_Prefix {
		prefix += " OR "
	}
	return prefix + fmt.Sprintf("%s ILIKE %s || '%%'", builder.Key, pl.Push(builder.Value))
}

type FLikeBuilder struct {
	Key          string
	Value        any
	And_Prefix   bool
	Or_Prefix    bool
	Where_Prefix bool
}

func (builder FLikeBuilder) Build(pl *Placeholder) string {
	if reflect.ValueOf(builder.Value).IsZero() {
		return ""
	}
	prefix := ""
	if builder.Where_Prefix {
		prefix += " WHERE "
	} else if builder.And_Prefix {
		prefix += " AND "
	} else if builder.Or_Prefix {
		prefix += " OR "
	}
	return prefix + fmt.Sprintf("%s ILIKE '%%' || %s || '%%'", builder.Key, pl.Push(builder.Value))
}

type InBuilder struct {
	Key          string
	Value        []any
	And_Prefix   bool
	Or_Prefix    bool
	Where_Prefix bool
}

func (builder InBuilder) Build(pl *Placeholder) string {
	if reflect.ValueOf(builder.Value).IsZero() {
		return ""
	}
	array := make([]string, len(builder.Value))
	for i, value := range builder.Value {
		array[i] = pl.Push(value)
	}
	prefix := ""
	if builder.Where_Prefix {
		prefix += " WHERE "
	} else if builder.And_Prefix {
		prefix += " AND "
	} else if builder.Or_Prefix {
		prefix += " OR "
	}
	return prefix + fmt.Sprintf("%s IN (%s)", builder.Key, strings.Join(array, ","))
}

type AndBuilder struct {
	Value        []SqlBuilder
	And_Prefix   bool
	Or_Prefix    bool
	Where_Prefix bool
}

func (builder *AndBuilder) And(sqlb SqlBuilder) {
	builder.Value = append(builder.Value, sqlb)
}

func (builder *AndBuilder) AndAll(builders ...SqlBuilder) {
	for _, b := range builders {
		builder.Value = append(builder.Value, b)
	}
}

func (builder AndBuilder) Build(pl *Placeholder) string {
	if len(builder.Value) == 0 {
		return ""
	}
	array := make([]string, 0, len(builder.Value))
	for _, value := range builder.Value {
		if sql := value.Build(pl); sql != "" {
			array = append(array, sql)
		}
	}
	prefix := ""
	if builder.Where_Prefix {
		prefix += " WHERE "
	} else if builder.And_Prefix {
		prefix += " AND "
	} else if builder.Or_Prefix {
		prefix += " OR "
	}
	return prefix + fmt.Sprintf("(%s)", strings.Join(array, " AND "))
}

type OrBuilder struct {
	Value        []SqlBuilder
	And_Prefix   bool
	Or_Prefix    bool
	Where_Prefix bool
}

func (builder *OrBuilder) Or(sqlb SqlBuilder) {
	builder.Value = append(builder.Value, sqlb)
}
func (builder *OrBuilder) OrAll(builders ...SqlBuilder) {
	for _, b := range builders {
		builder.Value = append(builder.Value, b)
	}
}

func (builder OrBuilder) Build(pl *Placeholder) string {
	if len(builder.Value) == 0 {
		return ""
	}
	array := make([]string, 0, len(builder.Value))
	for _, value := range builder.Value {
		if sql := value.Build(pl); sql != "" {
			array = append(array, sql)
		}
	}
	prefix := ""
	if builder.Where_Prefix {
		prefix += " WHERE "
	} else if builder.And_Prefix {
		prefix += " AND "
	} else if builder.Or_Prefix {
		prefix += " OR "
	}
	return prefix + fmt.Sprintf("(%s)", strings.Join(array, " OR "))
}
