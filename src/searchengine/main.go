package main

import (
	"downloader"
	"logger"
	"net/http"
)

func init() {
	var dl downloader.Downloader
	dl.Init(12)
	go dl.StartDownload()
}

func main() {
	http.HandleFunc("/", MainPage)

	if err := http.ListenAndServe(":80", nil); err != nil {
		logger.Log("服务器启动失败")
	}
}
