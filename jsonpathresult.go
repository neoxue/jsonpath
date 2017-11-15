package jsonpath

type JsonPathResult struct {
	Collections []interface{}
	Err         error
}

func (jpr *JsonPathResult) First() (interface{}, bool, error) {
	if jpr.Err != nil {
		return nil, false, jpr.Err
	}
	if len(jpr.Collections) > 0 {
		return jpr.Collections[0], true, nil
	} else {
		return nil, false, nil
	}
}
func (jpr *JsonPathResult) All() ([]interface{}, bool, error) {
	if jpr.Err != nil {
		return nil, false, jpr.Err
	}
	if len(jpr.Collections) > 0 {
		return jpr.Collections, true, nil
	} else {
		return nil, false, nil
	}
}
