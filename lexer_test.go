package jsonpath

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewLexer(t *testing.T) {
	lexer, error := NewLexer("$.a")
	if error != nil {
		t.Error(error)
	}
	a := Lexer{Expr: ".a"}
	if lexer != a {
		t.Error("lexer not expected")
	}
}

func TestLexer_ParseExpressionTokens(t *testing.T) {

}

func Test_index_wildcard(t *testing.T) {
	wildlexer, _ := NewLexer("$.*")
	tokens, _ := wildlexer.ParseExpressionTokens()
	assert.True(t, T_INDEX == tokens[0].Type)
	assert.True(t, "*" == tokens[0].Value)
}
