package main

import (
	"downloader"
)

func main() {
	var dloader downloader.Downloader
	dloader.Init(3)
	dloader.StartDownload()
}
