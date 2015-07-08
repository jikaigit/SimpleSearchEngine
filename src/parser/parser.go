package parser

import (
	"downloader"
)

func ParsePage(pagedata string, rooturl string) (contents []string, suburls []string) {
	var (
		err     error
		length  int = len(pagedata)
		curtag  string
		content string
		suburl  string
		quo     byte
	)
	for i := 0; i < length; i++ {
		if pagedata[i] == '<' {
			// 跳过闭合标签
			if i+1 >= length || pagedata[i] == '/' {
				i += 2
				continue
			}
			// 如果有要解析的标签名，前跳过空白符
			for pagedata[i] == ' ' || pagedata[i] == '\t' || pagedata[i] == '\r' || pagedata[i] == '\n' {
				i++
			}
			curtag = ""
			for pagedata[i] != ' ' && pagedata[i] != '\t' && pagedata[i] != '\r' && pagedata[i] != '\n' {
				curtag += string(pagedata[i])
				i++
			}
		}

		if pagedata[i] == '>' && curtag != "style" && curtag != "script" {
			// 单标签，例如<.../>这样的但标签是没有内容的
			if i-1 < 0 || pagedata[i-1] == '/' {
				continue
			}
			// 跳过没用的空白符
			for {
				if i++; i >= length {
					return contents, suburls
				}
				if pagedata[i] != ' ' && pagedata[i] != '\t' && pagedata[i] != '\r' && pagedata[i] != '\n' {
					break
				}
			}
			if pagedata[i] == '<' {
				continue
			}
			content = ""
			for pagedata[i] != '<' {
				content += string(pagedata[i])
				i++
			}
			contents = append(contents, content)
		}

		if pagedata[i] == 'h' || pagedata[i] == 'H' {
			if i+1 >= length || (pagedata[i+1] != 'r' && pagedata[i+1] != 'R') {
				continue
			}
			i++
			if i+1 >= length || (pagedata[i+1] != 'e' && pagedata[i+1] != 'E') {
				continue
			}
			i++
			if i+1 >= length || (pagedata[i+1] != 'f' && pagedata[i+1] != 'F') {
				continue
			}
			i++
			for pagedata[i] == ' ' || pagedata[i] == '\t' || pagedata[i] == '\r' || pagedata[i] == '\n' {
				i++
			}
			if pagedata[i] != '=' {
				continue
			}
			i++
			for pagedata[i] == ' ' || pagedata[i] == '\t' || pagedata[i] == '\r' || pagedata[i] == '\n' {
				i++
			}
			if pagedata[i] != '"' && pagedata[i] != '\'' {
				continue
			}
			quo = pagedata[i]
			i++
			suburl = ""
			for pagedata[i] != quo {
				suburl += string(pagedata[i])
				i++
			}
			if suburl, err = downloader.UrlNormalize(rooturl, suburl); err == nil {
				suburls = append(suburls, suburl)
			}
		}
	}
	return contents, suburls
}
