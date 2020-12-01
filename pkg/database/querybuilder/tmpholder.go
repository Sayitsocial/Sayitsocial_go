package querybuilder

import "reflect"

type tmpHolder struct {
	name   string
	typeOf string
	value  reflect.Value
}
