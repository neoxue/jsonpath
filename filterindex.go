package jsonpath

type filterIndex struct {
	t *pathtoken
}

func (f *filterIndex) eval(action string, cv interface{}, optionalValue interface{}) ([]interface{}, bool) {
	ah := newaccessins(cv)
	k := f.t.v.(string)
	switch action {
	case actionFind:
		return ah.get(k)
	case actionSet:
		return nil, ah.set(k, optionalValue)
	case actionUnset:
		return nil, ah.unset(k)
	default:
		return nil, true
	}
}
