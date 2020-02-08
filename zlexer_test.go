package jsonpath

import (
	"fmt"
	"testing"

	"github.com/pkg/errors"

	"github.com/stretchr/testify/assert"
)

func TestNewLexer(t *testing.T) {
	lexer, err := newLexer("$.a")
	if err != nil {
		t.Error(err)
	}
	a := pathlexer{Expr: ".a"}
	if lexer != a {
		t.Error("lexer not expected")
	}
}

func Test_index_wildcard(t *testing.T) {
	lexer, _ := newLexer("$.*")
	tokens, _ := lexer.parseExpressionTokens()
	assert.True(t, tokenIndex == tokens[0].typ)
	assert.True(t, "*" == tokens[0].v)
}

func Test_Index_Simple(t *testing.T) {
	lexer, _ := newLexer("$.foo")
	tokens, _ := lexer.parseExpressionTokens()
	assert.True(t, tokenIndex == tokens[0].typ)
	assert.True(t, "foo" == tokens[0].v)
}
func Test_Index_Recursive(t *testing.T) {
	lexer, _ := newLexer("$..teams.*")
	tokens, err := lexer.parseExpressionTokens()
	assert.Equal(t, nil, err)
	assert.Equal(t, 3, len(tokens))
	assert.Equal(t, tokenRecursive, tokens[0].typ)
	assert.Equal(t, "", tokens[0].v)
	assert.Equal(t, tokenIndex, tokens[1].typ)
	assert.Equal(t, "teams", tokens[1].v)
	assert.Equal(t, tokenIndex, tokens[2].typ)
	assert.Equal(t, "*", tokens[2].v)
}

func Test_Index_Complex(t *testing.T) {
	lexer, _ := newLexer(`$["'b.^*_'"]`)
	tokens, err := lexer.parseExpressionTokens()
	assert.Equal(t, nil, err)
	assert.Equal(t, tokenIndex, tokens[0].typ)
	assert.Equal(t, `'b.^*_'`, tokens[0].v)
}
func Test_Index_BadlyFormed(t *testing.T) {
	lexer, _ := newLexer(`$.hello*`)
	_, err := lexer.parseExpressionTokens()
	erra := errors.New(`unable to parse pathtoken {hello*} in expression:.hello*`)
	assert.Equal(t, erra.Error(), err.Error())
}

func Test_Index_Integer(t *testing.T) {
	lexer, _ := newLexer(`$[0]`)
	tokens, err := lexer.parseExpressionTokens()
	assert.Equal(t, nil, err)
	assert.Equal(t, tokenIndex, tokens[0].typ)
	assert.Equal(t, "0", tokens[0].v)
}

func Test_Index_IntegerAfterDotNotation(t *testing.T) {
	lexer, _ := newLexer(`$.books[0]`)
	tokens, err := lexer.parseExpressionTokens()
	assert.Equal(t, nil, err)
	assert.Equal(t, tokenIndex, tokens[0].typ)
	assert.Equal(t, "books", tokens[0].v)
	assert.Equal(t, tokenIndex, tokens[1].typ)
	assert.Equal(t, "0", tokens[1].v)
}

func Test_Index_Word(t *testing.T) {
	lexer, _ := newLexer(`$.books["foo$-/'"]`)
	tokens, err := lexer.parseExpressionTokens()
	assert.Equal(t, nil, err)
	assert.Equal(t, tokenIndex, tokens[0].typ)
	assert.Equal(t, "books", tokens[0].v)
	assert.Equal(t, tokenIndex, tokens[1].typ)
	assert.Equal(t, "foo$-/'", tokens[1].v)
}

func Test_Index_WordWithWhitespace(t *testing.T) {
	lexer, _ := newLexer(`$.books[     "foo$-/'"    ]`)
	tokens, err := lexer.parseExpressionTokens()
	assert.Equal(t, nil, err)
	assert.Equal(t, tokenIndex, tokens[0].typ)
	assert.Equal(t, "books", tokens[0].v)
	assert.Equal(t, tokenIndex, tokens[1].typ)
	assert.Equal(t, "foo$-/'", tokens[1].v)
}

//func Test_Slice_Simple(t *testing.T) {
//	lexer, _ := newLexer(`$.books[0:1:2]`)
//	tokens, err := lexer.parseExpressionTokens()
//	assert.Equal(t, nil, err)
//	assert.Equal(t, tokenIndex, tokens[0].typ)
//	assert.Equal(t, "books", tokens[0].v)
//	assert.Equal(t, tokenSlice, tokens[1].typ)
//	expected := make(map[string]int, 3)
//	expected[`start`] = 0
//	expected[`end`] = 1
//	expected[`step`] = 2
//	assert.Equal(t, expected, tokens[1].v)
//}
//func Test_Slice_NegativeIndex(t *testing.T) {
//	lexer, _ := newLexer(`$[-1]`)
//	tokens, err := lexer.parseExpressionTokens()
//	assert.Equal(t, nil, err)
//	assert.Equal(t, tokenSlice, tokens[0].typ)
//	expected := make(map[string]int, 3)
//	expected[`start`] = -1
//	assert.Equal(t, expected, tokens[0].v)
//}
//
//func Test_Slice_AllNull(t *testing.T) {
//	lexer, _ := newLexer(`$[:]`)
//	tokens, err := lexer.parseExpressionTokens()
//	assert.Equal(t, nil, err)
//	assert.Equal(t, tokenSlice, tokens[0].typ)
//	expected := make(map[string]int, 3)
//	expected[`start`] = 0
//	expected[`end`] = 0
//	assert.Equal(t, expected, tokens[0].v)
//}

