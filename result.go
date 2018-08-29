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
