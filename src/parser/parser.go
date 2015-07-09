package parser

import (
	"errors"
	. "net/url"
	"strconv"
	"strings"
)

// 解析下载来的页面，将页面中的文章和链接给分解出来
func ParsePage(pagedata []byte, rooturl string) (contents []string, suburls []string) {
	var (
		length  int = len(pagedata)
		curtag  string
		content string
		suburl  string
		quo     byte
	)
	for i := 0; i < length; {
		if pagedata[i] == '<' {
			// 跳过!大头的标签
			if i+1 < length && pagedata[i+1] == '!' {
				if i++; i >= length {
					return contents, suburls
				}
				for {
					if i++; i >= length {
						return contents, suburls
					}
					if pagedata[i] == '>' {
						break
					}
				}
				if i++; i >= length {
					return contents, suburls
				}
			}
			// 跳过闭合标签
			if i+1 >= length || pagedata[i+1] == '/' {
				i += 2
				continue
			}
			if i++; i >= length {
				return contents, suburls
			}
			// 如果有要解析的标签名，前跳过空白符
			for pagedata[i] == ' ' || pagedata[i] == '\t' || pagedata[i] == '\r' || pagedata[i] == '\n' {
				if i++; i >= length {
					return contents, suburls
				}
			}
			curtag = ""
			for ('a' <= pagedata[i] && pagedata[i] <= 'z') || ('A' <= pagedata[i] && pagedata[i] <= 'Z') {
				curtag += string(pagedata[i])
				if i++; i >= length {
					return contents, suburls
				}
			}
			continue
		}

		if pagedata[i] == '>' && curtag != "style" && curtag != "script" {
			// 单标签，例如<.../>这样的但标签是没有内容的
			if i-1 < 0 || pagedata[i-1] == '/' {
				if i++; i >= length {
					return contents, suburls
				}
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
				if i++; i >= length {
					return contents, suburls
				}
			}
			contents = append(contents, content)
		}

		if (pagedata[i] == 'h' || pagedata[i] == 'H') && curtag != "link" {
			if i-1 < 0 || (pagedata[i-1] != ' ' && pagedata[i-1] != '\t' && pagedata[i-1] != '\r' && pagedata[i-1] != '\n') {
				if i++; i >= length {
					return contents, suburls
				}
				continue
			}
			if i+1 >= length || (pagedata[i+1] != 'r' && pagedata[i+1] != 'R') {
				if i++; i >= length {
					return contents, suburls
				}
				continue
			}
			if i++; i >= length {
				return contents, suburls
			}
			if i+1 >= length || (pagedata[i+1] != 'e' && pagedata[i+1] != 'E') {
				if i++; i >= length {
					return contents, suburls
				}
				continue
			}
			if i++; i >= length {
				return contents, suburls
			}
			if i+1 >= length || (pagedata[i+1] != 'f' && pagedata[i+1] != 'F') {
				if i++; i >= length {
					return contents, suburls
				}
				continue
			}
			if i++; i >= length {
				return contents, suburls
			}
			for {
				if i++; i+1 >= length {
					return contents, suburls
				}
				if pagedata[i] != ' ' && pagedata[i] != '\t' && pagedata[i] != '\r' && pagedata[i] != '\n' {
					break
				}
			}
			if pagedata[i] != '=' {
				continue
			}
			if i++; i >= length {
				return contents, suburls
			}
			for pagedata[i] == ' ' || pagedata[i] == '\t' || pagedata[i] == '\r' || pagedata[i] == '\n' {
				if i++; i >= length {
					return contents, suburls
				}
			}
			if pagedata[i] != '"' && pagedata[i] != '\'' {
				continue
			}
			quo = pagedata[i]
			if i++; i >= length {
				return contents, suburls
			}
			suburl = ""
			for pagedata[i] != quo {
				suburl += string(pagedata[i])
				if i++; i >= length {
					return contents, suburls
				}
			}
			if suburl, _ = UrlNormalize(rooturl, suburl); suburl != "" {
				suburls = append(suburls, suburl)
			}
		}

		if i++; i >= length {
			return contents, suburls
		}
	}
	return contents, suburls
}

