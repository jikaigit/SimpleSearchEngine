package indexer

import (
	"github.com/huichen/sego"
	"os"
)

var SearchEngineIndexer Indexer

func init() {
	SearchEngineIndexer.Init()
}

const (
	dic_file string = "../github.com/huichen/sego/data/dictionary.txt"
	idx_dir  string = "../../index"
)

var (
	index_file_count int = 0
)

// 索引器用来分析下载器传来的页面数据，然后建立索引以供搜索引擎
// 执行查询操作
type Indexer struct {
	index_file *os.File
	index      IndexCache
	segmenter  sego.Segmenter
}

func (this *Indexer) Init() {
	this.segmenter.LoadDictionary(dic_file)
}

func (this *Indexer) AnalyseAndGenerateIndex(contents []string, source string) {
	var index_cache IndexCache

	for _, content := range contents {
		segs := this.segmenter.Segment([]byte(content))
		for _, seg := range segs {
			word := seg.Token().Text()
			switch word {
			case "", " ", "\t", "\r", "\n", ".", ",", "\"", "'", "?", "!":
				// do nothing
			default:
				index_cache.Add(word, source)
			}
		}
	}
}

func (this Indexer) Search(question string) (sources []string) {
	node := this.index.Search(question)
	if node != nil && node.sources != nil {
		return node.sources
	}
	return nil
}

func (this Indexer) Debug() {
	this.index.InOrderTravel()
}
