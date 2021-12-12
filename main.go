package main

import (
	"database/sql"
	"github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"log"
	"time"
)

var router = mux.NewRouter()

var db *sql.DB

func initDB() {
	var err error
	config := mysql.Config{
		User: 						"root",
		Passwd: 					"root",
		Addr:						"127.0.0.1:3306",
		Net: 						"tcp",
		DBName: 					"goblog",
		AllowNativePasswords: 		true,
	}

	// 准备数据库连接池
	db, err = sql.Open("mysql", config.FormatDSN())
	checkError(err)

	// 设置最大连接数
	db.SetMaxOpenConns(100)
	// 设置最大空闲连接数
	db.SetMaxIdleConns(25)
	// 设置每个链接过期时间
	db.SetConnMaxLifetime(5 * time.Minute)


	// 尝试连接，失败会报错
	err = db.Ping()
	checkError(err)
}

func checkError(err error){
	if err != nil {
		log.Fatal(err)
	}
}

func createTables() {
	createArticlesSQL := `CREATE TABLE IF NOT EXISTS articles(
	id bigint(20) PRIMARY KEY AUTO_INCREMENT NOT NULL,
	title varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL,
	body longtext COLLATE utf8mb4_unicode_ci
);`
	_, err := db.Exec(createArticlesSQL)
	checkError(err)
}

func main() {
	initDB()
	createTables()
}