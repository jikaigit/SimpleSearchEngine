package indexer

import (
	"github.com/huichen/sego"
	"os"
	"strconv"
)

const (
	dic_file string = "../github.com/huichen/sego/data/dictionary.txt"
	idx_dir  string = "../../index"
)

var (
	index_file_count    int      = 0
	index_seg_file_chan chan int = make(chan int, 50)
)

// 索引器用来分析下载器传来的页面数据，然后建立索引以供搜索引擎
// 执行查询操作
type Indexer struct {
	index_file *os.File
	segmenter  sego.Segmenter
}

func (this *Indexer) Init() {
	this.segmenter.LoadDictionary(dic_file)
}

func (this *Indexer) AnalyseAndGenerateIndex(contents []string, source string) {
	var index_cache IndexCache

	index_cache.Init(source)
	for _, content := range contents {
		segs := this.segmenter.Segment([]byte(content))
		for _, seg := range segs {
			word := seg.Token().Text()
			switch word {
			case " ", "\t", "\r", "\n", ".", ",", "\"", "'", "?", "!":
				// do nothing
			default:
				index_cache.Add(word)
			}
		}
	}

	index_file_count++
	index_seg_file_chan <- index_file_count
	index_cache.WriteToFile(idx_dir + "/" + strconv.Itoa(index_file_count))
}
