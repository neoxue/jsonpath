package jsonpath

type filterinterface interface {
	eval(action string, collection interface{}, optionalValue interface{}) ([]interface{}, bool)
}
