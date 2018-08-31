package jsonpath

import "strings"

type filterIndexes struct {
	t *pathtoken
}

func (f *filterIndexes) eval(action string, cv interface{}, optionalValue interface{}) ([]interface{}, bool) {
	ah := newaccessins(cv)
	var tokens []string
	var ok bool
	if tokens, ok = f.getIndexes(f.t.v); !ok {
		return nil, false
	}
	switch action {
	case actionFind:
		return ah.getByList(tokens)
	case actionSet:
		for _, k := range tokens {
			ah.set(k, optionalValue)
		}
		return nil, true
	case actionUnset:
		for _, k := range tokens {
			ah.unset(k)
		}
		return nil, true
	default:
		return nil, true
	}
}

func (f *filterIndexes) getIndexes(v string) ([]string, bool) {
	list := strings.Split(v, ",")
	listi := make([]string, 0)
	for _, v := range list {
		v = strings.TrimSpace(v)
		listi = append(listi, v)
	}
	return listi, true
}
