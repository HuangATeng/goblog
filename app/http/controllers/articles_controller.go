package controllers

import (
	"fmt"
	"goblog/app/models"
	"goblog/app/models/article"
	"goblog/app/policies"
	"goblog/app/requests"
	"goblog/pkg/auth"
	"goblog/pkg/flash"
	"goblog/pkg/logger"
	"goblog/pkg/route"
	"goblog/pkg/view"
	"net/http"
	"unicode/utf8"
)

// ArticlesController 文章相关页面
type ArticlesController struct {
	BaseController
	models.BaseModel

	Title	string `gorm:"type:varchar(255);not null;" valid:"title"`
	Body 	string `gorm:"type:longtext;not null"; valid:"body"`
}

// 文章列表页
func (ac *ArticlesController) Index(w http.ResponseWriter, r *http.Request)  {
	//fmt.Fprint(w, config.Get("app.name"))
	// 获取结果集
	articles,pagerData, err := article.GetAll(r, 10)

	if err != nil {
		// 数据库错误
		ac.ResponseForSQLError(w, err)
	} else {
		view.Render(w,  view.D{
			"Articles": articles,
			"PagerData": pagerData,
		}, "articles.index", "articles._article_meta")
	}
}

// Show 文章详情页面
func (ac *ArticlesController) Show(w http.ResponseWriter, r *http.Request)  {
	// 获取URL 参数
	//id := route.GetRouteVariable("id", r)
	id := route.GetRouteVariable("id", r)
	// 读取对应文章
	article, err := article.Get(id)
	if err != nil {
		ac.ResponseForSQLError(w, err)
	} else {
		//
		view.Render(w,  view.D{
			"Article": article,
		}, "articles.show", "articles._article_meta")
	}
}


// 创建文章页面
func (*ArticlesController) Create(w http.ResponseWriter, r *http.Request)  {
	// 注意 form 元素那里我们使用了 RouteName2URL 因为不需要传参 URL 参数，模板里我们直接使用 RouteName2URL 生成 URL，代码可以变得很简洁：
	view.Render(w,  view.D{}, "articles.create", "articles._form_field")
}

// 保存文章
func (*ArticlesController) Store(w http.ResponseWriter, r *http.Request)  {
	// 初始化数据
	currentUser := auth.User()
	_article := article.Article{
		Title:	r.PostFormValue("title"),
		Body:	r.PostFormValue("body"),
		UserID: currentUser.ID,
	}

	// 表的验证
	errors := requests.ValidateArticleForm(_article)

	// 检查是否有错误
	if len(errors) == 0 {

		_article.Create()
		if _article.ID > 0 {
			fmt.Println(_article.GetStringID())
			indexUrl := route.Name2Url("articles.show", "id", _article.GetStringID())
			http.Redirect(w, r, indexUrl, http.StatusFound)
		}else {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, "文章创建失败，请联系管理员")
		}
	} else {
		view.Render(w, view.D{
			"Article": _article,
			"Errors": errors,
		},"articles.create","articles._form_field")

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
func (ac *ArticlesController) Edit(w http.ResponseWriter,r *http.Request) {
	// 获取URL参数
	id := route.GetRouteVariable("id", r)
	fmt.Println(id)
	// 读取对应文章
	_article, err := article.Get(id)
	if err != nil {
		ac.ResponseForSQLError(w, err)
	} else {
		// 权限检查
		if !policies.CanModifyArticle(_article) {
			flash.Warning("未授权操作！")
			http.Redirect(w, r, "/", http.StatusFound)
		} else {
			// 4 读取成功，显示表单
			view.Render(w, view.D{
				"Errors":	view.D{},
				"Article": 	_article,
				//"Errors": nil,
			},"articles.edit", "articles._form_field")
		}

	}

}


// 更新文章

func (ac *ArticlesController) Update(w http.ResponseWriter, r *http.Request) {
	// 获取URL参数
	id := route.GetRouteVariable("id", r)

	// 读取对应文章
	_article, err := article.Get(id)

	// 判断是否出错
	if err != nil {
		ac.ResponseForSQLError(w, err)
	} else {
		// 权限检查
		if !policies.CanModifyArticle(_article) {
			flash.Warning("未授权操作")
			http.Redirect(w, r, "/", http.StatusForbidden)
		} else {
			//表单验证通过
			title := r.PostFormValue("title")
			body := r.PostFormValue("body")

			errors := validateArticleFormData(title, body)
			//errors := requests.ValidateArticleForm(_article)
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
					"Article": _article,
					"Errors": errors,
				},"articles.edit", "articles._form_field")

			}
		}

	}
}


// Delete 删除文章
func (ac *ArticlesController) Delete(w http.ResponseWriter, r *http.Request)  {
	// 获取URL 参数
	id := route.GetRouteVariable("id", r)

	// 读取对应文章
	_article, err := article.Get(id)

	if err != nil {
		ac.ResponseForSQLError(w, err)
	} else {
		// 权限检查
		if !policies.CanModifyArticle(_article) {
			flash.Warning("您没有权限执行此操作！")
			http.Redirect(w, r, "/", http.StatusFound)
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
}