package jsonpath

import (
	"crypto/md5"
	"encoding/hex"
)

/*

Flow/JSONPath
*/

type JsonPath struct {
	Data interface{}
}

var parsedTokenCache = make(map[string][]pathtoken)

func (jp *JsonPath) Find(expression string) *JsonPathResult {
	var (
		tokens        []pathtoken
		err           error
		t             pathtoken
		cv            interface{}
		collections   []interface{}
		filter        filterbase
		filterData    []interface{}
		filteredValue []interface{}
		ok            bool
	)
	if tokens, err = jp.parseTokens(expression); err != nil {
		return &JsonPathResult{Err: err}
	}
	collections = []interface{}{jp.Data}
	for _, t = range tokens {
		filter = t.buildFilter()
		filterData = []interface{}{}
		for _, cv = range collections {
			if filteredValue, ok = filter.filter(cv); ok {
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
