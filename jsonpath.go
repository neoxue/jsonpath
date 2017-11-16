package jsonpath

import (
	"crypto/md5"
	"encoding/hex"
)

/*

Flow/JSONPath
*/
const (
	ACTION_SET   = "set"
	ACTION_UNSET = "unset"
	ACTION_FIND  = "find"
)

type JsonPath struct {
	Data interface{}
}

var parsedTokenCache = make(map[string][]pathtoken)

func (jp *JsonPath) UnsetValue(expr string) ([]interface{}, error) {
	jpr := jp.eval("unset", expr, nil)
	return jpr.Collections, jpr.Err
}
func (jp *JsonPath) SetValue(expr string, v interface{}) ([]interface{}, error) {
	jpr := jp.eval("set", expr, v)
	return jpr.Collections, jpr.Err
}
func (jp *JsonPath) Find(expression string) *JsonPathResult {
	return jp.eval("find", expression, nil)
}

func (jp *JsonPath) eval(action string, expression string, optionalValue interface{}) *JsonPathResult {
	var (
		ok            bool
		err           error
		t             pathtoken
		tokens        []pathtoken
		filter        filterbase
		i             int
		cv            interface{}
		collections   []interface{}
		filterData    []interface{}
		filteredValue []interface{}
	)
	if tokens, err = jp.parseTokens(expression); err != nil {
		return &JsonPathResult{Err: err}
	}
	collections = []interface{}{jp.Data}
	for i, t = range tokens {
		filter = t.buildFilter()
		filterData = []interface{}{}
		for _, cv = range collections {
			theAction := "find"
			if i == len(tokens)-1 {
				theAction = action
			}
			if filteredValue, ok = filter.filter(theAction, cv, optionalValue); ok {
				filterData = append(filterData, filteredValue...)
			}
		}
		collections = filterData
	}
	return &JsonPathResult{Collections: collections}
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
