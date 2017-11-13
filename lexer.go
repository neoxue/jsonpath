package jsonpath

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

const (
	MATCH_INDEX                  = `\W+ | \*`          // foo
	MATCH_INDEXES                = `\s* \d+ [\d,\s]+`  // 0,1,2
	MATCH_SLICE                  = `[-\d:]+ | :`       // [0:2:1]
	MATCH_QUERY_RESULT           = `\s* \( .+? \) \s*` // ?(@.length - 1)
	MATCH_QUERY_MATCH            = `\s* \?\(.+?\) \s*` // ?(@.foo = "bar")
	MATCH_INDEX_IN_SINGLE_QUETES = `\s* ' (.+?) ' \s*` // 'bar'
	MATCH_INDEX_IN_DOUBLE_QUETES = `\s* " (.+?) " \s*` // "bar"
)

type Lexer struct {
	Expr string
}

func NewLexer(expression string) (Lexer, error) {
	if len(expression) < 1 {
		return Lexer{}, errors.New("lexer error: expression empty")
	}
	if expression[0] != '$' {
		return Lexer{}, errors.New("lexer error: the first char is not $")
	}
	expression = expression[1:]
	if expression[0] != '.' && expression[0] != '[' {
		return Lexer{}, errors.New("lexer error: the second char is not '.' or '['")
	}
	return Lexer{Expr: expression}, nil
}

func (lexer *Lexer) ParseExpressionTokens() ([]Token, error) {
	var (
		//dotIndexDepth     = 0
		squareBraketDepth = 0
		//capturing         = false
		tokenValue = ""
		length     = len(lexer.Expr)

		token  Token
		err    error
		tokens []Token
	)
	tokens = make([]Token, 10)
	for i := 0; i < length; i++ {
		char := lexer.Expr[i]
		if squareBraketDepth == 0 {
			if char == '.' {
				if lexer.lookAhead(i, 1) == '.' {
					if token, err = NewToken(T_RECURSIVE, ""); err != nil {
						tokens = append(tokens, token)
					}
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
		}

	}
	fmt.Println(tokens)
	return tokens, err
}

func (lexer *Lexer) lookAhead(pos int, forward int) byte {
	return lexer.Expr[pos+forward]
}

func (lexer *Lexer) CreateToken(value string) (Token, error) {
	var (
		matched bool
		err     error
		a       string
		r       *regexp.Regexp
	)
	if matched, err = regexp.Match(`/^(`+MATCH_INDEX+`)$/x`, []byte(value)); matched == true {
		return NewToken(T_INDEX, value)
	}

	if matched, err = regexp.Match(`/^(`+MATCH_INDEXES+`)$/x`, []byte(value)); matched == true {
		list := strings.Split(value, ",")
		listi := make([]float64, 10)
		fmt.Println(list)
		for _, a = range list {
			b, _ := strconv.ParseFloat(a, 10)
			listi = append(listi, b)
		}
		return NewToken(T_INDEX, listi)
	}

	if matched, err = regexp.Match(`/^(`+MATCH_SLICE+`)$/x`, []byte(value)); matched == true {
		parts := strings.Split(value, ":")
		a := make(map[string]string, 3)
		// TODO int  空
		a[`start`] = parts[0]
		a[`end`] = parts[1]
		a[`step`] = parts[2]
		return NewToken(T_SLICE, a)
	}

	if matched, err = regexp.Match(`/^(`+MATCH_QUERY_RESULT+`)$/x`, []byte(value)); matched == true {
		a = value[1 : len(value)-1]
		fmt.Println(a)
		return NewToken(T_QUERY_RESULT, a)
	}
	if matched, err = regexp.Match(`/^(`+MATCH_QUERY_MATCH+`)$/x`, []byte(value)); matched == true {
		a = value[2 : len(value)-1]
		fmt.Println(a)
		return NewToken(T_QUERY_MATCH, a)
	}

	r, _ = regexp.Compile(`/^(` + MATCH_INDEX_IN_SINGLE_QUETES + `)$x`)
	if matches := r.FindString(value); matches != "" {
		var by byte = matches[1]
		fmt.Println(by)
		return NewToken(T_INDEX, by)
	}

	r, _ = regexp.Compile(`/^(` + MATCH_INDEX_IN_DOUBLE_QUETES + `)$x`)
	if matches := r.FindString(value); matches != "" {
		var by byte = matches[1]
		fmt.Println(by)
		return NewToken(T_INDEX, by)
	}

	fmt.Println(err)
	return Token{}, errors.New("unable to parse token {" + value + "} in expression:" + lexer.Expr)
}
