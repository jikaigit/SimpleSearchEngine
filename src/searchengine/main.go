package main

import (
	"fmt"
	"indexer"
	"io/ioutil"
	"parser"
	"unicode/utf8"
)

func main() {
	data, err := ioutil.ReadFile("C:\\Users\\Administrator\\Desktop\\input.txt")
	if err != nil {
		fmt.Println("打开文件失败")
		return
	}

	contents, _ := parser.ParsePage(data, "host.com")

	var cacher indexer.Indexer
	cacher.Init()
	cacher.AnalyseAndGenerateIndex(contents, "host.com")
}
