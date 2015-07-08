package downloader

import (
	"fmt"
	"time"
)

func StartDownload(maxsite int, deepth int) {
	var (
		seeds []string = GetSeeds()
	)
	for {
		for _, seed := range seeds {
			fmt.Println(seed)
		}

		// 隔一段时间才重新从种子连接开始爬行
		time.Sleep(time.Duration(time.Hour * 1))
	}
}
