package user

import (
	"goblog/app/models"
	"goblog/pkg/logger"
	"goblog/pkg/model"
	"goblog/pkg/password"
	"goblog/pkg/types"
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

// Get 通过 ID 获取用户
func Get(idstr string) (User, error)  {
	var user User

	id := types.StringToUnit64(idstr)
	if err := model.DB.First(&user, id).Error; err != nil {
		return user, err
	}
	return user, nil
}

// GetByEmail 通过 Email 来获取用户
func GetByEmail(email string) (User, error)  {
	var user User
	if err := model.DB.Where("email = ?", email).First(&user).Error; err != nil {
		return user, err
	}

	return user, nil
}

// ComparePassword 对比密码是否匹配
func (user *User) ComparePassword(_password string) bool  {
	//return user.Password == password
	return password.CheckHash(_password, user.Password)
}

// UpdateUser 更新用户信息
func (user *User) UpdatePassword(_password string) (err error)  {
	// Updates,Update 更新操作会自动运行 model hooks 的 BeforeUpdate, AfterUpdate 方法，更新 updated_at 时间戳,
	// 如果你不想调用这些方法，你可以使用 UpdateColumn， UpdateColumns
	if err := model.DB.Model(&user).UpdateColumn("Password", password.Hash(_password)).Error; err != nil {
		logger.LogError(err)
		return err
	}
	return nil
	//db.Model(&product).Updates(map[string]interface{}{"price": gorm.Expr("price * ? + ?", 2, 100)})

}

// Link 生成用户连接
func (user *User) Link() string {
	return ""
}