func Test_Slice_QueryResult_Simple(t *testing.T) {
	lexer, _ := newLexer(`$[(@.foo + 2)]`)
	tokens, err := lexer.parseExpressionTokens()
	fmt.Println(tokens)
	assert.Equal(t, errors.New("jsonpath lexer: do not support query script now: {"+`@.foo + 2`+"}").Error(), err.Error())
}

//func Test_Slice_QueryMatch_Simple(t *testing.T) {
//	lexer, _ := newLexer(`$[?(@.foo < 'bar')]`)
//	tokens, err := lexer.parseExpressionTokens()
//	assert.Equal(t, nil, err)
//	assert.Equal(t, tokenQueryFilterExpression, tokens[0].typ)
//	assert.Equal(t, `@.foo < 'bar'`, tokens[0].v)
//}

//func Test_Slice_QueryMatch_NotEqualTo(t *testing.T) {
//	lexer, _ := newLexer(`$[?(@.foo != 'bar')]`)
//	tokens, err := lexer.parseExpressionTokens()
//	assert.Equal(t, nil, err)
//	assert.Equal(t, tokenQueryFilterExpression, tokens[0].typ)
//	assert.Equal(t, `@.foo != 'bar'`, tokens[0].v)
//}
//
//func Test_Slice_QueryMatch_Brackets(t *testing.T) {
//	lexer, _ := newLexer(`$[?(@['@language']='en')]`)
//	tokens, err := lexer.parseExpressionTokens()
//	assert.Equal(t, nil, err)
//	assert.Equal(t, tokenQueryFilterExpression, tokens[0].typ)
//	assert.Equal(t, `@['@language']='en'`, tokens[0].v)
//}

func Test_Recursive_Simple(t *testing.T) {
	lexer, _ := newLexer(`$..foo`)
	tokens, err := lexer.parseExpressionTokens()
	assert.Equal(t, nil, err)
	assert.Equal(t, tokenRecursive, tokens[0].typ)
	assert.Equal(t, tokenIndex, tokens[1].typ)
	assert.Equal(t, ``, tokens[0].v)
	assert.Equal(t, `foo`, tokens[1].v)
}

func Test_Recursive_Wildcard(t *testing.T) {
	lexer, _ := newLexer(`$..*`)
	tokens, err := lexer.parseExpressionTokens()
	assert.Equal(t, nil, err)
	assert.Equal(t, tokenRecursive, tokens[0].typ)
	assert.Equal(t, tokenIndex, tokens[1].typ)
	assert.Equal(t, ``, tokens[0].v)
	assert.Equal(t, `*`, tokens[1].v)
}
func Test_Recursive_BadlyFormed(t *testing.T) {
	lexer, _ := newLexer(`$..ba^r`)
	_, err := lexer.parseExpressionTokens()
	errExpected := errors.New("unable to parse pathtoken {ba^r} in expression:..ba^r")
	assert.Equal(t, errExpected.Error(), err.Error())
}

//func Test_Indexes_Simple(t *testing.T) {
//	lexer, _ := newLexer(`$[1,2,3]`)
//	tokens, err := lexer.parseExpressionTokens()
//	assert.Equal(t, nil, err)
//	expected := []string{"1", "2", "3"}
//
//	assert.Equal(t, tokenIndexes, tokens[0].typ)
//	assert.Equal(t, expected, tokens[0].v)
//}

//func Test_Indexes_WhiteSpace(t *testing.T) {
//	lexer, _ := newLexer(`$[1,2   ,  3]`)
//	tokens, err := lexer.parseExpressionTokens()
//	assert.Equal(t, nil, err)
//	expected := []string{"1", "2", "3"}
//
//	assert.Equal(t, tokenIndexes, tokens[0].typ)
//	assert.Equal(t, expected, tokens[0].v)
//}
//func Test_Indexes_word(t *testing.T) {
//	lexer, _ := newLexer(`$[test,second   ,  3]`)
//	tokens, err := lexer.parseExpressionTokens()
//	assert.Equal(t, nil, err)
//	expected := []string{"test", "second", "3"}
//
//	assert.Equal(t, tokenIndexes, tokens[0].typ)
//	assert.Equal(t, expected, tokens[0].v)
//}
