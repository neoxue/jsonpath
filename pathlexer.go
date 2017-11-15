package jsonpath

import (
	"regexp"
	"strconv"
	"strings"

	"github.com/pkg/errors"
)

/*
 TODO to be decided, whether we chould use MATCH_KEY(for only string) diffrent to MATCH_INDEX
*/

const (
	MATCH_INDEX                  = `^(\w+|\*)$`         // foo
	MATCH_INDEXES                = `^(\s*\w+[\w,\s]+)$` // 0,1,2
	MATCH_SLICE                  = `^([-\d:]+|:)$`      // [0:2:1]
	MATCH_QUERY_RESULT           = `^\s*\(.+?\)\s*$`    // ?(@.length - 1)
	MATCH_QUERY_MATCH            = `^\s*\?\(.+?\)\s*$`  // ?(@.foo = "bar")
	MATCH_INDEX_IN_SINGLE_QUETES = `^\s*'(.+?)'\s*$`    // 'bar'
	MATCH_INDEX_IN_DOUBLE_QUETES = `^\s*"(.+?)"\s*$`    // "bar"
)

type pathlexer struct {
	Expr string
}

func newLexer(expression string) (pathlexer, error) {
	if len(expression) < 1 {
		return pathlexer{}, errors.New("lexer error: expression empty")
	}
	if expression[0] == '$' {
		expression = expression[1:]
	}
	if expression[0] != '.' && expression[0] != '[' {
		return pathlexer{}, errors.New("lexer error: the second char is not '.' or '['")
	}
	return pathlexer{Expr: expression}, nil
}

func (lexer *pathlexer) ParseExpressionTokens() ([]pathtoken, error) {
	var (
		squareBraketDepth = 0
		tokenValue        = ""
		length            = len(lexer.Expr)

		t      pathtoken
		err    error
		tokens []pathtoken
	)
	tokens = make([]pathtoken, 0)
	for i := 0; i < length; i++ {
		char := lexer.Expr[i]
		if squareBraketDepth == 0 && char == '.' {
			if lexer.lookAhead(i, 1) == '.' {
				if t, err = newToken(T_RECURSIVE, ""); err != nil {
					return tokens, err
				}
				tokens = append(tokens, t)
			}
			continue
		}
		if char == '[' {
			squareBraketDepth += 1
			if squareBraketDepth == 1 {
				continue
			}
		}
		if char == ']' {
			squareBraketDepth -= 1
			if squareBraketDepth == 0 {
				continue
			}
		}
		if squareBraketDepth > 0 {
			tokenValue += string(char)
			if lexer.lookAhead(i, 1) == ']' && squareBraketDepth == 1 {
				if t, err = lexer.CreateToken(tokenValue); err != nil {
					return []pathtoken{}, err
				}
				tokens = append(tokens, t)
				tokenValue = ""
			}
		}
		if squareBraketDepth == 0 {
			tokenValue += string(char)
			// double dot ".."
			if char == '.' {
				if t, err = lexer.CreateToken(tokenValue); err != nil {
					return []pathtoken{}, err
				}
				tokens = append(tokens, t)
				continue
			}

			if lexer.lookAhead(i, 1) == '.' || lexer.lookAhead(i, 1) == '[' || lexer.atEnd(i) {
				if t, err = lexer.CreateToken(tokenValue); err != nil {
					return []pathtoken{}, err
				}
				tokenValue = ""
				tokens = append(tokens, t)
			}
		}
	}
	if tokenValue != "" {
		if t, err = lexer.CreateToken(tokenValue); err != nil {
			return []pathtoken{}, err
		}
		tokens = append(tokens, t)
	}
	return tokens, err
}

func (lexer *pathlexer) atEnd(i int) bool {
	return i+1 == len(lexer.Expr)

}
func (lexer *pathlexer) lookAhead(pos int, forward int) byte {
	if pos+forward >= len(lexer.Expr) {
		return 0
	}
	return lexer.Expr[pos+forward]
}

func (lexer *pathlexer) CreateToken(value string) (pathtoken, error) {
	var (
		matched bool
		err     error
		a       string
		r       *regexp.Regexp
		i       int
		v       string
		vi      int
	)
	if matched, err = regexp.Match(MATCH_INDEX, []byte(value)); matched == true {
		return newToken(T_INDEX, value)
	}

	if matched, err = regexp.Match(MATCH_INDEXES, []byte(value)); matched == true {
		list := strings.Split(value, ",")
		listi := make([]string, 0)
		for _, v := range list {
			v = strings.TrimSpace(v)
			listi = append(listi, v)
		}
		return newToken(T_INDEXES, listi)
	}

	if matched, err = regexp.Match(MATCH_SLICE, []byte(value)); matched == true {
		a := make(map[string]int, 3)
		parts := strings.Split(value, ":")
		var word string
		for i, v = range parts {
			switch i {
			case 0:
				word = "start"
			case 1:
				word = "end"
			case 2:
				word = "step"
			default:
				continue
			}
			if parts[i] == "" {
				a[word] = 0
			} else {
				v = strings.TrimSpace(v)
				if vi, err = strconv.Atoi(v); err != nil {
					return pathtoken{}, errors.Wrap(err, "unable to parse pathtoken {"+value+"}, strconv.Atoi failed in expression:"+lexer.Expr)
				}
				a[word] = vi
			}
		}
		return newToken(T_SLICE, a)
	}

	if matched, err = regexp.Match(MATCH_QUERY_RESULT, []byte(value)); matched == true {
		a = value[1 : len(value)-1]
		return newToken(T_QUERY_RESULT, a)
	}
	if matched, err = regexp.Match(MATCH_QUERY_MATCH, []byte(value)); matched == true {
		a = value[2 : len(value)-1]
		return newToken(T_QUERY_MATCH, a)
	}

	r, _ = regexp.Compile(MATCH_INDEX_IN_SINGLE_QUETES)
	if matches := r.FindStringSubmatch(value); len(matches) > 1 {
		return newToken(T_INDEX, matches[1])
	}

	r, _ = regexp.Compile(MATCH_INDEX_IN_DOUBLE_QUETES)
	if matches := r.FindStringSubmatch(value); len(matches) > 1 {
		return newToken(T_INDEX, matches[1])
	}

	return pathtoken{}, errors.New("unable to parse pathtoken {" + value + "} in expression:" + lexer.Expr)
}
