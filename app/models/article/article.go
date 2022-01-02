package article

import (
	"goblog/app/models"
	"goblog/app/models/user"
	"goblog/pkg/route"
)

// Article 文章模型

type Article struct {
	models.BaseModel
	Title	string `gorm:"type:varchar(255);not null;" valid:"title"`
	Body 	string `gorm:"type:longtext;not nul;" valid:"body"`

	UserID 	uint64 `gorm:"not null;index"`
	User 	user.User
}

// Link 方法生成文章链接
func (article Article) Link() string {
	return route.Name2Url("articles.show", "id", article.GetStringID())
}

// CreateAtDate 创建日期
func (article Article) CreatedAtDate() string  {
	return article.CreatedAt.Format("2006-01-02")
}