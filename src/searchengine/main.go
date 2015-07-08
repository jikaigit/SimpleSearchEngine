package main

import (
	"fmt"
)

func main() {
	var m map[int]string = make(map[int]string)
	for i, data := range m {
		fmt.Println(i, data)
	}
}
