package jsonpath

import (
	"crypto/md5"
	"encoding/hex"
	"github.com/lestrrat/go-file-rotatelogs"
	"github.com/sirupsen/logrus"
)

const (
	actionSet   = "set"
	actionUnset = "unset"
	actionFind  = "find"
)

func init() {
	logrus.SetFormatter(&logrus.TextFormatter{QuoteEmptyFields: false, ForceColors: true, FullTimestamp: true, DisableColors: false})
	rl, _ := rotatelogs.New("/data1/ms/log/logrus.%Y%m%d")
	logrus.SetOutput(rl)
	logrus.SetLevel(logrus.DebugLevel)
}

//necessary to decide whether log jsonpath problems;
//JsonPath is exported
type JsonPath struct {
	Data interface{}
}

var parsedTokenCache = make(map[string][]pathtoken)

//UnsetValue unsets json's value and return json
func (jp *JsonPath) UnsetValue(expr string) ([]interface{}, error) {
	jpr := jp.eval(actionUnset, expr, nil)
	return jpr.collection, jpr.err
}

//SetValue sets json's value and return json
func (jp *JsonPath) SetValue(expr string, v interface{}) ([]interface{}, error) {
	jpr := jp.eval(actionSet, expr, v)
	return jpr.collection, jpr.err
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
	if tokens, err = lexer.parseExpressionTokens(); err != nil {
		return nil, err
	}
	parsedTokenCache[cacheKeyStr] = tokens
	return tokens, nil
}
