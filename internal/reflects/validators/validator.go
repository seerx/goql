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

const (
	fmtMax  = "%s's %s must little then %s, but got "
	fmtMaxE = "%s's %s must little then or equal %s, but got "

	fmtMin  = "%s's %s must great then %s, but got "
	fmtMinE = "%s's %s must great then or equal %s, but got "

	fmtBetween   = "%s's %s must great then %s and little then %s, but got "
	fmtBetweenEE = "%s's %s must great then or equal %s and little then or equal %s, but got "
	fmtBetween_E = "%s's %s must great then %s and little then or equal %s, but got "
	fmtBetweenE_ = "%s's %s must great then or equal %s and little then %s, but got "
)

func getFmt(field string, desc string,
	limitMax bool,
	max string,
	includeMax bool,
	limitMin bool,
	min string,
	includeMin bool, got string) string {

	errfmt := field + "'s " + desc

	minCond := ""
	if limitMin {
		if includeMin {
			minCond = "great then or equal " + min
		} else {
			minCond = "great then " + min
		}
	}

	maxCond := ""
	if limitMax {
		if includeMax {
			maxCond = "little then or equal " + max
		} else {
			maxCond = "little then " + max
		}
	}

	if minCond != "" {
		if maxCond != "" {
			errfmt = fmt.Sprintf("%s must %s and %s", errfmt, minCond, maxCond)
		} else {
			errfmt = fmt.Sprintf("%s must %s", errfmt, minCond)
		}
	} else {
		errfmt = fmt.Sprintf("%s must %s", errfmt, maxCond)
	}

	return errfmt + ", but got " + got
}
