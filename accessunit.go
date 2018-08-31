package jsonpath

//  could not do anything

type accessunit struct {
	v interface{}
}

func (ah *accessunit) get(key string) ([]interface{}, bool) {
	return nil, false
}
func (ah *accessunit) set(key string, data interface{}) bool {
	return false
}
func (ah *accessunit) unset(key string) bool {
	return false
}

func (ah *accessunit) getByList(keys interface{}) ([]interface{}, bool) {
	return nil, false
}
