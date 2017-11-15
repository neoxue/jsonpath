package jsonpath

type filterIndex struct {
	t *pathtoken
}

func (f *filterIndex) filter(collection interface{}) ([]interface{}, bool) {
	var (
		ah *access
		ok bool
		v  interface{}
	)
	if ah, ok = newaccess(collection); !ok {
		return nil, false
	}
	if f.t.v == "*" {
		return ah.arrayValues(), true
	} else {
		if v, ok = ah.getValue(f.t.v.(string)); ok {
			return []interface{}{v}, true
		} else {
			return nil, false
		}
	}
}
