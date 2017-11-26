package fun

import "fmt"

func FunPtrTest(){
	conn := ConnA{}
	test(&conn)
	fmt.Println(conn)
}

func test(c Conn){
	c.Set(9)
}

type ConnA struct {
	I int
}
func (conn ConnA) Set(i int){
	conn.I = i
}

type Conn interface {
	Set(i int)
}