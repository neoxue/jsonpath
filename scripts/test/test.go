package main

import (
	"container/ring"
	"fmt"
	"reflect"
	"strconv"
)

func main() {

	fmt.Println("abc" > "ac")
	fmt.Println("abc" < "ac")

	str1 := "20.99"
	test1, err := strconv.Atoi(str1)
	test2, err2 := strconv.ParseFloat(str1, 64)
	fmt.Println(err2)
	fmt.Println(test2)
	fmt.Println(err)
	fmt.Println(test1)

	str := "test'"
	for _, c := range str {
		fmt.Println(reflect.TypeOf(c))
		fmt.Println(c == '\'')
		fmt.Println(c)
	}
	fmt.Println(str[4] == '\'')
	fmt.Println(reflect.TypeOf(str[1]))

	//RingFunc()

}
func RingFunc() {
	r := ring.New(10) //初始长度10
	for i := 0; i < r.Len(); i++ {
		r.Value = i
		r = r.Next()
	}
	for i := 0; i < r.Len(); i++ {
		fmt.Println(r.Value)
		r = r.Next()
	}
	r = r.Move(6)
	fmt.Println(r.Value) //6
	r1 := r.Unlink(19)   //移除19%10=9个元素
	for i := 0; i < r1.Len(); i++ {
		fmt.Println(r1.Value)
		r1 = r1.Next()
	}
	fmt.Println(r.Len())  //10-9=1
	fmt.Println(r1.Len()) //9
}
