package main

import (
	"downloader"
	"indexer"
	"logger"
	"net/http"
)

func init() {
	var dl downloader.Downloader
	dl.Init(12, &indexer.SearchEngineIndexer)
	go dl.StartDownload()
}

func main() {
	http.HandleFunc("/", MainPage)
	http.HandleFunc("/search", Search)

	if err := http.ListenAndServe(":80", nil); err != nil {
		logger.Log("服务器启动失败")
	}

}
