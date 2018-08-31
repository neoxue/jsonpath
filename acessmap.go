package jsonpath

type accessmap struct {
	v map[string]interface{}
}

func (ah *accessmap) get(key string) ([]interface{}, bool) {
	if key == "*" {
		arr := []interface{}{}
		for _, v := range ah.v {
			arr = append(arr, v)
		}
		return arr, true
	}
	if v, ok := ah.v[key]; ok {
		return []interface{}{v}, true
	}
	return nil, false
}

func (ah *accessmap) set(key string, v interface{}) bool {
	if key == "*" {
		for k := range ah.v {
			ah.v[k] = v
		}
		return true
	}
	ah.v[key] = v
	return true
}

func (ah *accessmap) unset(key string) bool {
	if key == "*" {
		for k := range ah.v {
			delete(ah.v, k)
		}
		return true
	}
	delete(ah.v, key)
	return true
}

func (ah *accessmap) getByList(keys interface{}) ([]interface{}, bool) {
	var vs = []interface{}{}
	switch keys.(type) {
	case []string:
		for _, k := range keys.([]string) {
			if v, ok := ah.v[k]; ok {
				vs = append(vs, v)
			}
		}
	}
	return vs, true
}

// TODO
func (ah *accessmap) setByList(keys []string) bool {

	return true
}

// TODO
func (ah *accessmap) unsetByList(keys []string) bool {
	return true

}
