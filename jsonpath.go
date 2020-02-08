package jsonpath

import (
	"crypto/md5"
	"encoding/hex"
	"sync"
)

const (
	actionSet   = "set"
	actionUnset = "unset"
	actionFind  = "find"
)

//necessary to decide whether log jsonpath problems;
//JsonPath is exported
type JsonPath struct {
	Data interface{}
}

var parsedTokenCache = make(map[string][]pathtoken)

//UnsetValue unsets json's value and return bool,err
func (jp *JsonPath) UnsetValue(expr string) (bool, error) {
	jpr := jp.eval(actionUnset, expr, nil)
	return jpr.err == nil, jpr.err
}

//SetValue sets json's value and return bool,err
func (jp *JsonPath) SetValue(expr string, v interface{}) (bool, error) {
	jpr := jp.eval(actionSet, expr, v)
	return jpr.err == nil, jpr.err
}

//Find json's value and return JsonResult
func (jp *JsonPath) Find(expr string) *Result {
	return jp.eval(actionFind, expr, nil)
}

func (jp *JsonPath) eval(action string, expression string, optionalValue interface{}) *Result {
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
		return &Result{err: err}
	}
	collection = []interface{}{jp.Data}

	for i, t = range tokens {
		if filter, err = t.buildFilter(jp); err != nil {
			return &Result{err: err}
		}
		values := []interface{}{}
		for _, cv = range collection {
			theAction := actionFind
			if i == len(tokens)-1 {
				theAction = action
			}
			if filteredValue, ok := filter.eval(theAction, cv, optionalValue); ok {
				values = append(values, filteredValue...)
			}
		}
		collection = values
	}
	return &Result{collection: collection}
}

var mu sync.Mutex

func (jp *JsonPath) parseTokens(expr string) ([]pathtoken, error) {
	mu.Lock()
	defer mu.Unlock()
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
	if tokens, err = lexer.parseExpressionTokens(); err != nil {
		return nil, err
	}
	parsedTokenCache[cacheKeyStr] = tokens
	return tokens, nil
}
