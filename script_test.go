package jsonpath

import (
	"testing"

	"encoding/json"

	"fmt"

	"github.com/stretchr/testify/assert"
)

func Test_a(t *testing.T) {
	lexer, _ := newLexer(`$[1,2   ,  3]`)
	tokens, err := lexer.parseExpressionTokens()
	assert.Equal(t, nil, err)
	expected := "1,2   ,  3"
	assert.Equal(t, expected, tokens[0].v)
	assert.Equal(t, "indexes", tokens[0].typ)

	var stra = `["aaa","bbb","ccc","ddd"]`
	var myjsona interface{}
	json.Unmarshal([]byte(stra), &myjsona)
	jpa := JsonPath{Data: myjsona}
	resulta := jpa.Find(`$[1,2   ,  3]`)

	assert.Equal(t, "bbb", resulta.collection[0].(string))
	fmt.Println(resulta, err)
	fmt.Println(resulta.collection)

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
	jp := JsonPath{Data: myjson}
	result, finded, err := jp.Find("$.c[e][3][d]").First()
	fmt.Println(result, finded, err)
	result = nil
	fmt.Println(myjson)
}

func testint() int {
	return 0
}
