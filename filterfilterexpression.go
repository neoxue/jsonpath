package jsonpath

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"reflect"
	"regexp"
	"strconv"
)

//filter expression
//?($.name == "xue")

type filterFilterExpression struct {
	root *JsonPath
	t    *pathtoken
}

// TODO action set and unset
func (f *filterFilterExpression) eval(action string, cv interface{}, optionalValue interface{}) ([]interface{}, bool) {
	varr := f.getCandidates(cv)
	expr := f.t.v.(*expression)
	// do filter
	op := expr.items[1].val
	result := []interface{}{}
	for _, v := range varr {
		sjp := &JsonPath{Data: v}
		val, got, _ := sjp.Find(expr.items[0].val).First()

		if got && reflect.TypeOf(val) == reflect.TypeOf("string") {
			switch expr.items[2].typ {
			case "jsonpath":
				var njp *JsonPath
				if expr.items[2].val[0] == '$' {
					njp = f.root
				}
				if expr.items[2].val[0] == '@' {
					njp = &JsonPath{Data: v}
				}
				if rval, got, err := njp.Find(expr.items[2].val).First(); err != nil {
					logrus.Error(err)
				} else if got == true {
					if f.doOp(val, rval, op) {
						result = append(result, v)
					}
				} else {
					logrus.Error("rv jsonpath {" + expr.items[2].val + "} nothing found")
				}
			default:
				if f.doOp(val, expr.items[2].val, op) {
					result = append(result, v)
				}
			}
		}
	}
	return result, true
}

func (f *filterFilterExpression) getCandidates(cv interface{}) []interface{} {
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

/*
	number  -> number compare
	string  -> string compare
	if =~ all trans to string
*/
func (f *filterFilterExpression) doOp(lv interface{}, rv interface{}, op string) bool {
	fmt.Println(lv)
	fmt.Println(rv)

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
			logrus.Error(err)
			// handle err
		} else {
			return ret
		}

	case "=~":
		// trans to string

		if r, err := regexp.Compile(rv); err != nil {
			logrus.Error(err)
		} else {
			return r.MatchString(lv)
		}
	}
	// error op not supported
	logrus.Error("error: filter expression op{" + op + "} not supported")
	return false
}
