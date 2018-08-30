package jsonpath

import (
	"regexp"
	"strconv"
	"strings"

	"github.com/pkg/errors"
)

/*
 TODO to be decided, whether we should use MATCH_KEY(for only string) diffrent to matchIndex
*/
const (
	matchIndex                 = `^(\w+|\*)$`         // foo
	matchIndexes               = `^(\s*\w+[\w,\s]+)$` // 0,1,2
	matchSlice                 = `^([-\d:]+|:)$`      // [0:2:1]
	matchQueryScript           = `^\s*\(.+?\)\s*$`    // (@.length - 1)   script
	matchQueryFilterExpression = `^\s*\?\(.+?\)\s*$`  // ?(@.foo = "bar")
	matchIndexInSingleQuetes   = `^\s*'(.+?)'\s*$`    // 'bar'
	matchIndexInDoubleQuetes   = `^\s*"(.+?)"\s*$`    // "bar"
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

func (lexer *pathlexer) parseExpressionTokens() ([]pathtoken, error) {
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
				if t, err = newToken(tokenRecursive, ""); err != nil {
					return tokens, err
				}
				tokens = append(tokens, t)
			}
			continue
		}
		if char == '[' {
			squareBraketDepth++
			if squareBraketDepth == 1 {
				continue
			}
		}
		if char == ']' {
			squareBraketDepth--
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
	if matched, err = regexp.Match(matchIndex, []byte(value)); matched == true {
		return newToken(tokenIndex, value)
	}

	if matched, err = regexp.Match(matchIndexes, []byte(value)); matched == true {
		return newToken(tokenIndexes, value)
	}

	if matched, err = regexp.Match(matchSlice, []byte(value)); matched == true {
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
		return newToken(tokenSlice, a)
	}

	if matched, err = regexp.Match(matchQueryScript, []byte(value)); matched == true {
		a = value[1 : len(value)-1]
		return newToken(tokenQueryScript, a)
	}
	if matched, err = regexp.Match(matchQueryFilterExpression, []byte(value)); matched == true {
		a = value[2 : len(value)-1]
		return newToken(tokenQueryFilterExpression, a)
	}

	r, _ = regexp.Compile(matchIndexInSingleQuetes)
	if matches := r.FindStringSubmatch(value); len(matches) > 1 {
		return newToken(tokenIndex, matches[1])
	}

	r, _ = regexp.Compile(matchIndexInDoubleQuetes)
	if matches := r.FindStringSubmatch(value); len(matches) > 1 {
		return newToken(tokenIndex, matches[1])
	}

	return pathtoken{}, errors.New("unable to parse pathtoken {" + value + "} in expression:" + lexer.Expr)
}
