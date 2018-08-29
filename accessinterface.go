package jsonpath

type accessinterface interface {
	getValue(key string) interface{}
	setValue(key string, v interface{}) interface{}
	unsetValue(key string) interface{}
}
