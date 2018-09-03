package jsonpath

// this part should be rewrited to use a finite-state-machine
// ofcourse use cache

// lp op rp
type expression struct {
	sentence string
	lp       interface{}
	op       string
	rp       interface{}
}

// special character  ! = ~ <>
func (expr *expression) parse() error {
	stack := []int32{}
	tmp := ""
	strarr := []string{}
	a := 0
	for idx, c := range expr.sentence {
		if c == '\'' || c == '"' {
			tmp += string(c)
			if idx > 1 && expr.sentence[idx-1] != '\\' {
				continue
			}
			if len(stack) > 0 {
				if stack[len(stack)-1] != c {
					stack = append(stack, c)
				} else {
					if len(stack) == 1 {
						stack = []int32{}
					} else {
						stack = stack[0 : len(stack)-1]
					}
				}
			}
		}
		if c == ' ' {
			if len(stack) == 0 {
				strarr[a] = tmp
				tmp = ""
			} else {
				tmp += string(c)
			}
		}

		if

	}

}

// verify whether is a number
func isnumber(a interface{}) bool {
	return true
}

func verifyJsonPathStartChar(a interface{}) bool {
	return a == '$' || a == '@'
}

func isOp() {

}