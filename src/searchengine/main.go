package main

import (
	"downloader"
)

func main() {
	var filter downloader.DownloadFilter
	filter.Init("www.baidu.com")
	if err := filter.AddFootPrint("http://www.baidu.com/"); err != nil {
		println(err.Error())
	}
	if err := filter.AddFootPrint("www.baidu.com"); err != nil {
		println(err.Error())
	}
}
