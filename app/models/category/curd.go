package category

import (
	"goblog/pkg/logger"
	"goblog/pkg/model"
)

// Create 创建分类

func (category *Category) Create() (err error)  {
	if err = model.DB.Create(&category).Error; err != nil {
		logger.LogError(err)
		return err
	}
	return nil
}
