package jsonpath

import (
	"strconv"
	"strings"
)

// @.length -1
type filterScript struct {
	t *pathtoken
}

// TODO action set and unset
func (f *filterScript) eval(action string, cv interface{}, optionalValue interface{}) ([]interface{}, bool) {
	f.doScript(cv)
	ah := newaccessins(cv)
	length := 0
	switch cv.(type) {
	case []interface{}:
		length = len(cv.([]interface{}))
	}
	ks, _ := f.getIndexes(f.t.v, length)
	switch action {
	case actionFind:
		return ah.getByList(ks)
	case actionSet:
		return nil, false
	case actionUnset:
		return nil, false
	default:
		return nil, true
	}
}

func (f *filterScript) doScript(cv interface{}) []string {
	if f.t.v[0] == '$'


}








