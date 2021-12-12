package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"html/template"
	"net/http"
	"net/url"
	"strings"
	"unicode/utf8"
	_ "github.com/go-sql-driver/mysql"
)
// 包级别变量声明 不能 ":=" 语法声明
//router := mux.NewRouter().StrictSlash(true)
var router = mux.NewRouter()

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
	vars := mux.Vars(r)
	id := vars["id"]
	fmt.Fprint(w, "文章ID: " + id)
}

func articlesIndexHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "访问文章列表")
}


// ArticlesFormData 创建博文表单数据
type ArticlesFormData struct {
	Title, Body string
	URL 	*url.URL
	Errors	map[string]string
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

//func articlesStoreHandler(w http.ResponseWriter, r *http.Request) {
//	// 表的验证
//	title := r.PostFormValue("title")
//	body := r.PostFormValue("body")
//
//	errors := make(map[string]string)
//	// 内建函数len() 可以用来获取 切片、字符串、通道(channel) 等长度
//	// go 语言的字符都是以UTF-8格式保存，每个中文占用3个字节，因此使用len() 获得长度为 字符个数 * 3
//	// 如果需要获取字符个数需要用Go语言中utf8 提供的RuneCountInString() 函数计算
//	// 验证标题
//	//if title == "" {
//	//	errors["title"] = "标题不能为空"
//	//} else if len(title) < 3 || len(title) > 40 {
//	//	errors["title"] = "标题长度需介于 3-40"
//	//}
//
//	if title == "" {
//		errors["title"] = "标题不能为空"
//	} else if utf8.RuneCountInString(title) < 3 || utf8.RuneCountInString(title) > 40 {
//		errors["title"] = "标题长度需介于 3-40"
//	}
//
//	// 验证内容
//	if body == "" {
//		errors["body"] = "内容不能为空"
//	} else if len(body) < 10 {
//		errors["body"] = "内容长度须大于或等于 10 字节"
//	}
//
//	// 验证是否通过
//	if len(errors) == 0 {
//		fmt.Fprint(w, "验证通过！<br>")
//		fmt.Fprintf(w, "title 的值为: %v <br>", title)
//		fmt.Fprintf(w, "title 的长度为: %v <br>", len(title))
//		fmt.Fprintf(w, "body 的值为: %v <br>", body)
//		fmt.Fprintf(w, "body 的长度为: %v <br>", len(body))
//	} else {
//		fmt.Fprintf(w, "有错误发生, errors 的值为： %v <br>", errors)
//	}
//}
//func articlesStoreHandler(w http.ResponseWriter, r *http.Request) {
//	// http 包提供 从请求中解析请求参数必须执行完此代码， 后续 r.PostForm 和 r.Form 才能读取到数据 否则为空数组
//	err := r.ParseForm()
//
//	if err != nil {
//		// 解析错误处理
//		fmt.Fprint(w, "请提供正确数据")
//		return
//	}
//
//	//title := r.PostForm.Get("title")
//	//
//	//fmt.Fprintf(w, "POST PostForm: %v <br>", r.PostForm)
//	//fmt.Fprintf(w, "POST Form: %v <br>", r.Form)
//	//fmt.Fprintf(w, "title 的值为： %v", title)
//
//	fmt.Fprintf(w, "r.Form 中 title 的值为: %v <br>", r.FormValue("title"))
//	fmt.Fprintf(w, "r.PostForm 中 title 的值为: %v <br>", r.PostFormValue("title"))
//	fmt.Fprintf(w, "r.Form 中 test 的值为: %v <br>", r.FormValue("test"))
//	fmt.Fprintf(w, "r.PostForm 中 test 的值为: %v <br>", r.PostFormValue("test"))
//
//	// POST PostForm: map[body:[这里是内容] title:[这里是标题]]
//	// POST Form: map[body:[这里是内容] title:[这里是标题]]
//	// title 的值为： 这里是标题
//
//	// 根据打印结果可见 r.PostForm 和 r.Form 的数据是一样的
//	// Form 存储了 post、put和get 参数， 使用之前需要调用 parseForm 方法
//	// PostForm 存储了 post、put参数 使用之前需要调用 parseForm 方法
//	// 可见 r.Form 比 r.PostForm 多了 URL 参数里的数据。
//
//
//}

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


func main() {

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

	// 通过路由命名获取 URL
	homeURL, _ := router.Get("home").URL()
	if homeURL != nil {
		fmt.Println("homeURL: ", homeURL)
	}
	//fmt.Println("homeURL: ", homeURL)
	articleURL, _ := router.Get("article.show").URL("id", "23")
	if articleURL != nil {
		fmt.Println("articleURL: ", articleURL)
	}
	//fmt.Println("articleURL: ", articleURL)
	// url 后缀处理
	http.ListenAndServe(":3000", removeTrailingSlash(router))
}