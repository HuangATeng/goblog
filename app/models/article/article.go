package article

import (
	"goblog/pkg/route"
	"strconv"
)

// Article 文章模型

type Article struct {
	ID		uint64
	Title	string
	Body 	string
}

// Link 方法生成文章链接
func (article Article) Link() string {
	return route.Name2Url("articles.show", "id", strconv.FormatUint(article.ID, 10))
}