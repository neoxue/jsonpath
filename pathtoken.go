package jsonpath

import "errors"

const (
	tokenIndex                 = "index"
	tokenIndexes               = "indexes"
	tokenRecursive             = "recursive"
	tokenQueryFilterExpression = "queryFilterExpression"
	tokenQueryScript           = "queryScript"
	tokenSlice                 = "slice"
)

type pathtoken struct {
	typ string
	v   string
}

func newToken(t string, v string) (pathtoken, error) {
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
	case tokenIndexes:
		return &filterIndexes{t: t}
	case tokenRecursive:
		return &filterRecursive{t: t}
	case tokenSlice:
		return &filterSlice{t: t}
	}
	return nil
}
