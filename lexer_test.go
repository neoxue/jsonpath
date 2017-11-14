package jsonpath

import (
	"testing"

	"github.com/pkg/errors"

	"github.com/stretchr/testify/assert"
)

func TestNewLexer(t *testing.T) {
	lexer, err := NewLexer("$.a")
	if err != nil {
		t.Error(err)
	}
	a := Lexer{Expr: ".a"}
	if lexer != a {
		t.Error("lexer not expected")
	}
}

func Test_index_wildcard(t *testing.T) {
	lexer, _ := NewLexer("$.*")
	tokens, _ := lexer.ParseExpressionTokens()
	assert.True(t, T_INDEX == tokens[0].Type)
	assert.True(t, "*" == tokens[0].Value)
}

func Test_Index_Simple(t *testing.T) {
	lexer, _ := NewLexer("$.foo")
	tokens, _ := lexer.ParseExpressionTokens()
	assert.True(t, T_INDEX == tokens[0].Type)
	assert.True(t, "foo" == tokens[0].Value)
}
func Test_Index_Recursive(t *testing.T) {
	lexer, _ := NewLexer("$..teams.*")
	tokens, err := lexer.ParseExpressionTokens()
	assert.Equal(t, nil, err)
	assert.Equal(t, 3, len(tokens))
	assert.Equal(t, T_RECURSIVE, tokens[0].Type)
	assert.Equal(t, "", tokens[0].Value)
	assert.Equal(t, T_INDEX, tokens[1].Type)
	assert.Equal(t, "teams", tokens[1].Value)
	assert.Equal(t, T_INDEX, tokens[2].Type)
	assert.Equal(t, "*", tokens[2].Value)
}

func Test_Index_Complex(t *testing.T) {
	lexer, _ := NewLexer(`$["'b.^*_'"]`)
	tokens, err := lexer.ParseExpressionTokens()
	assert.Equal(t, nil, err)
	assert.Equal(t, T_INDEX, tokens[0].Type)
	assert.Equal(t, `'b.^*_'`, tokens[0].Value)
}
func Test_Index_BadlyFormed(t *testing.T) {
	lexer, _ := NewLexer(`$.hello*`)
	_, err := lexer.ParseExpressionTokens()
	err_e := errors.New(`unable to parse token {hello*} in expression:.hello*`)
	assert.Equal(t, err_e.Error(), err.Error())
}

func Test_Index_Integer(t *testing.T) {
	lexer, _ := NewLexer(`$[0]`)
	tokens, err := lexer.ParseExpressionTokens()
	assert.Equal(t, nil, err)
	assert.Equal(t, T_INDEX, tokens[0].Type)
	assert.Equal(t, "0", tokens[0].Value)
}

func Test_Index_IntegerAfterDotNotation(t *testing.T) {
	lexer, _ := NewLexer(`$.books[0]`)
	tokens, err := lexer.ParseExpressionTokens()
	assert.Equal(t, nil, err)
	assert.Equal(t, T_INDEX, tokens[0].Type)
	assert.Equal(t, "books", tokens[0].Value)
	assert.Equal(t, T_INDEX, tokens[1].Type)
	assert.Equal(t, "0", tokens[1].Value)
}

func Test_Index_Word(t *testing.T) {
	lexer, _ := NewLexer(`$.books["foo$-/'"]`)
	tokens, err := lexer.ParseExpressionTokens()
	assert.Equal(t, nil, err)
	assert.Equal(t, T_INDEX, tokens[0].Type)
	assert.Equal(t, "books", tokens[0].Value)
	assert.Equal(t, T_INDEX, tokens[1].Type)
	assert.Equal(t, "foo$-/'", tokens[1].Value)
}

func Test_Index_WordWithWhitespace(t *testing.T) {
	lexer, _ := NewLexer(`$.books[     "foo$-/'"    ]`)
	tokens, err := lexer.ParseExpressionTokens()
	assert.Equal(t, nil, err)
	assert.Equal(t, T_INDEX, tokens[0].Type)
	assert.Equal(t, "books", tokens[0].Value)
	assert.Equal(t, T_INDEX, tokens[1].Type)
	assert.Equal(t, "foo$-/'", tokens[1].Value)
}
func Test_Slice_Simple(t *testing.T) {
	lexer, _ := NewLexer(`$.books[0:1:2]`)
	tokens, err := lexer.ParseExpressionTokens()
	assert.Equal(t, nil, err)
	assert.Equal(t, T_INDEX, tokens[0].Type)
	assert.Equal(t, "books", tokens[0].Value)
	assert.Equal(t, T_SLICE, tokens[1].Type)
	expected := make(map[string]int, 3)
	expected[`start`] = 0
	expected[`end`] = 1
	expected[`step`] = 2
	assert.Equal(t, expected, tokens[1].Value)
}
func Test_Slice_NegativeIndex(t *testing.T) {
	lexer, _ := NewLexer(`$[-1]`)
	tokens, err := lexer.ParseExpressionTokens()
	assert.Equal(t, nil, err)
	assert.Equal(t, T_SLICE, tokens[0].Type)
	expected := make(map[string]int, 3)
	expected[`start`] = -1
	assert.Equal(t, expected, tokens[0].Value)
}