// 自动解析页面中的URL并且自动将相对URL转换成绝对URL
func ParseURL(pagecontent string, rooturl string) (urls []string) {
	var (
		url        string
		length     int = len(pagecontent)
		match_flag int
		quotation  byte
	)

	for i := 0; i < length; {
		// 匹配关键字'href'或者'src'
		if pagecontent[i] == 'h' || pagecontent[i] == 'H' {
			for match_flag = 0; match_flag < 3; match_flag++ {
				if i+1 >= length {
					return urls
				}
				i++
				if match_flag == 0 && (pagecontent[i] != 'r' && pagecontent[i] != 'R') {
					break
				} else if match_flag == 1 && (pagecontent[i] != 'e' && pagecontent[i] != 'E') {
					break
				} else if match_flag == 2 && (pagecontent[i] != 'f' && pagecontent[i] != 'F') {
					break
				}
			}
			if match_flag != 3 {
				// 没有匹配则进入下一次循环
				continue
			}

		} else if pagecontent[i] == 's' || pagecontent[i] == 'S' {
			for match_flag = 0; match_flag < 2; match_flag++ {
				if i+1 >= length {
					return urls
				}
				i++
				if match_flag == 0 && (pagecontent[i] != 'r' && pagecontent[i] != 'R') {
					break
				} else if match_flag == 1 && (pagecontent[i] != 'c' && pagecontent[i] != 'C') {
					break
				}
			}
			if match_flag != 2 {
				// 如果不匹配，进入下一次i循环
				continue
			}

		} else {
			i++
			continue
		}

		// 跳过空格和'='
		for {
			if i+1 >= length {
				return urls
			}
			i++
			if pagecontent[i] != ' ' {
				break
			}
		}

		if pagecontent[i] != '=' {
			continue
		}

		for {
			if i+1 >= length {
				return urls
			}
			i++
			if pagecontent[i] != ' ' {
				break
			}
		}

		// 记录引号的类型以方便我们之后进行闭合操作
		if pagecontent[i] == '"' || pagecontent[i] == '\'' {
			quotation = pagecontent[i]
		} else {
			continue
		}

		// 开始解析URL
		for {
			if i+1 >= length {
				return urls
			}
			i++

			if pagecontent[i] == quotation {
				if url, _ = UrlNormalize(rooturl, url); url != "" {
					urls = append(urls, url)
				}
				url = ""
				break
			}
			url += string(pagecontent[i])
		}
	}

	return urls
}

