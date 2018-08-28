package jsonpath

import (
	"crypto/md5"
	"encoding/hex"
)

const (
	actionSet   = "set"
	actionUnset = "unset"
	actionFind  = "find"
)

//	type JsonPath is a struct
type JsonPath struct {
	Data interface{}
}

var parsedTokenCache = make(map[string][]pathtoken)

func (jp *JsonPath) UnsetValue(expr string) ([]interface{}, error) {
	jpr := jp.eval(actionUnset, expr, nil)
	return jpr.Collection, jpr.Err
}
func (jp *JsonPath) SetValue(expr string, v interface{}) ([]interface{}, error) {
	jpr := jp.eval(actionSet, expr, v)
	return jpr.Collection, jpr.Err
}
func (jp *JsonPath) Find(expr string) *JsonPathResult {
	return jp.eval(actionFind, expr, nil)
}

func (jp *JsonPath) eval(action string, expression string, optionalValue interface{}) *JsonPathResult {
	var (
		err        error
		t          pathtoken
		tokens     []pathtoken
		filter     filterinterface
		i          int
		cv         interface{}
		collection []interface{}
	)
	if tokens, err = jp.parseTokens(expression); err != nil {
		return &JsonPathResult{Err: err}
	}
	collection = []interface{}{jp.Data}
	for i, t = range tokens {
		filter = t.buildFilter()
		values := []interface{}{}
		for _, cv = range collection {
			theAction := "find"
			if i == len(tokens)-1 {
				theAction = action
			}
			if filteredValue, ok := filter.filter(theAction, cv, optionalValue); ok {
				values = append(values, filteredValue...)
			}
		}
		collection = values
	}
	return &JsonPathResult{Collection: collection}
}

func (jp *JsonPath) parseTokens(expr string) ([]pathtoken, error) {
	var (
		tokens []pathtoken
		err    error
		ok     bool
		lexer  pathlexer
	)
	md5Ctx := md5.New()
	md5Ctx.Write([]byte(expr))
	cacheKeyMD5 := md5Ctx.Sum(nil)
	cacheKeyStr := hex.EncodeToString(cacheKeyMD5)
	if tokens, ok = parsedTokenCache[cacheKeyStr]; ok {
		return tokens, nil
	}
	if lexer, err = newLexer(expr); err != nil {
		return nil, err
	}
	if tokens, err = lexer.ParseExpressionTokens(); err != nil {
		return nil, err
	}

	parsedTokenCache[cacheKeyStr] = tokens
	return tokens, nil
}
