package downloader

import (
	"bytes"
	"io"
	"logger"
	"net/http"
)

// 下载单个Internet资源
func Download(uri string) (data []byte, err error) {
	var (
		res  *http.Response
		buff bytes.Buffer
	)
	if res, err = http.Get(uri); err != nil {
		// logger.Log("无法下载资源:" + uri)
		return nil, err
	}
	defer func() {
		if err = res.Body.Close(); err != nil {
			logger.Log("下载资源:" + uri + "的HTTP响应关闭失败")
		}
	}()
	if _, err = io.Copy(&buff, res.Body); err != nil {
		logger.Log("将资源数据写入缓冲区失败")
		return nil, err
	}
	return buff.Bytes(), nil
}
