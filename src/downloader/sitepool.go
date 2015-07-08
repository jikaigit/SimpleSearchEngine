package downloader

import (
	"errors"
)

// 站点池用来存储下载器从页面中解析出的不同域的域名，用来并行
// 下载其他站点的页面
type SitePool struct {
	sites     map[string]bool
	footprint map[string]bool
}

func (this *SitePool) Init() {
	this.sites = make(map[string]bool)
	this.footprint = make(map[string]bool)
}

// 向站点池中增加一个域名，自动去重哦
func (this *SitePool) AddSite(domain string) {
	if _, ok := this.footprint[domain]; ok {
		return
	}
	if _, ok := this.sites[domain]; ok {
		return
	}
	this.sites[domain] = true
}

func (this *SitePool) GetSite() (uri string, err error) {
	for domain, _ := range this.sites {
		this.footprint[domain] = true
		delete(this.sites, domain)
		return "http://" + domain + "/", nil
	}
	return "", errors.New("站点池已空")
}

func (this *SitePool) Close() {
	this.sites = nil
}
