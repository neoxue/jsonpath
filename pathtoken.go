package jsonpath

import (
	"errors"
	"github.com/sirupsen/logrus"
)

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
func (t *pathtoken) buildFilter(root *JsonPath) (filterinterface, error) {
	switch t.typ {
	case tokenIndex:
		return &filterIndex{t: t}, nil
	case tokenIndexes:
		return &filterIndexes{t: t}, nil
	case tokenRecursive:
		return &filterRecursive{t: t}, nil
	case tokenSlice:
		return &filterSlice{t: t}, nil
	case tokenQueryFilterExpression:
		return &filterFilterExpression{t: t, root: root}, nil
	case tokenQueryScript:
		return &filterScript{t: t, root: root}, nil
	}
	logrus.Error(t)
	return nil, errors.New("jsonpath build filter unsupported filter")
}
