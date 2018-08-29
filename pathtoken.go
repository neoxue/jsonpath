package jsonpath

import "errors"

const (
	tokenIndex                 = "index"
	tokenRecursive             = "recursive"
	tokenQueryFilterExpression = "queryFilterExpression"
	tokenQueryScript           = "queryScript"
	tokenSlice                 = "slice"
	tokenIndexes               = "indexes"
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
	list := []string{tokenIndex, tokenRecursive, tokenQueryScript, tokenQueryFilterExpression, tokenSlice, tokenIndexes}
	for _, b := range list {
		if t == b {
			return nil
		}
	}
	return errors.New("Invalid pathtoken: " + t)
}
func (t *pathtoken) buildFilter() filterinterface {
	switch t.typ {
	case tokenIndex:
		return &filterIndex{t: t}
	}
	return nil
}
