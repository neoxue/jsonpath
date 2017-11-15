package jsonpath

import (
	"testing"

	"encoding/json"

	"fmt"

	"github.com/stretchr/testify/assert"
)

func Test_a(t *testing.T) {
	lexer, _ := newLexer(`$[1,2   ,  3]`)
	tokens, err := lexer.ParseExpressionTokens()
	assert.Equal(t, nil, err)
	expected := []string{"1", "2", "3"}

	assert.Equal(t, expected, tokens[0].v)

	str := `{
		"a":"b",
		"c":{
			"e":[
				1,
				4,
				"rrr",
				{"d":"f"}
			]
		}
	}
	`
	var myjson interface{}
	json.Unmarshal([]byte(str), &myjson)
	jp := JsonPath{data: myjson}
	result, _ := jp.Find("$.c[e][3][d]")
	fmt.Println(result)
}

func testint() int {
	return 0
}
