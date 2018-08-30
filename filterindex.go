package jsonpath

type filterIndex struct {
	t *pathtoken
}

func (f *filterIndex) eval(action string, cv interface{}, optionalValue interface{}) ([]interface{}, bool) {
	ah := newaccessins(cv)
	switch action {
	case actionFind:
		return ah.get(f.t.v)
	case actionSet:
		return nil, ah.set(f.t.v, optionalValue)
	case actionUnset:
		return nil, ah.unset(f.t.v)
	default:
		return nil, true
	}
}
