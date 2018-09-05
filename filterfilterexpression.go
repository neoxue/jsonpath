package jsonpath

import (
	"strconv"
	"strings"
)

//filter expression
//?($.name == "xue")

type filterFilterExpression struct {
	t *pathtoken
}

// TODO action set and unset
func (f *filterFilterExpression) eval(action string, cv interface{}, optionalValue interface{}) ([]interface{}, bool) {

}
