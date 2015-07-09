package main

import (
	"io/ioutil"
	"logger"
	"net/http"
)

func MainPage(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		w.WriteHeader(200)
		w.Header().Set("Content-Type", "text/html")
		pagedata, err := ioutil.ReadFile("../frontend/index.html")
		if err != nil {
			logger.Log("读取页面文件失败")
			return
		}
		w.Write(pagedata)
	}
}
