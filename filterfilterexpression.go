package jsonpath

import (
	"strconv"
	"strings"
	"testing"
)

//filter expression
//?($.name == "xue")

type filterFilterExpression struct {
	t *pathtoken
}

// TODO action set and unset
func (f *filterFilterExpression) eval(action string, cv interface{}, optionalValue interface{}) ([]interface{}, bool) {

	expr := f.t.v.(expression)
	// do filter
	op := expr.items[1].val
	lv := expr.items[0]
	rv := expr.items[2]
	result := []interface{}{}
	varr := []interface{}{}
	switch cv.(type) {
	case map[string]interface{}:
		for _, v := range cv.(map[string]interface{}) {
			varr = append(varr, v)
		}
	case []interface{}:
		for _, v := range cv.([]interface{}) {
			varr = append(varr, v)
		}
	default:
		return nil, false
	}
	switch op {
	case "==":
		fallthrough
	case "!=":
		fallthrough
	case ">":
		fallthrough
	case ">=":
		fallthrough
	case "<=":
		fallthrough
	case "<":
		for _, v := range varr {
			sjp := &JsonPath{Data: cv}
			val, got, _ := sjp.Find(expr.items[0].val).First()
			if got {
				switch rv.typ {
				case "string":
					if val == rv.val {
						result = append(result, v)
					}
				case "num":
					if !isnumber(val) {
						continue
					}
					if valnum, err := strconv.Atoi(rv.val); err != nil {
						continue
					} else if val == valnum {
						result = append(result, v)
					}
				case "jsonpath":
					if rv.val[0] == '$' {
						// root jsonpath
					}
					if rv.val[0] == '@' {
						njp := &JsonPath{Data: v}
						if rval, got, _ := njp.Find(rv.val).First(); got == true && val == rval {
							result = append(result, v)
						}
					}
				}
			}
		}

	}
}
