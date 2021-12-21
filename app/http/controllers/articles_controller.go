package controllers

import (
	"fmt"
	"goblog/app/models/article"
	"goblog/pkg/logger"
	"goblog/pkg/route"
	"goblog/pkg/view"
	"gorm.io/gorm"
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
		view.Render(w,  articles, "articles.index")
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
		view.Render(w,  article, "articles.show")
	}
}

// ArticlesFormData 创建博文表单数据
//type ArticlesFormData struct {
//	Title, Body string
//	Article article.Article
//	Errors 		map[string]string
//}

// 创建文章页面
func (*ArticlesController) Create(w http.ResponseWriter, r *http.Request)  {
	// 注意 form 元素那里我们使用了 RouteName2URL 因为不需要传参 URL 参数，模板里我们直接使用 RouteName2URL 生成 URL，代码可以变得很简洁：
	//storeURL := route.Name2Url("articles.store")
	//data := ArticlesFormData{
	//	Title: "",
	//	Body: "",
	//	URL:	storeURL,
	//	Errors: nil,
	//}
	//tmpl, err := template.ParseFiles("resources/views/articles/create.gohtml")
	//if err != nil {
	//	panic(err)
	//}
	//
	//err = tmpl.Execute(w, data)
	//if err != nil {
	//	panic(err)
	//}
	view.Render(w,  view.D{}, "articles.create", "articles._form_field")
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
		view.Render(w, view.D{
			"Title": title,
			"Body": body,
			"Errors": errors,
		},"articles.create","articles._form_field")
		//storeURL := route.Name2Url("articles.store")
		//
		//data := ArticlesFormData{
		//	Title: title,
		//	Body: body,
		//	URL: storeURL,
		//	Errors: errors,
		//}
		//
		//tmpl, err := template.ParseFiles("resources/views/articles/create.gohtml")
		//if err != nil {
		//	panic(err)
		//}
		//
		//err = tmpl.Execute(w, data)
		//if err != nil {
		//	panic(err)
		//}

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


// 编辑文章
func (*ArticlesController) Edit(w http.ResponseWriter,r *http.Request) {
	// 获取URL参数
	id := route.GetRouteVariable("id", r)

	// 读取对应文章
	_article, err := article.Get(id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			// 数据未找到
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprint(w, "404 文章未找到")
		} else {
			// 数据库错误
			logger.LogError(err)
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, "500服务器内部错误")
		}
	} else {
		// 4 读取成功，显示表单
		//updateURL := route.Name2Url("articles.update", "id", id)
		//
		//data := ArticlesFormData{
		//	Title: article.Title,
		//	Body: article.Body,
		//	URL: updateURL,
		//	Errors: nil,
		//}
		//tmpl, err := template.ParseFiles("resources/views/articles/edit.gohtml")
		//logger.LogError(err)
		//err = tmpl.Execute(w, data)
		//logger.LogError(err)
		view.Render(w, view.D{
			"Title": _article.Title,
			"Body": _article.Body,
			"Article": _article,
			//"Errors": nil,
		},"articles.edit", "articles._form_field")
	}

}


// 更新文章

func (*ArticlesController) Update(w http.ResponseWriter, r *http.Request) {
	// 获取URL参数
	id := route.GetRouteVariable("id", r)

	// 读取对应文章
	_article, err := article.Get(id)

	// 判断是否出错
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			// 未查询到数据
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprint(w, "404 文章不存在")
		} else {
			// 数据库错误
			logger.LogError(err)
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, "500 服务器内部错误")
		}
	} else {
		//表单验证
		title := r.PostFormValue("title")
		body := r.PostFormValue("body")

		errors := validateArticleFormData(title, body)
		if len(errors) == 0 {
			// 更新数据库
			_article.Title 	= title
			_article.Body 	= body

			rowsAffected, err := _article.Update()
			if err != nil {
				logger.LogError(err)
				w.WriteHeader(http.StatusInternalServerError)
				fmt.Fprint(w, "500 服务器内部错误")
			}

			// 更新成功跳转文章详情页
			if rowsAffected > 0 {
				showURL := route.Name2Url("articles.show", "id", id)
				http.Redirect(w, r, showURL, http.StatusFound)
			} else {
				fmt.Fprint(w, "未做任何更新")
			}
		} else {
			// 表单验证不通过， 显示理由
			view.Render(w, view.D{
				"Title": title,
				"Body": body,
				"Article": _article,
				"Errors": errors,
			},"articles.edit", "articles._form_field")
			//updateURL := route.Name2Url("articles.update", "id", id)
			//
			//data := ArticlesFormData{
			//	Title: title,
			//	Body: body,
			//	URL: updateURL,
			//	Errors: errors,
			//}
			//tmpl, err := template.ParseFiles("resources/views/articles/edit.gohtml")
			//logger.LogError(err)
			//
			//
			//err = tmpl.Execute(w, data)
			//logger.LogError(err)

		}
	}
}


// Delete 删除文章
func (*ArticlesController) Delete(w http.ResponseWriter, r *http.Request)  {
	// 获取URL 参数
	id := route.GetRouteVariable("id", r)

	// 读取对应文章
	_article, err := article.Get(id)

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprint(w, "404 文章未找到")
		} else {
			logger.LogError(err)
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, "500 服务器内部错误")
		}
	} else {
		// 未出现错误，指向删除操作
		rowsAffected, err := _article.Delete()

		// 发生错误
		if err != nil {
			// SQL 错误
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, "500 服务器内部错误")
		} else {
			// 未发生错误
			if rowsAffected > 0 {
				// 重定向文章列表页
				indexURL := route.Name2Url("articles.index")
				http.Redirect(w, r, indexURL, http.StatusFound)
			} else {
				w.WriteHeader(http.StatusNotFound)
				fmt.Fprint(w, "404 文章为找到")
			}
		}
	}
}