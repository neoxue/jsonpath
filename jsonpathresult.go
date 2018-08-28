package jsonpath

type JsonPathResult struct {
	Collection []interface{}
	Err        error
}

func (jpr *JsonPathResult) First() (interface{}, bool, error) {
	if len(jpr.Collection) > 0 {
		return jpr.Collection[0], true, nil
	} else {
		return nil, false, jpr.Err
	}
}
func (jpr *JsonPathResult) All() ([]interface{}, bool, error) {
	return jpr.Collection, len(jpr.Collection) > 0, jpr.Err
}
