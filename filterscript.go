package jsonpath

import (
	"github.com/sirupsen/logrus"
	"reflect"
	"regexp"
)

//filter expression
//?($.name == "xue")

type filterScript struct {
	root *JsonPath
	t    *pathtoken
}

// TODO action set and unset
func (f *filterScript) eval(action string, cv interface{}, optionalValue interface{}) ([]interface{}, bool) {
	varr := f.getCandidates(cv)
	expr := f.t.v.(expression)
	// do filter
	op := expr.items[1].val
	result := []interface{}{}
	for _, v := range varr {
		sjp := &JsonPath{Data: v}
		val, got, _ := sjp.Find(expr.items[0].val).First()

		if got && reflect.TypeOf(val) == reflect.TypeOf("string") {
			switch expr.items[2].typ {
			case "jsonpath":
				if expr.items[2].val[0] == '$' {
					// root jsonpath
				}
				if expr.items[2].val[0] == '@' {
					njp := &JsonPath{Data: v}
					if rval, got, _ := njp.Find(expr.items[2].val).First(); got == true && reflect.TypeOf(rval) == reflect.TypeOf("string") {
						if f.doOp(val.(string), rval.(string), op) {
							result = append(result, v)
						}
					}
				}
			default:
				if f.doOp(val.(string), expr.items[2].val, op) {
					result = append(result, v)
				}
			}
		}
	}
	return result, true
}

func (f *filterScript) getCandidates(cv interface{}) []interface{} {
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
		return nil
	}
	return varr
}

func (f *filterScript) doOp(lv string, rv string, op string) bool {

	switch op {
	case "==":
		fallthrough
	case "!=":
		fallthrough
	case ">":
		fallthrough
	case ">=":
		fallthrough
	case "<":
		fallthrough
	case "<=":
		if ret, err := compare_valstring(lv, rv, op); err != nil {

			// handle err
		} else {
			return ret
		}

	case "=~":
		if r, err := regexp.Compile(rv); err != nil {
			// handle err
		} else {
			return r.MatchString(lv)
		}
	}
	// error op not supported
	logrus.Error("error: filter expression op{" + op + "} not supported")
	return false
}
