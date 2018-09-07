package jsonpath

import "strconv"

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
func compare_valstring(lv string, rv string, op string) bool {
	lv1, err1 := strconv.Atoi(lv)
	rv1, err2 := strconv.Atoi(rv)
	if err1 == nil && err2 == nil {
		return compareInt()

	}
	if isnumber(lv) && isnumber(rv) {
		switch op {
		case "==":
			return lv == rv
		case "!=":
			return lv != rv
		case "<":
			return lv < rv
		case ">":
			return lv > rv
		case ">=":
			return lv >= rv
		case "<=":
			return lv <= rv
		}

	}

}
