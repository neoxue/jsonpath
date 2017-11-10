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

type Token struct {
	Type  string
	Value string
}

func NewToken(t string, v string) (Token, error) {
	if err := validateType(t); err != nil {
		return Token{}, err
	}
	return Token{Type: t, Value: v}, nil
}

func validateType(t string) error {
	list := []string{T_INDEX, T_RECURSIVE, T_QUERY_RESULT, T_QUERY_MATCH, T_SLICE, T_INDEXES}
	for _, b := range list {
		if t == b {
			return nil
		}
	}
	return errors.New("Invalid token: " + t)
}
