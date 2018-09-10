package jsonpath

import (
	"github.com/sirupsen/logrus"
	"reflect"
	"regexp"
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
					logrus.Warn(err)
				} else if got == true {
					if f.doOp(val, rval, op) {
						result = append(result, v)
					}
				} else {
					logrus.Warn("rv jsonpath {" + expr.items[2].val + "} nothing found")
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
	only
	number  -> number compare
	string  -> string compare
	if =~ all trans to string
*/
func (f *filterFilterExpression) doOp(lv interface{}, rv interface{}, op string) bool {
	switch op {
	case "==", "!=", ">", ">=", "<", "<=":
		lv1, ok1 := tryConvertFloat64(lv)
		rv1, ok2 := tryConvertFloat64(rv)
		if ok1 && ok2 {
			if got, err := compareFloat(lv1, rv1, op); err != nil {
				logrus.Warn(err)
				return false
			} else {
				return got
			}
		}
		switch lv.(type) {
		case string:
			switch rv.(type) {
			case string:
				if got, err := compareString(lv.(string), rv.(string), op); err != nil {
					logrus.Warn(err)
					return false
				} else {
					return got
				}
			}
		}
		return false
	case "=~":
		// trans to string
		var lv1, rv1 string
		if reflect.TypeOf(lv) == reflect.TypeOf("string") {
			lv1 = lv.(string)
		} else {
			return false
		}
		if reflect.TypeOf(rv) == reflect.TypeOf("string") {
			rv1 = rv.(string)
		} else {
			return false
		}
		if r, err := regexp.Compile(rv1); err != nil {
			logrus.Warn(err)
		} else {
			return r.MatchString(lv1)
		}
	}
	// error op not supported
	logrus.Warn("error: filter expression op{" + op + "} not supported")
	return false
}
