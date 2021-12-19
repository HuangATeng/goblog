package controllers

import (
	"database/sql"
	"fmt"
	"goblog/pkg/logger"
	"goblog/pkg/route"
	"goblog/pkg/types"
	"net/http"
)

// ArticlesController 文章相关页面
type ArticlesController struct {

}

// Show 文章详情页面
func (* ArticlesController) Show(w http.ResponseWriter, r *http.Request)  {
	// 获取URL 参数
	id := route.GetRouteVariable("id", r)

	// 读取对应文章
	article, err := getArticleById(id)

	if err != nil {
		if err == sql.ErrNoRows {
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
				"Int64ToString": types.Int64ToString,
			}).ParseFiles("resources/views/articles/show.gohtml")

		logger.LogError(err)

		err = tmpl.Execute(w, article)
		logger.LogError(err)
	}
}