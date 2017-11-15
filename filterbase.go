package jsonpath

type filterbase interface {
	filter(collection interface{}) ([]interface{}, bool)
}
