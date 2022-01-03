package category

import (
	"goblog/app/models"
	"goblog/pkg/route"
)

// 文章分类

type Category struct {
	models.BaseModel

	Name string `gorm:"type:varchar(255);not null;" valid:"name"`
}

// Link 方法生成文章链接
func (c Category) Link() string  {
	return route.Name2Url("categories.show", "id", c.GetStringID())
}

