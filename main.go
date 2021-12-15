package main

import (
	"database/sql"
	"fmt"
	"github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"html/template"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"
)

var router = mux.NewRouter()

var db *sql.DB

func homeHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "<h1>Hello welcome to goblog</h1>")
}

func aboutHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "此博客是用以记录编程笔记"+
		"<a href=\"https://huangateng.github.io/\">huangateng.github.io/</a>")
}

func notFundHandler(w http.ResponseWriter, r *http.Request){
	w.WriteHeader(http.StatusNotFound)
	fmt.Fprint(w,"<h1>not found</h1>")
}

func articlesShowHandler(w http.ResponseWriter, r *http.Request)  {
	// 获取URL 参数
	vars := mux.Vars(r)
	id := vars["id"]

	// 读取对应文章
	article := Article{}
	query := "SELECT * FROM articles WHERE id = ?"
	err := db.QueryRow(query, id).Scan(&article.ID, &article.Title, &article.Body)

	if err != nil {
		if err == sql.ErrNoRows {
			// 3.1 未找到数据
			w.WriteHeader(http.StatusNotFound)
		} else {
			// 数据错误
			checkError(err)
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, "500 服务器内部错误")
		}
	} else {
		//
		tmpl, err := template.ParseFiles("resources/views/articles/show.gohtml")
		//fmt.Fprint(w, "读取成功，文章标题 -- " + article.Title)
		checkError(err)

		err = tmpl.Execute(w, article)
		checkError(err)
	}

	//fmt.Fprint(w, "文章ID: " + id)
}

func articlesIndexHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "访问文章列表")
}

// 创建文章列表t
func articlesCreateHandler(w http.ResponseWriter, r *http.Request) {
	html := `<!DOCTYPE html>
	<html lang="en">
	<head>
		<title>创建文章 —— 我的技术博客</title>
	</head>
	<body>
		<form action="%s?test=data" method="post">
			<p><input type="text" name="title"></p>
			<p><textarea name="body" cols="30" rows="10"></textarea></p>
			<p><button type="submit">提交</button></p>
		</form>
	</body>
	</html>`

	storeURL, _ := router.Get("articles.store").URL()
	fmt.Fprintf(w, html, storeURL)
}


// ArticlesFormData 创建博文表单数据
type ArticlesFormData struct {
	Title, Body string
	URL 	*url.URL
	Errors	map[string]string
}

// Article 对应一条文章记录
type Article struct {
	Title, Body string
	ID 			int64
}


func articlesStoreHandler(w http.ResponseWriter, r *http.Request) {
	// 表的验证
	title := r.PostFormValue("title")
	body := r.PostFormValue("body")

	errors := make(map[string]string)
	// 内建函数len() 可以用来获取 切片、字符串、通道(channel) 等长度
	// go 语言的字符都是以UTF-8格式保存，每个中文占用3个字节，因此使用len() 获得长度为 字符个数 * 3
	// 如果需要获取字符个数需要用Go语言中utf8 提供的RuneCountInString() 函数计算
	// 验证标题
	//if title == "" {
	//	errors["title"] = "标题不能为空"
	//} else if len(title) < 3 || len(title) > 40 {
	//	errors["title"] = "标题长度需介于 3-40"
	//}

	if title == "" {
		errors["title"] = "标题不能为空"
	} else if utf8.RuneCountInString(title) < 3 || utf8.RuneCountInString(title) > 40 {
		errors["title"] = "标题字符个数需介于 3-40"
	}

	// 验证内容
	if body == "" {
		errors["body"] = "内容不能为空"
	} else if len(body) < 10 {
		errors["body"] = "内容长度须大于或等于 10 字节"
	}
	// 检查是否有错误
	if len(errors) == 0 {
		lastInsertID, err := saveArticleToDB(title, body)
		if lastInsertID > 0 {
			fmt.Fprintf(w, "插入成功， ID 为" + strconv.FormatInt(lastInsertID,10))
		}else {
			checkError(err)
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, "500 服务器内部错误")
		}
	} else {

		storeURL, _ := router.Get("articles.store").URL()

		data := ArticlesFormData{
			Title: title,
			Body: body,
			URL: storeURL,
			Errors: errors,
		}

		//tmpl, err := template.New("create-form").Parse(html)
		tmpl, err := template.ParseFiles("resources/views/articles/create.gohtml")
		if err != nil {
			panic(err)
		}

		err = tmpl.Execute(w, data)
		if err != nil {
			panic(err)
		}

	}
}

// 保存文章到数据库
func saveArticleToDB(title string, body string) (int64, error) {
	// 变量初始化
	var (
		id 	int64
		err error
		rs	sql.Result
		stmt *sql.Stmt
	)
	fmt.Println("INSERT INTO articles (title, body) VALUES (?,?)")
	// 获取一个 prepare 声明语句
	stmt, err = db.Prepare("INSERT INTO articles (title, body) VALUES (?,?)")
	// 列行错误检测
	if err != nil {
		return 0, err
	}
	// 在此函数运行结束后关闭此语句，防止占用 SQL连接
	defer stmt.Close()

	// 执行请求，传参进入绑定的内容
	rs, err = stmt.Exec(title, body)
	if err != nil {
		return 0, err
	}

	// 插入成功返回 自增长ID
	if id, err = rs.LastInsertId(); id > 0 {
		return id, nil
	}
	return 0, err
}

// html标头中间件
func forceHTMLMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 设置标头
		w.Header().Set("Content-Type", "text/html; charset=utf-8")

		// 继续处理请求
		next.ServeHTTP(w, r)
	})
}

// 去除URL后缀 "/"
func removeTrailingSlash(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 除首页以外，移除所有请求路径后面的斜杠
		if r.URL.Path != "/" {
			r.URL.Path = strings.TrimSuffix(r.URL.Path, "/")
		}
		// 将请求传递下去
		next.ServeHTTP(w, r)
	})
}


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

	//router := mux.NewRouter()
	//处理斜杠问题 localhost:3000/about/ 出现404解决
	// 可以看到当请求 about/ 时产生了两个请求，第一个是 301 跳转，第二个是跳转到的 about 去掉斜杆的链接。
	// 这个解决方案看起来不错，然而有一个严重的问题 —— 当请求方式为 POST 的时候，遇到服务端的 301 跳转，将会变成 GET 方式。很明显，这并非所愿，我们需要一个更好的方案

	router.HandleFunc("/", homeHandler).Methods("GET").Name("home")
	router.HandleFunc("/about", aboutHandler).Methods("GET").Name("about")
	router.HandleFunc("/articles/{id:[0-9]+}",articlesShowHandler).Methods("GET").Name("article.show")
	router.HandleFunc("/articles", articlesIndexHandler).Methods("GET").Name("articles.index")

	// 创建表单
	router.HandleFunc("/articles/create", articlesCreateHandler).Methods("GET").Name("articles.create")


	router.HandleFunc("/articles", articlesStoreHandler).Methods("POST").Name("articles.store")
	// 自定义 404页面
	router.NotFoundHandler = http.HandlerFunc(notFundHandler)

	// 中间件： 强制内容类型为 HTML
	router.Use(forceHTMLMiddleware)


	// url 后缀处理
	http.ListenAndServe(":3000", removeTrailingSlash(router))
}