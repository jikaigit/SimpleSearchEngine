package pq

import (
	"database/sql"
	"log"
)

var (
	// DB用来保持服务器和数据库的连接
	DBConn *sql.DB
)

// 初始化一些数据库相关的数据结构以加速对数据库的访问和操作
func init() {
	var (
		err         error
		db_conn_str string
	)

	db_conn_str += "dbname=creatogether "
	db_conn_str += "user=postgres "
	db_conn_str += "password=test "
	db_conn_str += "sslmode=disable"

	if DBConn, err = sql.Open("postgres", db_conn_str); err != nil {
		log.Fatal("database error: ", err.Error())
	}

	// TODO: 检查各表是否都已存在(不存在则创建)

}
