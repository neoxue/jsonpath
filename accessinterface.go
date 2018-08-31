package jsonpath

type accessinterface interface {
	get(key string) ([]interface{}, bool)
	set(key string, v interface{}) bool
	unset(key string) bool
	getByList(keys interface{}) ([]interface{}, bool)
}

func newaccessins(cv interface{}) accessinterface {
	switch cv.(type) {
	case []interface{}:
		return accessinterface(&accessarray{v: cv.([]interface{})})
	case map[string]interface{}:
		return accessinterface(&accessmap{v: cv.(map[string]interface{})})
	default:
		return accessinterface(&accessunit{v: cv})
	}
}
