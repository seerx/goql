package validators

import (
	"fmt"
)

//内置类型检查

// ValueChecker 数据检查
type Validator interface {
	Check(val interface{}) error
}

func ignoreCh(ch string) bool {
	return " " == ch
}

func typeError(expected string) error {
	return fmt.Errorf("expect type %s", expected)
}
