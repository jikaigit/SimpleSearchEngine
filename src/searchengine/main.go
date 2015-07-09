package main

import (
	"downloader"
)

func main() {
	var dl downloader.Downloader
	dl.Init(12)
	dl.StartDownload()
}