func Test_Slice_AllNull(t *testing.T) {
	lexer, _ := NewLexer(`$[:]`)
	tokens, err := lexer.ParseExpressionTokens()
	assert.Equal(t, nil, err)
	assert.Equal(t, T_SLICE, tokens[0].Type)
	expected := make(map[string]int, 3)
	expected[`start`] = 0
	expected[`end`] = 0
	assert.Equal(t, expected, tokens[0].Value)
}

func Test_Slice_QueryResult_Simple(t *testing.T) {
	lexer, _ := NewLexer(`$[(@.foo + 2)]`)
	tokens, err := lexer.ParseExpressionTokens()
	assert.Equal(t, nil, err)
	assert.Equal(t, T_QUERY_RESULT, tokens[0].Type)
	assert.Equal(t, `@.foo + 2`, tokens[0].Value)
}

func Test_Slice_QueryMatch_Simple(t *testing.T) {
	lexer, _ := NewLexer(`$[?(@.foo < 'bar')]`)
	tokens, err := lexer.ParseExpressionTokens()
	assert.Equal(t, nil, err)
	assert.Equal(t, T_QUERY_MATCH, tokens[0].Type)
	assert.Equal(t, `@.foo < 'bar'`, tokens[0].Value)
}

func Test_Slice_QueryMatch_NotEqualTo(t *testing.T) {
	lexer, _ := NewLexer(`$[?(@.foo != 'bar')]`)
	tokens, err := lexer.ParseExpressionTokens()
	assert.Equal(t, nil, err)
	assert.Equal(t, T_QUERY_MATCH, tokens[0].Type)
	assert.Equal(t, `@.foo != 'bar'`, tokens[0].Value)
}

func Test_Slice_QueryMatch_Brackets(t *testing.T) {
	lexer, _ := NewLexer(`$[?(@['@language']='en')]`)
	tokens, err := lexer.ParseExpressionTokens()
	assert.Equal(t, nil, err)
	assert.Equal(t, T_QUERY_MATCH, tokens[0].Type)
	assert.Equal(t, `@['@language']='en'`, tokens[0].Value)
}

func Test_Recursive_Simple(t *testing.T) {
	lexer, _ := NewLexer(`$..foo`)
	tokens, err := lexer.ParseExpressionTokens()
	assert.Equal(t, nil, err)
	assert.Equal(t, T_RECURSIVE, tokens[0].Type)
	assert.Equal(t, T_INDEX, tokens[1].Type)
	assert.Equal(t, ``, tokens[0].Value)
	assert.Equal(t, `foo`, tokens[1].Value)
}

func Test_Recursive_Wildcard(t *testing.T) {
	lexer, _ := NewLexer(`$..*`)
	tokens, err := lexer.ParseExpressionTokens()
	assert.Equal(t, nil, err)
	assert.Equal(t, T_RECURSIVE, tokens[0].Type)
	assert.Equal(t, T_INDEX, tokens[1].Type)
	assert.Equal(t, ``, tokens[0].Value)
	assert.Equal(t, `*`, tokens[1].Value)
}
func Test_Recursive_BadlyFormed(t *testing.T) {
	lexer, _ := NewLexer(`$..ba^r`)
	_, err := lexer.ParseExpressionTokens()
	err_expected := errors.New("unable to parse token {ba^r} in expression:..ba^r")
	assert.Equal(t, err_expected.Error(), err.Error())
}

func Test_Indexes_Simple(t *testing.T) {
	lexer, _ := NewLexer(`$[1,2,3]`)
	tokens, err := lexer.ParseExpressionTokens()
	assert.Equal(t, nil, err)
	expected := []int{1, 2, 3}

	assert.Equal(t, T_INDEXES, tokens[0].Type)
	assert.Equal(t, expected, tokens[0].Value)
}

func Test_Indexes_WhiteSpace(t *testing.T) {
	lexer, _ := NewLexer(`$[1,2   ,  3]`)
	tokens, err := lexer.ParseExpressionTokens()
	assert.Equal(t, nil, err)
	expected := []int{1, 2, 3}

	assert.Equal(t, T_INDEXES, tokens[0].Type)
	assert.Equal(t, expected, tokens[0].Value)
}
