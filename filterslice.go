package jsonpath

type filterSlice struct {
	t *pathtoken
}

func (f *filterSlice) eval(action string, cv interface{}, optionalValue interface{}) ([]interface{}, bool) {
	ah := newaccessins(cv)
	switch action {
	case actionFind:
		return ah.get("*")
	case actionSet:
		return nil, ah.set("*", optionalValue)
	case actionUnset:
		return nil, ah.unset("*")
	default:
		return nil, true
	}
}

func (f *filterRecursive)
