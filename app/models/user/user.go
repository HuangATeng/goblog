package user

import (
	"goblog/app/models"
	"goblog/pkg/logger"
	"goblog/pkg/model"
)

// User 用户模型
type User struct {
	models.BaseModel

	Name 		string `gorm:"type:varchar(255);not null;unique" valid:"name"`
	Email		string `gorm:"type:varchar(255);unique" valid:"email"`
	Password	string `gorm:"type:varchar(255)" valid:"password"`

	// gorm:"-" -- 设置 GORM 在读写时略过此字段， 仅用于表单验证
	PasswordConfirm string `gorm:"-" valid:"password_confirm"`

}


// Create 创建用户， 通过 User.ID 判断是否成功
func (user *User) Create() (err error) {
	if err = model.DB.Create(&user).Error; err != nil {
		logger.LogError(err)
		return err
	}

	return nil
}