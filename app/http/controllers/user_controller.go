package controllers

import (
	"fmt"
	"goblog/app/models/article"
	"goblog/app/models/user"
	"goblog/pkg/logger"
	"goblog/pkg/route"
	"goblog/pkg/view"
	"net/http"
)

type UserController struct {
	BaseController
}

func (uc *UserController) Show(w http.ResponseWriter, r *http.Request)  {
	// 获取 URL 参数
	id := route.GetRouteVariable("id", r)

	// 读取对应文章
	_user, err := user.Get(id)

	// 是否异常
	if err != nil {
		uc.ResponseForSQLError(w, err)
	} else {
		// 显示用户文章列表
		articles, err := article.GetByUserID(_user.GetStringID())

		if err != nil {
			logger.LogError(err)
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, "500 system error")
		} else {
			view.Render(w, view.D{
				"Articles": articles,
			}, "articles.index", "articles._article_meta")
		}
	}
}
