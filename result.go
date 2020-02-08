package jsonpath

// Result represents return value
type Result struct {
	collection []interface{}
	err        error
}

// First returns the first element of result collection
func (jpr *Result) First() (interface{}, bool, error) {
	if len(jpr.collection) > 0 {
		return jpr.collection[0], true, nil
	}
	return nil, false, jpr.err
}

// All returns collection
func (jpr *Result) All() ([]interface{}, bool, error) {
	return jpr.collection, len(jpr.collection) > 0, jpr.err
}

func (jpr *Result) FilterString() ([]string, bool, error) {
	ret := []string{}
	for _, a := range jpr.collection {
		switch a.(type) {
		case string:
			ret = append(ret, a.(string))
		}
	}
	return ret,len(ret) > 0,jpr.err
}
func (jpr *Result) FilterFloat64() ([]float64, bool, error) {
	ret := []float64{}
	for _, a := range jpr.collection {
		switch a.(type) {
		case float64:
			ret = append(ret, a.(float64))
		}
	}
	return ret,len(ret) > 0,jpr.err
}
