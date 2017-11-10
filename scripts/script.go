package main

import (
	"fmt"

	"github.com/neoxue/jsonpath"
)

func main() {
	lexer, error := jsonpath.NewLexer("$.a")
	fmt.Println(error)
	fmt.Println(lexer)

}
