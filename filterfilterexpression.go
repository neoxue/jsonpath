package jsonpath

import (
	"strconv"
	"strings"
)

//filter expression
//?($.name == "xue")

type filterFilterExpression struct {
	t *pathtoken
}

// TODO action set and unset
func (f *filterFilterExpression) eval(action string, cv interface{}, optionalValue interface{}) ([]interface{}, bool) {
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

// only support limit scripts
func (f *filterFilterExpression) getIndexes(v string, length int) ([]int, bool) {
	slice := strings.Split(v, ":")
	start, _ := strconv.Atoi(slice[0])
	end, _ := strconv.Atoi(slice[1])
	step, _ := strconv.Atoi(slice[2])
	k := start
	ks := []int{k}
	if step < 1 {
		step = 1
	}
	if end == 0 {
		end = length
	}

	for true {
		if k+step < end {
			k += step
			ks = append(ks, k)
		}
	}
	return ks, true
}
