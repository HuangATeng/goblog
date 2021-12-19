package controllers

import (
	"fmt"
	"goblog/app/models/article"
	"goblog/pkg/logger"
	"goblog/pkg/route"
	"goblog/pkg/types"
	"gorm.io/gorm"
	"html/template"
	"net/http"
	"strconv"
	"unicode/utf8"
)

// ArticlesController 文章相关页面
type ArticlesController struct {

}

// 文章列表页
func (*ArticlesController) Index(w http.ResponseWriter, r *http.Request)  {
	// 获取结果集
	articles, err := article.GetAll()

	if err != nil {
		// 数据库错误
		logger.LogError(err)
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, "500 服务器内部错误")
	} else {
		// 加载模板
		tmpl, err := template.ParseFiles("resources/views/articles/index.gohtml")
		logger.LogError(err)

		// 渲染模板
		err = tmpl.Execute(w, articles)
		logger.LogError(err)
	}
}

// Show 文章详情页面
func (* ArticlesController) Show(w http.ResponseWriter, r *http.Request)  {
	// 获取URL 参数
	//id := route.GetRouteVariable("id", r)
	id := route.GetRouteVariable("id", r)
	// 读取对应文章
	article, err := article.Get(id)

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			// 3.1 未找到数据
			w.WriteHeader(http.StatusNotFound)
		} else {
			// 数据错误
			logger.LogError(err)
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, "500 服务器内部错误")
		}
	} else {
		//
		//tmpl, err := template.ParseFiles("resources/views/articles/show.gohtml")
		tmpl, err := template.New("show.gohtml").
			Funcs(template.FuncMap{
				"RouteName2URL": route.Name2Url,
				"Uint64ToString": types.Uint64ToString,
			}).ParseFiles("resources/views/articles/show.gohtml")

		logger.LogError(err)

		err = tmpl.Execute(w, article)
		logger.LogError(err)
	}
}

// ArticlesFormData 创建博文表单数据
type ArticlesFormData struct {
	Title, Body string
	URL 		string
	Errors 		map[string]string
}

// 创建文章页面
func (*ArticlesController) Create(w http.ResponseWriter, r *http.Request)  {
	storeURL := route.Name2Url("articles.store")
	data := ArticlesFormData{
		Title: "",
		Body: "",
		URL:	storeURL,
		Errors: nil,
	}
	tmpl, err := template.ParseFiles("resources/views/articles/create.gohtml")
	if err != nil {
		panic(err)
	}

	err = tmpl.Execute(w, data)
	if err != nil {
		panic(err)
	}
}

// 保存文章
func (*ArticlesController) Store(w http.ResponseWriter, r *http.Request)  {
	// 表的验证
	title := r.PostFormValue("title")
	body := r.PostFormValue("body")

	errors := validateArticleFormData(title, body)

	// 检查是否有错误
	if len(errors) == 0 {
		_article := article.Article{
			Title: title,
			Body: body,
		}
		_article.Create()

		if _article.ID > 0 {
			fmt.Fprintf(w, "插入成功， ID 为" + strconv.FormatUint(_article.ID,10))
		}else {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, "文章创建失败，请联系管理员")
		}
	} else {

		storeURL := route.Name2Url("articles.store")

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

// 表单验证
func validateArticleFormData(title string, body string) map[string]string {
	errors := make(map[string]string)

	//验证标题
	if title == "" {
		errors["title"] = "标题不能为空"
	} else if utf8.RuneCountInString(title) < 3 || utf8.RuneCountInString(title) > 40 {
		errors["title"] = "标题长度需介于3-40个字符"
	}

	// 验证body
	if body == "" {
		errors["body"] = "内容不能为空"
	} else if utf8.RuneCountInString(body) < 10 {
		errors["body"] = "内容长度须大于等于 10 个字符"
	}

	return errors
}