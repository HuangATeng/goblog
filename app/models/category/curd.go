package category

import (
	"goblog/pkg/logger"
	"goblog/pkg/model"
	"goblog/pkg/types"
)

// Create 创建分类

func (category *Category) Create() (err error)  {
	if err = model.DB.Create(&category).Error; err != nil {
		logger.LogError(err)
		return err
	}
	return nil
}

// Get 通过 ID 获取分类
func Get(idstr string) (Category, error) {
	var category Category
	id := types.StringToUnit64(idstr)
	if err := model.DB.First(&category, id).Error; err != nil {
		return category, err
	}

	return category, nil
}

// All 获取分类数据
func All() ([]Category, error)  {
	var categories []Category
	if err := model.DB.Find(&categories).Error; err != nil {
		return categories, err
	}
	return categories, nil
}