// 将URL编程更加标准的形式，如果在转换过程中发生错误，那么将会返回空字符串""
func UrlNormalize(rooturl string, relativeurl string) (absoluteurl string, err error) {
	var (
		rootu   *URL
		rltvu   *URL
		tempurl string
	)

	if rootu, err = Parse(rooturl); err != nil {
		return "", errors.New("parse/util.UrlNormalize: " + err.Error())
	}
	if rltvu, err = Parse(relativeurl); err != nil {
		return "", errors.New("parse/util.UrlNormalize: " + err.Error())
	}

	// 尝试着将根路径和相对路径进行连接
	// 如果相对路径是一个绝对路径，就忽略跟URL
	// 如果相对路径是一个相对路径，那么根路径是必须的
	if rooturl == "" && relativeurl == "" {
		return "", errors.New("parse/util.UrlNormalize: url is empty.")
	} else if rltvu.Scheme != "" {
		tempurl = relativeurl
	} else if rltvu.Scheme == "" && (rooturl == "" || rootu.Scheme == "") {
		return "", errors.New("parse/util.UrlNormalize: root url required.")
	} else {
		tempurl = rooturl + "/" + relativeurl
	}

	if strings.Index(tempurl, "javascript:") != -1 {
		return "", errors.New("parse/util.UrlNormalize: find 'jsvascript:' snippet.")
	}

	var (
		length               int  = len(tempurl)
		lower                bool = true
		path_first_slash_pos int  = -1
		percent_code         string
		char                 int64
	)
	for i := 0; i < length; i++ {
		if tempurl[i] == '#' {
			// 即使有#也表示相同的URL，你懂得~
			break

		} else if i+2 < length && tempurl[i] == '%' {
			// 解码%开头的字符编码
			percent_code = string(tempurl[i+1]) + string(tempurl[i+2])
			char, _ = strconv.ParseInt(percent_code, 16, 64)

			if ('A' <= char && char <= 'Z') || ('a' <= char && char <= 'z') || ('0' <= char && char <= '9') || (char == '-' || char == '.' || char == '_' || char == '~') {
				absoluteurl += string(rune(char))
			} else {
				absoluteurl += "%"

				if 'a' <= tempurl[i+1] && tempurl[i+1] <= 'z' {
					absoluteurl += string(rune(tempurl[i+1] - 32))
				} else {
					absoluteurl += string(tempurl[i+1])
				}

				if 'a' <= tempurl[i+2] && tempurl[i+2] <= 'z' {
					absoluteurl += string(rune(tempurl[i+2] - 32))
				} else {
					absoluteurl += string(tempurl[i+2])
				}
			}

			if i+2 >= length {
				break
			}
			i += 2

		} else if path_first_slash_pos == -1 && i-2 >= 0 && (tempurl[i] == '/' && tempurl[i-1] != ':' && tempurl[i-2] != ':') {
			// 标记第一个斜杠的位置作为路径部分的开始，这个斜杠的位置也是../的
			// 边界和大小写转换的边界
			absoluteurl += "/"
			path_first_slash_pos = i
			lower = false

		} else if tempurl[i] == '/' {
			// 将双斜杠(//)转换成但斜杠(/)
			if i-1 >= 0 && tempurl[i-1] == '/' {
				// 协议部分的://中的双斜杠不要转换
				if i-2 >= 0 && tempurl[i-2] == ':' {
					absoluteurl += "/"
				}
			} else {
				absoluteurl += "/"
			}

		} else if i+2 < length && ((tempurl[i] == '/' && tempurl[i+1] == '.' && tempurl[i+2] == '.') || (tempurl[i] == '.' && tempurl[i+1] == '.' && tempurl[i+2] == '/')) {
			// 处理像./或/.或../或/..这些形式的相对路径
			j := len(absoluteurl) - 1
			if absoluteurl[j] == '/' {
				j--
			}
			for ; j >= path_first_slash_pos; j-- {
				if absoluteurl[j] == '/' {
					absoluteurl = absoluteurl[:j+1]
					break
				}
			}
			i += 2

		} else if i+1 < length && (tempurl[i] == '/' && tempurl[i+1] == '.') {
			i++

		} else if i+1 < length && (tempurl[i] == '.' && tempurl[i+1] == '/') {
			if i-1 >= 0 && tempurl[i-1] != '/' {
				absoluteurl += "."
			} else {
				i++
			}

		} else if lower && (tempurl[i] >= 'A' && tempurl[i] <= 'Z') {
			// 将协议和域名部分转换成小写形式
			absoluteurl += string(tempurl[i] + 32)

		} else if tempurl[i] == '?' && i == length-1 {
			// 如果查询参数是空的就去掉?
			break

		} else {
			absoluteurl += string(tempurl[i])
		}
	}

	// 如果是个绝对路径，在返回之前先测试它的格式是否正确
	if _, err = Parse(absoluteurl); err != nil {
		return "", err
	}

	return absoluteurl, nil
}

// 逐渐解析一个URL的路径
//
// 比如我们现在要解析URL:
// "http://www.host.com/path1/path2/index.html"
//
// 返回的结果会是像这样:
// "http://www.host.com/"
// "http://www.host.com/path1/"
// "http://www.host.com/path1/path2/"
func PathAscend(rawurl string) (urls []string, err error) {
	var (
		u       *URL
		length  int
		tempurl string
		subpath string
	)

	if u, err = Parse(rawurl); err != nil {
		return urls, err
	}

	if u.Path == "" || u.Path == "/" {
		tempurl = u.Scheme + "://" + u.Host + "/"
		urls = append(urls, tempurl)
		return urls, nil
	}

	length = len(u.Path)
	for i := 0; i < length; i++ {
		if u.Path[i] != '/' {
			subpath += string(u.Path[i])
		} else {
			if i-1 >= 0 && u.Path[i-1] == '/' {
				continue
			}
			subpath += string(u.Path[i])
			tempurl = u.Scheme + "://" + u.Host + subpath
			urls = append(urls, tempurl)
		}
	}

	return urls, nil
}
