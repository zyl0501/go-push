package main

import "fmt"

func main() {
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
