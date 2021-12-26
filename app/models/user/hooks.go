package user

import (
	"goblog/pkg/logger"
	"goblog/pkg/password"
	"gorm.io/gorm"
)

/*
// BeforeCreate GORM 的模型钩子， 创建模型前调用
func (user *User) BeforeCreate(tx *gorm.DB) (err error)  {
	user.Password = password.Hash(user.Password)
	fmt.Println(user)
	return
}

// BeforeUpdate GORM 的模型钩子，更新模型前调用
func (user *User) BeforeUpdate(tx *gorm.DB) (err error)  {
	if !password.IsHashed(user.Password) {
		user.Password = password.Hash(user.Password)
	}
	return
}
*/

/**
	注意BeforeSave，AfterSave在Create和Update时也会调用。
	这意味着，如果你同时定义了BeforeSave和BeforeCreate，那么在执行Create时，两者都会被触发。
 */
// BeforeSave GORM 模型钩子，在保存和更新模型前调用
func (user *User) BeforeSave(tx *gorm.DB) (err error){
	logger.LogInfo(user)
	if !password.IsHashed(user.Password) {
		user.Password = password.Hash(user.Password)
	}
	return
}
