package jsonpath

import "errors"

// this part should be rewrited to use a finite-state-machine
// ofcourse use cache

// lp op rp
type expression struct {
	sentence string
	items    []struct {
		str string
		typ string
	}
}

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
			continue
		}
		if c == ' ' {
			if len(stack) == 0 {
				strarr[a] = tmp
				tmp = ""
			} else {
				tmp += string(c)
			}
			continue
		}

		// is op, then
		if isOp(c) {
			if tmp == "" {
				tmp += string(c)
			}
			if isOp
			continue
		}
		// normal characters
		if isNormalCharacter(c) {
			if isOperator(tmp) {

			}
			tmp += string(c)
			continue
		}
		return errors.New("expression contains unaccepted characters:" + string(c))
	}
	return nil
}

// verify whether is a number
func isnumber(a interface{}) bool {
	return true
}

func verifyJsonPathStartChar(a interface{}) bool {
	return a == '$' || a == '@'
}

// special character  ! = ~ <>
func isOpCharacter(c int32) bool {
	return c == '<' || c == '=' || c == '>' || c == '!' || c == '~'
}
func isOperator(tmp string) bool {
	if tmp != "" {
		return isOpCharacter(tmp[0])
	}
	return false
}

// normal characters
func isNormalCharacter(c int32) bool {
	return true
}
