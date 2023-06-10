package test

import (
	"fmt"
	"reflect"
	"regexp"
)

func ShouldMatchWithRegex(actual interface{}, expected ...interface{}) string {
	value, valueIsString := actual.(string)

	if !valueIsString {
		return fmt.Sprintf("Actual argument to this assertion must be string(you provided %v: %v)", reflect.TypeOf(actual), actual)
	}

	if len(expected) == 0 {
		return "Expected value to this assertion must be specified"
	}

	for _, e := range expected {
		regex, regexIsString := e.(string)
		if !regexIsString {
			return fmt.Sprintf("Expected argument to this assertion must be string(you provided %v: %v)", reflect.TypeOf(e), e)
		}

		matched, err := regexp.MatchString(regex, value)
		if err != nil || !matched {
			return fmt.Sprintf("Expected      '%v'\nregex match with '%v'\n(but it didn't '%v')!",
				regex, value, err)
		}
	}
	return ""
}
