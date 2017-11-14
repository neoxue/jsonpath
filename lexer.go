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
	MATCH_INDEXES                = `^(\s*\d+[\d,\s]+)$` // 0,1,2
	MATCH_SLICE                  = `^([-\d:]+|:)$`      // [0:2:1]
	MATCH_QUERY_RESULT           = `^\s*\(.+?\)\s*$`    // ?(@.length - 1)
	MATCH_QUERY_MATCH            = `^\s*\?\(.+?\)\s*$`  // ?(@.foo = "bar")
	MATCH_INDEX_IN_SINGLE_QUETES = `^\s*'(.+?)'\s*$`    // 'bar'
	MATCH_INDEX_IN_DOUBLE_QUETES = `^\s*"(.+?)"\s*$`    // "bar"
)

type Lexer struct {
	Expr string
}

func NewLexer(expression string) (Lexer, error) {
	if len(expression) < 1 {
		return Lexer{}, errors.New("lexer error: expression empty")
	}
	if expression[0] == '$' {
		expression = expression[1:]
	}
	if expression[0] != '.' && expression[0] != '[' {
		return Lexer{}, errors.New("lexer error: the second char is not '.' or '['")
	}
	return Lexer{Expr: expression}, nil
}

func (lexer *Lexer) ParseExpressionTokens() ([]Token, error) {
	var (
		squareBraketDepth = 0
		tokenValue        = ""
		length            = len(lexer.Expr)

		token  Token
		err    error
		tokens []Token
	)
	tokens = make([]Token, 0)
	for i := 0; i < length; i++ {
		char := lexer.Expr[i]
		if squareBraketDepth == 0 && char == '.' {
			if lexer.lookAhead(i, 1) == '.' {
				if token, err = NewToken(T_RECURSIVE, ""); err != nil {
					return tokens, err
				}
				tokens = append(tokens, token)
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
				if token, err = lexer.CreateToken(tokenValue); err != nil {
					return []Token{}, err
				}
				tokens = append(tokens, token)
				tokenValue = ""
			}
		}
		if squareBraketDepth == 0 {
			tokenValue += string(char)
			// double dot ".."
			if char == '.' {
				if token, err = lexer.CreateToken(tokenValue); err != nil {
					return []Token{}, err
				}
				tokens = append(tokens, token)
				continue
			}

			if lexer.lookAhead(i, 1) == '.' || lexer.lookAhead(i, 1) == '[' || lexer.atEnd(i) {
				if token, err = lexer.CreateToken(tokenValue); err != nil {
					return []Token{}, err
				}
				tokenValue = ""
				tokens = append(tokens, token)
			}
		}
	}
	if tokenValue != "" {
		if token, err = lexer.CreateToken(tokenValue); err != nil {
			return []Token{}, err
		}
		tokens = append(tokens, token)
	}
	return tokens, err
}

func (lexer *Lexer) atEnd(i int) bool {
	return i+1 == len(lexer.Expr)

}
func (lexer *Lexer) lookAhead(pos int, forward int) byte {
	if pos+forward >= len(lexer.Expr) {
		return 0
	}
	return lexer.Expr[pos+forward]
}

func (lexer *Lexer) CreateToken(value string) (Token, error) {
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
		return NewToken(T_INDEX, value)
	}

	if matched, err = regexp.Match(MATCH_INDEXES, []byte(value)); matched == true {
		list := strings.Split(value, ",")
		listi := make([]int, 0)
		for _, v := range list {
			v = strings.TrimSpace(v)
			if vi, err = strconv.Atoi(v); err != nil {
				return Token{}, errors.Wrap(err, "unable to parse token {"+value+"}, strconv.Atoi failed in expression:"+lexer.Expr)
			}
			listi = append(listi, vi)
		}
		return NewToken(T_INDEXES, listi)
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
					return Token{}, errors.Wrap(err, "unable to parse token {"+value+"}, strconv.Atoi failed in expression:"+lexer.Expr)
				}
				a[word] = vi
			}
		}
		return NewToken(T_SLICE, a)
	}

	if matched, err = regexp.Match(MATCH_QUERY_RESULT, []byte(value)); matched == true {
		a = value[1 : len(value)-1]
		return NewToken(T_QUERY_RESULT, a)
	}
	if matched, err = regexp.Match(MATCH_QUERY_MATCH, []byte(value)); matched == true {
		a = value[2 : len(value)-1]
		return NewToken(T_QUERY_MATCH, a)
	}

	r, _ = regexp.Compile(MATCH_INDEX_IN_SINGLE_QUETES)
	if matches := r.FindStringSubmatch(value); len(matches) > 1 {
		return NewToken(T_INDEX, matches[1])
	}

	r, _ = regexp.Compile(MATCH_INDEX_IN_DOUBLE_QUETES)
	if matches := r.FindStringSubmatch(value); len(matches) > 1 {
		return NewToken(T_INDEX, matches[1])
	}

	return Token{}, errors.New("unable to parse token {" + value + "} in expression:" + lexer.Expr)
}
