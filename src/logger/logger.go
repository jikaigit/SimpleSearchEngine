// 日志系统的目录组织形式：
// log
//  |year1
//     |month1
//         |day1.txt
//         |day2.txt
//         |day3.txt
//     |month2
//         |day1.txt
//         |day2.txt
//  |year2
//     |month1
//         |day1.txt
//         |day2.txt
//  ...
//
// 通过调用logger.Log函数来记录普通日志
// 通过调用logger.Panic函数来记录高危日志并终止服务器运行
//
package logger

import (
	"fmt"
	"os"
	"strconv"
	"sync"
	"time"
)

var (
	log_dir       string = "../../log"
	the_day       time.Time
	log_file_path string
	lock          sync.RWMutex
)

// Tip: 会在日志记录的末尾自动添加换行符
func Log(msg string) {
	var (
		err       error
		now       time.Time = time.Now()
		log_file  *os.File
		log_entry string
	)

	lock.Lock()
	defer lock.Unlock()

	// 每天都会创建一个新的日志文件
	if now.Day() != the_day.Day() {
		the_day = now
		createLogFile(the_day)
		log_file_path = getLogFilePath(the_day)
	}
	if log_file, err = os.OpenFile(log_file_path, os.O_APPEND, os.ModePerm); err != nil {
		panic("打开日志文件" + log_file_path + "时发生错误: " + err.Error())
	}
	defer func() {
		if err = log_file.Close(); err != nil {
			panic("关闭日志文件" + log_file_path + "时发生错误: " + err.Error())
		}
	}()

	log_entry = fmt.Sprintf("%04d-%02d-%02d %02d:%02d:%02d %s", now.Year(), now.Month(), now.Day(), now.Hour(), now.Minute(), now.Second(), msg)
	fmt.Println(log_entry)
	fmt.Fprintf(log_file, "%s\r\n", log_entry)
}

func Panic(msg string) {
	Log(msg)
	panic("")
}

func createLogFile(t time.Time) error {
	var (
		err          error
		year_dir     string = log_dir + "/" + strconv.Itoa(t.Year())
		month_dir    string = year_dir + "/" + strconv.Itoa(int(t.Month()))
		day_log_file string = month_dir + "/" + strconv.Itoa(t.Day()) + ".txt"
	)
	if _, err = os.Open(year_dir); err != nil {
		if err = os.Mkdir(year_dir, os.ModeDir); err != nil {
			panic("创建年份为" + strconv.Itoa(t.Year()) + "的子日志目录时发生错误\r\n")
		}
	}
	if _, err = os.Open(month_dir); err != nil {
		if err = os.Mkdir(month_dir, os.ModeDir); err != nil {
			panic("在年份目录'/" + strconv.Itoa(t.Year()) + "'下创建月份为" + strconv.Itoa(int(t.Month())) + "的子日志目录时发生错误\r\n")
		}
	}
	if _, err = os.Open(day_log_file); err != nil {
		if _, err = os.Create(day_log_file); err != nil {
			panic("在年月目录'/" + strconv.Itoa(t.Year()) + "/" + strconv.Itoa(int(t.Month())) + "'下创建日期为" + strconv.Itoa(t.Day()) + "的日志文件时发生错误\r\n")
		}
	}
	return nil
}

func getLogFilePath(t time.Time) string {
	return log_dir + "/" + strconv.Itoa(t.Year()) + "/" + strconv.Itoa(int(t.Month())) + "/" + strconv.Itoa(t.Day()) + ".txt"
}

func init() {
	var err error

	// 检查日志目录
	fmt.Println("正在检查日志目录...")
	if _, err = os.Open(log_dir); err != nil {
		panic("日志目录不存在，请检查log_dir变量\r\n")
	}

	// 获取当天的日期并检查日志目录中是否已经建立了相应日期的子日志目录
	the_day = time.Now()
	createLogFile(the_day)
	log_file_path = getLogFilePath(the_day)

	fmt.Println("日志目录检查完毕... OK")
}
