package downloader

type DownloadFilter struct {
	domain     string
	footprintf map[string]bool
}

func (filter *DownloadFilter) Init(domain string) {
	filter.domain = domain
	filter.footprintf = make(map[string]bool)
}

func (filter *DownloadFilter) AddFootPrint(uri string) {

}
