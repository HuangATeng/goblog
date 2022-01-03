package controllers

import (
	"fmt"
	"goblog/app/models/article"
	"goblog/app/models/category"
	"goblog/app/requests"
	"goblog/pkg/flash"
	"goblog/pkg/route"
	"goblog/pkg/view"
	"net/http"
)

type CategoriesController struct {
	BaseController
}

// Show 显示文章分类
func (cc *CategoriesController) Show(w http.ResponseWriter, r *http.Request)  {
	// 获取URL 参数
	id := route.GetRouteVariable("id", r)

	// 读取数据
	_category, err := category.Get(id)

	// 获取结果
	articles, pagerData, err := article.GetByCategoryID(_category.GetStringID(), r, 10)

	if err != nil {
		cc.ResponseForSQLError(w, err)
	} else {
		// 加载模板
		view.Render(w, view.D{
			"Articles": articles,
			"PagerData": pagerData,
		}, "articles.index", "articles._article_meta")
	}
}

// Create 文章分类创建页面
func (*CategoriesController) Create(w http.ResponseWriter, r *http.Request) {
	view.Render(w, view.D{}, "categories.create")
}

// Store 保存文章分类
func (*CategoriesController) Store(w http.ResponseWriter, r *http.Request)  {
	// 初始化数据
	_category := category.Category{
		Name: r.PostFormValue("name"),
	}

	// 表单验证
	errors := requests.ValidateCategoryForm(_category)

	// 错误检测
	if len(errors) == 0 {
		// 创建文章分类
		_category.Create()
		if _category.ID > 0 {
			//fmt.Fprint(w, "创建成功！")
			flash.Success("分类创建成功")
			indexURL := route.Name2Url("home")
			http.Redirect(w, r, indexURL, http.StatusFound)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, "创建文章分类失败，请联系管理员")
		}
	} else {
		view.Render(w, view.D{
			"Category": _category,
			"Errors":	errors,
		}, "categories.create")
	}
}