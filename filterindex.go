package jsonpath

type filterIndex struct {
	t *pathtoken
}

func (f *filterIndex) filter(action string, collection interface{}, optionalValue interface{}) ([]interface{}, bool) {
	var (
		ah *access
		ok bool
		v  interface{}
	)
	if ah, ok = newaccess(collection); !ok {
		return nil, false
	}
	switch action {
	case "find":
		if f.t.v == "*" {
			return ah.arrayValues(), true
		} else {
			if v, ok = ah.getValue(f.t.v.(string)); ok {
				return []interface{}{v}, true
			} else {
				return nil, false
			}
		}
	case "set":
		if f.t.v == "*" {
			return []interface{}{ah.setAll(optionalValue)}, true
		} else {
			ok = ah.setValue(f.t.v.(string), optionalValue)
			return []interface{}{ok}, ok
		}
	case "unset":
		if f.t.v == "*" {
			return []interface{}{ah.unsetAll()}, true
		} else {
			ok = ah.unsetValue(f.t.v.(string))
			return []interface{}{ok}, ok
		}
	}
	return nil, false //should not go here
}
