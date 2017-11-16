package jsonpath

type filterbase interface {
	filter(action string, collection interface{}, optionalValue interface{}) ([]interface{}, bool)
}
