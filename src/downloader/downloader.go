package downloader

import (
	"logger"
)

// 搜索引擎的下载器，用来游荡在互联网上下载信息
type Downloader struct {
	site_cral_max_count_ctl_chan chan int
	sitepool                     SitePool
}

// 参数:
// maxsite用来表示能同时爬行的网站数目
func (this *Downloader) Init(maxsite int) {
	this.site_cral_max_count_ctl_chan = make(chan int, maxsite)
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
		if domain, err = this.sitepool.GetSite(); err != nil {
			logger.Log("下载已经全部完成")
			return
		}
		this.site_cral_max_count_ctl_chan <- 1
		go this.travelSiteAndDownload(domain, 5)
	}
}

func (this *Downloader) travelSiteAndDownload(domain string, maxdeepth int) {
	defer func() {
		<-this.site_cral_max_count_ctl_chan
	}()
	var (
		filter      DownloadFilter
		max_routine chan int = make(chan int, 5)
	)
	filter.Init(domain)

	this.download("http://"+domain+"/", &filter, max_routine)
}

func (this *Downloader) download(uri string, filter *DownloadFilter, max_routine chan int) {

}
