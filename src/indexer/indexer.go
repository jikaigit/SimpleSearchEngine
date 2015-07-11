package indexer

import (
	"github.com/huichen/sego"
	"sync"
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
	index     IndexCache
	segmenter sego.Segmenter
	lock      sync.RWMutex
}

func (this *Indexer) Init() {
	this.segmenter.LoadDictionary(dic_file)
}

func (this *Indexer) AnalyseAndGenerateIndex(contents []string, source string) {
	if contents == nil || len(contents) <= 0 || source == "" {
		return
	}
	for _, content := range contents {
		segs := this.segmenter.Segment([]byte(content))
		if segs == nil || len(segs) <= 0 {
			break
		}
		for _, seg := range segs {
			if seg.Token() == nil {
				continue
			}
			word := seg.Token().Text()
			switch word {
			case "", " ", "\t", "\r", "\n", ".", ",", "\"", "'", "?", "!":
				// do nothing
			default:
				this.lock.Lock()
				this.index.Add(word, source)
				this.lock.Unlock()
			}
		}
	}
}

func (this Indexer) Search(question string) (sources map[string]int) {
	this.lock.Lock()
	node := this.index.Search(question)
	this.lock.Unlock()
	if node != nil && node.sources != nil {
		return node.sources
	}
	if node == nil {
		return nil
	}
	return nil
}

func (this Indexer) Debug() {
	this.lock.RLock()
	this.index.InOrderTravel()
	this.lock.RUnlock()
}
