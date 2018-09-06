package jsonpath

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
func compare_any(lv string, rv string, op string) bool {
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
