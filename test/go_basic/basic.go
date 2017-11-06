package main

import "fmt"

const (
	Sunday    = 1 + iota //0
	Monday               //1
	Tuesday              //2
	Wednesday            //3
	Thursday             //4
	Friday               //5
	Saturday             //6
	Unknown   = -1
	Unknown2   //-1
)

func main() {
	//mapTest()
	enumTest()
}
func mapTest() {
	var pc map[string]string
	pc = make(map[string]string)
	pc["qingdao"] = "青岛"
	pc["jinan"] = "济南"
	pc["yantai"] = "烟台"
	fmt.Println("size:", len(pc))
	for k := range pc {
		fmt.Println("k:", k)
		delete(pc, k)
	}
	qingdao, ok := pc["qingdao"]
	if ok {
		fmt.Println(qingdao)
	} else {
		fmt.Println("元素不存在")
	}
	fmt.Println("size:", len(pc))
}

func enumTest() {
	fmt.Println("Sunday=", Sunday)
	fmt.Println("Monday=", Monday)
	fmt.Println("Tuesday=", Tuesday)
	fmt.Println("Wedenesday=", Wednesday)
	fmt.Println("Thursday=", Thursday)
	fmt.Println("Friday=", Friday)
	fmt.Println("Saturday=", Saturday)
	fmt.Println("Unknown=", Unknown)
	fmt.Println("Unknown2=", Unknown2)
}
