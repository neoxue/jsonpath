package jsonpath

import (
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
)

/*
func TestNewJsonPath(t *testing.T) {
	data := map[string]interface{}{"a": "test"}
	a := &JsonPath{Data: data}
	result := a.Find("$.a")
	respect := &Result{collection: []interface{}{"test"}, err: nil}
	assert.True(t, len(result.collection) == len(respect.collection))
	assert.True(t, reflect.TypeOf(result.collection[0]) == reflect.TypeOf(respect.collection[0]))
	assert.True(t, result.collection[0].(string) == respect.collection[0].(string))
	assert.True(t, result.collection[0] == respect.collection[0])
}

func TestJPGetIndex(t *testing.T) {
	data := map[string]interface{}{"a": "test"}
	a := &JsonPath{Data: data}
	result := a.Find("$.a")
	respect := &Result{collection: []interface{}{"test"}, err: nil}
	assert.True(t, len(result.collection) == len(respect.collection))
	assert.True(t, reflect.TypeOf(result.collection[0]) == reflect.TypeOf(respect.collection[0]))
	assert.True(t, result.collection[0].(string) == respect.collection[0].(string))
	assert.True(t, result.collection[0] == respect.collection[0])
}
func TestJPGetMapIndexes(t *testing.T) {
	data := map[string]interface{}{"a": "test", "b": "testb"}
	a := &JsonPath{Data: data}
	result := a.Find("$[a,b]")
	respect := &Result{collection: []interface{}{"test", "testb"}, err: nil}
	assert.True(t, len(result.collection) == len(respect.collection))
	assert.True(t, reflect.TypeOf(result.collection[0]) == reflect.TypeOf(respect.collection[0]))
	assert.True(t, result.collection[0].(string) == respect.collection[0].(string))
	assert.True(t, result.collection[0] == respect.collection[0])
}
func TestJPGetArraySlice(t *testing.T) {
	data := []interface{}{"test", "test1", "test2", "test3", "test4", "test5", "test6"}
	a := &JsonPath{Data: data}
	result := a.Find("$[2:5:2]")
	respect := &Result{collection: []interface{}{"test2", "test4"}, err: nil}
	assert.True(t, len(result.collection) == len(respect.collection))
	assert.True(t, reflect.TypeOf(result.collection[0]) == reflect.TypeOf(respect.collection[0]))
	assert.True(t, result.collection[0].(string) == respect.collection[0].(string))
	assert.True(t, result.collection[0] == respect.collection[0])
}

func TestJPGetArrayFilterExpressionGtEqual(t *testing.T) {
	data := []interface{}{"test", "test1", "test2", "test3", "test4", "test5", "test6"}
	a := &JsonPath{Data: data}
	result := a.Find("$[?(@ > \"test4\")]")
	respect := &Result{collection: []interface{}{"test5", "test6"}, err: nil}
	assert.True(t, len(result.collection) == len(respect.collection))
	assert.True(t, reflect.TypeOf(result.collection[0]) == reflect.TypeOf(respect.collection[0]))
	assert.True(t, result.collection[0].(string) == respect.collection[0].(string))
	assert.True(t, result.collection[0] == respect.collection[0])
}


func TestJPSetArraySetSimple(t *testing.T) {
	data := map[string]interface{}{}
	respect := map[string]interface{}{}
	json.Unmarshal([]byte("{\"a\":[\"test1\",\"test2\",\"test3\",\"test4\",\"test5\",\"test6\"]}"), &data)
	a := &JsonPath{Data: data}
	a.SetValue("$.a.1", "set1")
	json.Unmarshal([]byte("{\"a\":[\"test1\",\"set1\",\"test3\",\"test4\",\"test5\",\"test6\"]}"), &respect)
	assert.True(t, respect["a"].([]interface{})[1] == data["a"].([]interface{})[1])
}
func TestJPSetArrayUnSetSimple(t *testing.T) {
	data := map[string]interface{}{}
	respect := map[string]interface{}{}
	json.Unmarshal([]byte("{\"a\":[\"test1\",\"test2\",\"test3\",\"test4\",\"test5\",\"test6\"]}"), &data)
	a := &JsonPath{Data: data}
	a.UnsetValue("$.a.1")
	json.Unmarshal([]byte("{\"a\":[\"test1\",\"test3\",\"test4\",\"test5\",\"test6\"]}"), &respect)
	assert.True(t, respect["a"].([]interface{})[1] == data["a"].([]interface{})[1])
}

func TestJPGetArrayFilterScript(t *testing.T) {
	data := []interface{}{"test", "test1", "test2", "test3", "test4", "test5", "test6"}
	a := &JsonPath{Data: data}
	result := a.Find("$[(@.length - 1)]")
	assert.True(t, result.err.Error() == "jsonpath lexer: do not support query script now: {@.length - 1}")
}
*/
func TestJPGetArrayFilterExpressionGtEqualbeta(t *testing.T) {
	data := map[string]interface{}{}
	json.Unmarshal([]byte("{\"a\":[\"test1\",\"test2\",\"test3\",\"test4\",\"test5\",\"test6\"], \"b\":4}"), &data)
	a := &JsonPath{Data: data}
	result := a.Find("$.a[?(@ =~ $.b)]")
	respect := &Result{collection: []interface{}{"test4"}, err: nil}
	fmt.Println(result.err)
	assert.True(t, len(result.collection) == len(respect.collection))
	assert.True(t, reflect.TypeOf(result.collection[0]) == reflect.TypeOf(respect.collection[0]))
	assert.True(t, result.collection[0].(string) == respect.collection[0].(string))
	assert.True(t, result.collection[0] == respect.collection[0])
}
