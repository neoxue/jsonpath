package jsonpath

type filterIndex struct {
	t *pathtoken
}

func (f *filterIndex) eval(action string, cv interface{}, optionalValue interface{}) ([]interface{}, bool) {
	var (
		ah *access
		ok bool
	)
	if ah, ok = newaccess(cv); !ok {
		return nil, false
	}
	switch action {
	case actionFind:
		return f.find(ah)
	case actionSet:
		if f.t.v == "*" {
			return []interface{}{ah.setAll(optionalValue)}, true
		} else {
			ok = ah.setValue(f.t.v.(string), optionalValue)
			return []interface{}{ok}, ok
		}
	case actionUnset:
		if f.t.v == "*" {
			return []interface{}{ah.unsetAll()}, true
		} else {
			ok = ah.unsetValue(f.t.v.(string))
			return []interface{}{ok}, ok
		}
	}
	return nil, false //should not go here
}

func (f *filterIndex) find(ah *access) ([]interface{}, bool) {
	if f.t.v == "*" {
		return ah.arrayValues(), true
	}
	if v, ok := ah.getValue(f.t.v.(string)); ok {
		return []interface{}{v}, true
	}
	return nil, false
}
