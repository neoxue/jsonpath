package jsonpath

import "errors"

const (
	T_INDEX        = "index"
	T_RECURSIVE    = "recursive"
	T_QUERY_MATCH  = "queryMatch"
	T_QUERY_RESULT = "queryResult"
	T_SLICE        = "slice"
	T_INDEXES      = "indexes"
)

type pathtoken struct {
	typ string
	v   interface{}
}

func newToken(t string, v interface{}) (pathtoken, error) {
	if err := validateType(t); err != nil {
		return pathtoken{}, err
	}
	return pathtoken{typ: t, v: v}, nil
}

func validateType(t string) error {
	list := []string{T_INDEX, T_RECURSIVE, T_QUERY_RESULT, T_QUERY_MATCH, T_SLICE, T_INDEXES}
	for _, b := range list {
		if t == b {
			return nil
		}
	}
	return errors.New("Invalid pathtoken: " + t)
}
func (t *pathtoken) buildFilter() filterinterface {
	switch t.typ {
	case T_INDEX:
		return &filterIndex{t: t}
	}
	return nil
}
