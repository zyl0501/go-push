package _map

import "fmt"

type A struct {
	A int
	B string
	C []byte
}

func MapTest(){
	m :=make(map[string]A)
	m["1"] = A{A:2}
	v,ok:=m["2"]
	fmt.Printf("%v, %b",v,ok)
}