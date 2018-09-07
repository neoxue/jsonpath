package jsonpath

import (
	"github.com/kataras/go-errors"
	"github.com/labstack/gommon/log"
	"strconv"
)

func isstring(tmp string) bool {
	return tmp[0] == tmp[len(tmp)-1] && (tmp[0] == '"' || tmp[0] == '\'')
}
func isjsonpath(tmp string) bool {
	return tmp[0] == '$' || tmp[0] == '@'
}

// verify whether is a number
func isnumber(a interface{}) bool {
	return true
}

func convertnum(a interface{}) bool {

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

	}

}

func compareInt(lv1 int, rv1 int, op string) (bool, error) {
	switch op {
	case "==":
		return lv1 == rv1, nil
	case "!=":
		return lv1 != rv1, nil
	case "<":
		return lv1 < rv1, nil
	case ">":
		return lv1 > rv1, nil
	case ">=":
		return lv1 >= rv1, nil
	case "<=":
		return lv1 <= rv1, nil
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
