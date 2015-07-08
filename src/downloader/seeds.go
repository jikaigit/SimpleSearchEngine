package downloader

import (
	"io/ioutil"
	"logger"
	"os"
	"strings"
)

const seed_file string = "../../config/seeds.txt"

// 从爬行种子文件获取种子，下载器从这些种子开始下载
func GetSeeds() (uris []string) {
	data, err := ioutil.ReadFile(seed_file)
	if err != nil {
		logger.Panic("获取爬行种子失败")
	}
	uris = strings.Split(string(data), "\n")
	return uris
}

func AddSeed(uri string) (err error) {
	var (
		file *os.File
	)
	if file, err = os.OpenFile(seed_file, os.O_APPEND, os.ModePerm); err != nil {
		logger.Log("打开爬行种子文件失败")
		return err
	}
	defer func() {
		if err = file.Close(); err != nil {
			logger.Log("关闭爬行种子文件失败")
		}
	}()
	if _, err = file.WriteString(uri + "\r\n"); err != nil {
		logger.Log("向爬行种子文件中添加新种子失败")
		return err
	}
	return nil
}
