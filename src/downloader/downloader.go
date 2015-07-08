package downloader

import (
	"fmt"
	"logger"
)

// 搜索引擎的下载器，用来游荡在互联网上下载信息
type Downloader struct {
	site_crawl_max_count_ctl_chan chan int
	sitepool                      SitePool
	site_need_crawl               int64
	site_finish_crawl             int64
}

// 参数:
// maxsite用来表示能同时爬行的网站数目
func (this *Downloader) Init(maxsite int) {
	this.site_crawl_max_count_ctl_chan = make(chan int, maxsite)
	this.sitepool.Init()
}

func (this *Downloader) StartDownload() {
	var (
		err    error
		seeds  []string = GetSeeds()
		domain string
	)

	for _, seed := range seeds {
		this.sitepool.AddSite(seed)
	}

	for {
		if domain, err = this.sitepool.GetSite(); err != nil && this.site_need_crawl == this.site_finish_crawl {
			logger.Log("下载已经全部完成")
			return
		}
		this.site_crawl_max_count_ctl_chan <- 1
		this.site_need_crawl++
		go this.travelSiteAndDownload(domain, 5)
	}
}

func (this *Downloader) travelSiteAndDownload(domain string, maxdeepth int) {
	defer func() {
		<-this.site_crawl_max_count_ctl_chan
		this.site_finish_crawl++
	}()
	var (
		filter      DownloadFilter
		max_routine chan int = make(chan int, 10)
	)
	filter.Init(domain)

	max_routine <- 1
	this.download("http://"+domain+"/", &filter, max_routine, maxdeepth)
}

func (this *Downloader) download(uri string, filter *DownloadFilter, max_routine chan int, deepth int) {
	if filter.IsRepeat(uri) {
		<-max_routine
		return
	}

	var (
		data []byte
		err  error
	)
	if data, err = Download(uri); err != nil {
		<-max_routine
		return
	} else {
		filter.AddFootPrint(uri)
		<-max_routine
	}

	fmt.Println(string(uri))

	// 这里往下就不要再向max_routine里发送信息了

	deepth--
	if deepth <= 0 {
		return
	}

	suburls := ParseURL(string(data), uri)
	for _, suburl := range suburls {
		max_routine <- 1
		go this.download(suburl, filter, max_routine, deepth)
	}
}
