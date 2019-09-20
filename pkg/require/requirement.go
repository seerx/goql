package require

import (
	"fmt"
	"reflect"
)

type Requirement struct {
	missMap map[string]int
}

func New() *Requirement {
	return &Requirement{
		missMap: map[string]int{},
	}
}

var typeOfRequirement = reflect.TypeOf(Requirement{})

func IsRequirement(typ reflect.Type) bool {
	if typ.Kind() == reflect.Ptr {
		return typ.Elem() == typeOfRequirement
	}
	return typ == typeOfRequirement
}

func (r *Requirement) Add(param string) {
	r.missMap[param] = 0
}

func (r *Requirement) Requires(params ...string) {
	for _, p := range params {
		if _, ok := r.missMap[p]; ok {
			panic(fmt.Errorf("Require param %s", p))
		}
	}
}
