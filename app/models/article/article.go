package article

import (
	"goblog/app/models"
	"goblog/pkg/route"
)

// Article 文章模型

type Article struct {
	models.BaseModel
	Title	string
	Body 	string
}

// Link 方法生成文章链接
func (article Article) Link() string {
	return route.Name2Url("articles.show", "id", article.GetStringID())
}