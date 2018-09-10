package jsonpath

import (
	"github.com/kataras/go-errors"
	"strconv"
)

func isquotestring(tmp string) bool {
	return len(tmp) > 1 && tmp[0] == tmp[len(tmp)-1] && (tmp[0] == '"' || tmp[0] == '\'')
}
func isjsonpath(tmp string) bool {
	return tmp[0] == '$' || tmp[0] == '@'
}

func tryConvertFloat64(a interface{}) (float64, bool) {
	switch a.(type) {
	case int8:
		return float64(a.(int8)), true
	case int16:
		return float64(a.(int16)), true
	case int32:
		return float64(a.(int32)), true
	case int64:
		return float64(a.(int64)), true
	case uint8:
		return float64(a.(uint8)), true
	case uint16:
		return float64(a.(uint16)), true
	case uint32:
		return float64(a.(uint32)), true
	case uint64:
		return float64(a.(uint64)), true
	case float32:
		return float64(a.(float32)), true
	case float64:
		return a.(float64), true
	case bool:
		if true == a.(bool) {
			return 1, true
		} else {
			return 0, true
		}
	case string:
		if got, err := strconv.ParseFloat(a.(string), 64); err == nil {
			return got, true
		} else {
			return 0, false
		}
	}
	return 0, false
}

// verify whether is a number (string)
func isnumberstring(a string) bool {
	if _, err := strconv.Atoi(a); err == nil {
		return true
	}
	if _, err := strconv.ParseFloat(a, 64); err == nil {
		return true
	}
	return false
}

// rules:
// number -> convert
// string -> compare one byte by one byte
func compare_valstring(lv string, rv string, op string) (bool, error) {
	lv1, err1 := strconv.Atoi(lv)
	rv1, err2 := strconv.Atoi(rv)
	if err1 == nil && err2 == nil {
		return compareInt(lv1, rv1, op)
	}
	lv2, err3 := strconv.ParseFloat(lv, 64)
	rv2, err4 := strconv.ParseFloat(rv, 64)
	if err3 == nil && err4 == nil {
		return compareFloat(lv2, rv2, op)
	}
	return compareString(lv, rv, op)
}

//generic is necessary......
func compareString(lv string, rv string, op string) (bool, error) {
	switch op {
	case "==":
		return lv == rv, nil
	case "!=":
		return lv != rv, nil
	case "<":
		return lv < rv, nil
	case ">":
		return lv > rv, nil
	case ">=":
		return lv >= rv, nil
	case "<=":
		return lv <= rv, nil
	default:
		return false, errors.New("compare int operator {" + op + "} not supported")
	}
}

func compareInt(lv int, rv int, op string) (bool, error) {
	switch op {
	case "==":
		return lv == rv, nil
	case "!=":
		return lv != rv, nil
	case "<":
		return lv < rv, nil
	case ">":
		return lv > rv, nil
	case ">=":
		return lv >= rv, nil
	case "<=":
		return lv <= rv, nil
	default:
		return false, errors.New("compare int operator {" + op + "} not supported")
	}
}
func compareFloat(lv float64, rv float64, op string) (bool, error) {
	switch op {
	case "==":
		return lv == rv, nil
	case "!=":
		return lv != rv, nil
	case "<":
		return lv < rv, nil
	case ">":
		return lv > rv, nil
	case ">=":
		return lv >= rv, nil
	case "<=":
		return lv <= rv, nil
	default:
		return false, errors.New("compare float operator {" + op + "} not supported")
	}
}
