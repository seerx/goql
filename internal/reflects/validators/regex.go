package validators

import (
	"fmt"
	"regexp"
)

// RegexValidator 正则表达式验证
// 对应 tag 中的 regexp
type RegexpValidator struct {
	field  string
	regstr string
	regex  *regexp.Regexp
}

func (v *RegexpValidator) Check(val interface{}) error {
	str, ok := val.(string)
	if !ok {
		return typeError("string")
	}
	if !v.regex.Match([]byte(str)) {
		return fmt.Errorf("%s's value do not match with regular expression %s", v.field, v.regstr)
	}
	return nil
}

func CreateRegexpValidator(fieldName string, exp string) *RegexpValidator {
	reg, err := regexp.Compile(exp)
	if err != nil {
		panic(fmt.Errorf("Invalid regular expression %s = %s", fieldName, exp))
		return nil
	}
	return &RegexpValidator{
		regstr: exp,
		field:  fieldName,
		regex:  reg,
	}
}
