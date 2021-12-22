package controllers

import (
	"fmt"
	"goblog/app/models/user"
	"goblog/app/requests"
	"goblog/pkg/auth"
	"goblog/pkg/view"
	"net/http"
)

// AuthController 处理静态页面
type AuthController struct {

}

//type userForm struct {
//	Name 			string `valid:"name"`
//	Email			string `valid:"email"`
//	Password		string `valid:"password"`
//	PasswordConfirm	string `valid:"password_confirm"`
//}

// Register 注册页面
func (*AuthController) Register(w http.ResponseWriter, r *http.Request)  {
	view.Render(w, view.D{}, "auth.register")
}

// DoRegister 注册逻辑处理
func (*AuthController) DoRegister(w http.ResponseWriter, r *http.Request)  {
	// 初始化变量

	_user := user.User{
		Name:				r.PostFormValue("name"),
		Email:				r.PostFormValue("email"),
		Password:			r.PostFormValue("password"),
		PasswordConfirm:	r.PostFormValue("password_confirm"),
	}

	// 表单规则

	// 开始认证
	errs := requests.ValidateRegistrationForm(_user)

	if len(errs) > 0 {
		// 有错误
		//data, _ := json.MarshalIndent(errs, "", " ")
		//fmt.Fprint(w, string(data))
		// 表单验证不通过
		view.RenderSimple(w, view.D{
			"Errors": 	errs,
			"User": 	_user,
		},"auth.register")
	}else{
		_user.Create()

		if _user.ID > 0 {
			fmt.Fprint(w, "插入成功， ID 为"+ _user.GetStringID())
		} else {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w,"创建用户失败, 请联系管理员")
		}
	}
}

// Login 显示登录表单
func (*AuthController) Login(w http.ResponseWriter, r *http.Request) {
	view.RenderSimple(w, view.D{}, "auth.login")
}

// DoLogin 登录表单提交验证
func (*AuthController) DoLogin(w http.ResponseWriter, r *http.Request)  {
	//初始化表单数据
	email := r.PostFormValue("email")
	password := r.PostFormValue("password")

	// 尝试登录
	if err := auth.Attempt(email, password); err == nil {
		// 登录成功
		http.Redirect(w, r, "/", http.StatusFound)
	} else {
		// 失败，显示错误
		view.RenderSimple(w, view.D{
			"Error":		err.Error(),
			"Email":		email,
			"Password":		password,
		}, "auth.login")
	}
}