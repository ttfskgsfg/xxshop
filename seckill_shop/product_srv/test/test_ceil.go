package main

import (
	"fmt"
)

func main() {

	count := 25

	page_size := 8


	ret := (count + page_size -1) / page_size

	fmt.Println(ret)


	//fmt.Println(count/page_size)
	//
	//fmt.Println(math.Ceil(1.5))


}
