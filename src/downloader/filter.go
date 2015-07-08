package downloader

import (
	"errors"
	"logger"
	"net/url"
)

// 用于记录被下载过的页面以防止重复下载相同的页面
type DownloadFilter struct {
	domain     string
	footprintf map[string]bool
}

func (this *DownloadFilter) Init(domain string) {
	this.domain = domain
	this.footprintf = make(map[string]bool)
}

// true  -> 页面已经被下载过
// false -> 页面没有被下载过
func (this *DownloadFilter) IsRepeat(uri string) bool {
	if _, ok := this.footprintf[uri]; ok {
		return true
	}
	return false
}

// 记录指定下载完的页面，用来防止以后再次下载它
func (this *DownloadFilter) AddFootPrint(uri string) error {
	u, err := url.Parse(uri)
	if err != nil {
		logger.Log("无法被解析的资源定位符:" + uri)
		return err
	}
	if u.Host != this.domain {
		logger.Log("添加的资源定位与过滤器中设置的指定域不同, 主域:" + this.domain + ", 添加域:" + u.Host)
		return errors.New("different domain")
	}
	this.footprintf[uri] = true
	return nil
}
