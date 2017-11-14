package jsonpath

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_a(t *testing.T) {
	lexer, _ := NewLexer(`$[1,2   ,  3]`)
	tokens, err := lexer.ParseExpressionTokens()
	assert.Equal(t, nil, err)
	expected := []int{1, 2, 3}

	assert.Equal(t, expected, tokens[0].Value)

}

func testint() int {
	return 0
}
