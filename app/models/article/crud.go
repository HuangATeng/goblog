package article

import (
	"goblog/pkg/logger"
	"goblog/pkg/model"
	"goblog/pkg/pagination"
	"goblog/pkg/route"
	"goblog/pkg/types"
	"net/http"
)

// Get 通过 ID 获取文章

func Get(idstr string) (Article, error) {
	var article Article
	id := types.StringToUnit64(idstr)

	if err := model.DB.Debug().Preload("User").First(&article, id).Error; err != nil {
		return article, err
	}

	return article, nil
}

// GetByUserID 获取用户全部文章
func GetByUserID(uid string) ([]Article ,error)  {
	var articles []Article
	if err := model.DB.Where("user_id = ?", uid).Preload("User").Find(&articles).Error; err != nil {
		return articles, err
	}
	return articles, nil
}

// 获取全部文章
func GetAll(r *http.Request, perPage int) ([]Article, pagination.ViewData ,error)  {


	// 初始化分页实例
	db := model.DB.Model(Article{}).Order("created_at desc")
	_pager := pagination.New(r, db, route.Name2Url("articles.index"), perPage)

	// 获取视图数据
	viewData := _pager.Paging()

	// 获取数据
	var articles []Article
	_pager.Results(&articles)
	//if err := model.DB.Preload("User").Find(&articles).Error; err != nil {
	//	return articles, err
	//}
	return articles, viewData, nil
}


// Create 创建文章，通过article.ID 来判断是否创建成功
func (article *Article) Create() (err error) {
	result := model.DB.Debug().Create(&article)
	if err = result.Error; err != nil {
		return err
	}

	return nil
}

// 更新文章
func (article *Article) Update() (rowsAffected int64, err error)  {
	result := model.DB.Save(&article)

	if err = result.Error; err != nil {
		logger.LogError(err)
		return 0, err
	}

	return result.RowsAffected, nil
}

// 删除文章
func (article *Article) Delete() (rowsAffected int64, err error)  {
	result := model.DB.Delete(&article)
	if err = result.Error; err != nil {
		logger.LogError(err)
		return 0, err
	}
	return result.RowsAffected, nil
}