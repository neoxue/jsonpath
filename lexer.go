package jsonpath

import (
	"errors"
	"fmt"
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
	fmt.Println(expression)

	if expression[0] != '.' && expression[0] != '[' {
		return Lexer{}, errors.New("lexer error: the second char is not '.' or '['")
	}
	return Lexer{Expr: expression}, nil
}

func (lexer *Lexer) ParseExpressionTokens() err {
	var (
		dotIndexDepth     = 0
		squareBraketDepth = 0
		capturing         = false
		tokenValue        = ""
		tokens            = [...]Token{}

		length = len(lexer.Expr)
		token  Token
		err    error
	)
	for i := 0; i < length; i++ {
		char := lexer.Expr[i]
		if squareBraketDepth == 0 {
			if char == '.' {
				if lexer.lookAhead(i, 1) == '.' {
					if token, err = NewToken(T_RECURSIVE, ""); err != nil {
						append(tokens, token)
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
					if token, err = lexer.createToken(tokenValue); err != nil {
						return nil, err
					}
					append(tokens, token)
					tokenValue = ""
				}
			}
		}

	}

}

func (lexer *Lexer) lookAhead(pos int, forward int) byte {
	return lexer.Expr[pos+forward]
}

func (lexer *Lexer) createToken(tokenValue string) {

}
