package jsonpath

type filterinterface interface {
	filter(action string, collection interface{}, optionalValue interface{}) ([]interface{}, bool)
}
