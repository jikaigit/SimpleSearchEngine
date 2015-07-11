package main

import (
	"indexer"
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

func Search(w http.ResponseWriter, r *http.Request) {
	querys := r.URL.Query()
	question := querys.Get("question")
	sources := indexer.SearchEngineIndexer.Search(question)
	var page_data string = `
    <html>
    <head>
        <title>小凯搜索引擎-结果页面</title>
        <style>
            body {
                background-color: #E8E8E8;
            }
            .result {
                font-family: "Microsoft YaHei", SimHei;
                width: 88%;
                margin: 15px 0px 0px 6%;
                height: 50px;
                line-height: 50px;
                color: #555555;
                overflow: hidden;
                padding: 0px 15px 0px 15px;
                box-shadow: 0px 0px 5px 0px #AAAAAA;
                font-size: 18px;
                background-color: #FFFFFF;
                cursor: pointer;
            }
        </style>
    </head>
    <body>`
	if sources != nil {
		for _, source := range sources {
			page_data += "<div class='result'><a href='"
			page_data += source
			page_data += "'>"
			page_data += source
			page_data += "</a></div>"
		}
	}
	page_data += "</body></html>"

	w.WriteHeader(200)
	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte(page_data))
}
