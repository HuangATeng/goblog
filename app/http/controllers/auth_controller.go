package controllers

import (
	"fmt"
	"goblog/app/models/user"
	"goblog/app/requests"
	"goblog/pkg/auth"
	"goblog/pkg/flash"
	"goblog/pkg/logger"
	"goblog/pkg/sedemail"
	"goblog/pkg/session"
	"goblog/pkg/view"
	"math/rand"
	"net/http"
	"strconv"
	"time"
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
	view.RenderSimple(w, view.D{}, "auth.register")
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
		// 表单验证不通过
		view.RenderSimple(w, view.D{
			"Errors": 	errs,
			"User": 	_user,
		},"auth.register")
	}else{
		_user.Create()

		if _user.ID > 0 {
			// 注册成功跳，登录用户并跳转首页
			flash.Success("恭喜您注册成功！")
			auth.Login(_user)
			http.Redirect(w, r, "/", http.StatusFound)
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
		flash.Success("欢迎回来！")
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

// Logout 退出登录
func (*AuthController) Logout(w http.ResponseWriter, r *http.Request)  {
	auth.Logout()
	flash.Success("您已退出登录")
	http.Redirect(w, r, "/", http.StatusFound)
}

// Retrieve 密码找回
func (*AuthController) Retrieve(w http.ResponseWriter, r *http.Request)  {
	fmt.Println("1111")
	view.RenderSimple(w, view.D{}, "auth.retrieve")
}

// Doretrieve 发送验证码
func (*AuthController) Doretrieve(w http.ResponseWriter, r *http.Request)  {
	email := r.PostFormValue("email")

	_user, err := user.GetByEmail(email)

	if err != nil {
		// 账号不存在
		view.RenderSimple(w, view.D{
			"Error":		"邮箱不存在" + err.Error(),
			"Email":		email,
		}, "auth.retrieve")
	} else {
		// 账号验证通过发送验证码
		rand.Seed(time.Now().Unix())
		num := rand.Intn(10000)
		text := fmt.Sprintf("您的验证码是：%d", num)
		session.Put("code", num)
		sedemail.SendEmail("ht19910000@163.com",_user.Email, "goblog 博客密码找回", "TFXQXEGWNOJVWOVE", text)
		view.RenderSimple(w, view.D{"message":"验证码已发送至您邮箱","success":true}, "auth.retrieve")
	}
}

// 修改密码
func (*AuthController) Update(w http.ResponseWriter, r *http.Request)  {
	view.RenderSimple(w, view.D{}, "auth.update")
}

func (*AuthController) Doupdate(w http.ResponseWriter, r *http.Request)  {
	email := r.PostFormValue("email")
	code, _ := strconv.Atoi(r.PostFormValue("code"))
	password := r.PostFormValue("password")
	//_code := 5597
	_code := session.Get("code")
	logger.LogInfo(_code)
	if code != _code {
		view.RenderSimple(w, view.D{
			"errorMessage":		"验证码错误",
			"Email":		email,
			"code" :		code,
			"password":		password,
		}, "auth.update")
		return
	}

	_user, err := user.GetByEmail(email)

	if err != nil {
		// 账号不存在
		view.RenderSimple(w, view.D{
			"Error":		"邮箱不存在" + err.Error(),
			"Email":		email,
		}, "auth.update")
	} else {
		//_user.Password = password.Hash(_password)
		err := _user.UpdatePassword(password)
		if err != nil {
			view.RenderSimple(w, view.D{
				"Error":		"修改失败:" + err.Error(),
				"Email":		email,
			}, "auth.update")
		}else {
			// 修改成功跳转登录
			flash.Success("密码修改成功！")
			http.Redirect(w, r, "/auth/login", http.StatusFound)
		}
	}

}