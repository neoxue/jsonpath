package jsonpath

import "errors"

// this part should be rewrited to use a finite-state-machine
// ofcourse use cache

// lp op rp
// do not support backslash to escape characters
/*
supported normal characters:
0-9 a-z A-Z
_ . * ( ) []

supported op characters:
+ - * / < > = ! ~
*/
type expression struct {
	sentence string
	items    []struct {
		str string
		typ string
	}
}

func (expr *expression) parse() error {
	if err := expr.parseItems(); err != nil {
		return err
	}
	if err := expr.validateItems(); err != nil {
		return err
	}
	return nil
}

// support only   val op val
//           or   op val
func (expr *expression) validateItems() error {
	if len(expr.items) > 3 {
		return errors.New("expression " + expr.sentence + " does not support 4 or more items")
	}
	// validate typs
	for idx, item := range expr.items {
		if item.typ == "op" && !isAvailableOp(item.str) {
			return errors.New("expression " + expr.sentence + " operator {" + item.str + "} not valid")
		}
		if item.typ == "val" {
			if isnumber(item.str) {
				item.typ = "num"
			} else if isstring(item.str) {
				item.typ = "string"
			} else if isjsonpath(item.str) {
				item.typ = "jsonpath"
			} else {
				return errors.New("expression " + expr.sentence + " val {" + item.str + "} not valid")
			}
		}
		if idx > 0 {
			if item.typ == "op" && expr.items[idx-1].typ == "op" {
				return errors.New("expression " + expr.sentence + " two ops, not valid")
			}
			if item.typ != "op" && expr.items[idx-1].typ != "op" {
				return errors.New("expression " + expr.sentence + " two non-ops, not valid")
			}
		}
	}
	return nil
}

func (expr *expression) parseItems() error {
	stack := []byte{}
	tmp := []byte{}
	for _, c := range ([]byte)(expr.sentence) {
		if c == '\'' || c == '"' {
			tmp = append(tmp, c)
			if len(stack) > 0 {
				if stack[len(stack)-1] != c {
					stack = append(stack, c)
				} else {
					if len(stack) == 1 {
						stack = []byte{}
					} else {
						stack = stack[0 : len(stack)-1]
					}
				}
			}
			continue
		}
		if c == ' ' {
			if len(stack) == 0 {
				typ := "val"
				if isOperator(tmp) {
					typ = "op"
				}
				expr.items = append(expr.items, struct {
					str string
					typ string
				}{str: string(tmp), typ: typ})
				tmp = []byte{}
			} else {
				tmp = append(tmp, c)
			}
			continue
		}

		// is op, then
		if isOpCharacter(c) {
			if len(tmp) > 0 {
				if isOperator(tmp) {
					tmp = append(tmp, c)
				} else {
					expr.items = append(expr.items, struct {
						str string
						typ string
					}{str: string(tmp), typ: "val"})
					tmp = []byte{c}
				}
				continue
			}
			tmp = append(tmp, c)
			continue
		}
		// normal characters
		if isNormalCharacter(c) {
			if isOperator(tmp) {
				expr.items = append(expr.items, struct {
					str string
					typ string
				}{str: string(tmp), typ: "val"})
				tmp = []byte{c}
				continue
			}
			tmp = append(tmp, c)
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

// op character  ! = ~ <>  + -
func isOpCharacter(c byte) bool {
	return c == '<' || c == '=' || c == '>' || c == '!' || c == '~' || c == '+' || c == '-' || c == '*' || c == '/'
}

// operator should contains only opcharacter
func isOperator(tmp []byte) bool {
	if len(tmp) > 0 {
		return isOpCharacter(tmp[0])
	}
	return false
}

// normal characters
func isNormalCharacter(c byte) bool {
	if c == '_' || c == '.' || c < '*' || c == '@' || c == '$' || c == '(' || c == ')' || c == '[' || c == ']' {
		return true
	}
	if (c > 'a' && c < 'z') || (c > 'A' && c < 'Z') || (c > '0' && c < '9') {
		return true
	}
	return false
}

func isAvailableOp(tmp string) bool {
	return tmp == "!" || tmp == "==" || tmp == "<" || tmp == ">" || tmp == "<=" || tmp == ">=" || tmp == "=~" || tmp == "+" || tmp == "-" || tmp == "*" || tmp == "-"
}

func isstring(tmp string) bool {
	return tmp[0] == tmp[len(tmp)-1] && (tmp[0] == '"' || tmp[0] == '\'')
}
func isjsonpath(tmp string) bool {
	return tmp[0] == '$' || tmp[0] == '@'
}
