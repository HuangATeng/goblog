package view

import (
	"goblog/app/models/category"
	"goblog/app/models/user"
	"goblog/pkg/auth"
	"goblog/pkg/flash"
	"goblog/pkg/logger"
	"goblog/pkg/route"
	"html/template"
	"io"
	"path/filepath"
	"strings"
)

// D 是 map[string]interface{} 简写
type D map[string]interface{}


// render 渲染通用视图

func Render(w io.Writer, data D, tplFiles...string)  {
	RenderTemplate(w, "app", data, tplFiles...)
}

// RenderSimple 渲染简单视图
func RenderSimple(w io.Writer, data D, tplFiles...string)  {
	RenderTemplate(w, "simple", data, tplFiles...)
}

// RenderTemplate 渲染视图
func RenderTemplate(w io.Writer, name string, data D, tplFiles...string)  {
	// 1通用模板数据
	data["isLogined"] = auth.Check()
	data["loginUser"] = auth.User
	data["flash"] = flash.All()
	data["Users"], _ = user.All()
	data["Categories"], _ = category.All()
	// 生成模板文件
	allFiles := getTemplateFiles(tplFiles...)

	// 解析所有模板文件
	tmpl, err := template.New("").
		Funcs(template.FuncMap{
			"RouteName2URL": route.Name2Url,
		}).ParseFiles(allFiles...)
	logger.LogError(err)

	// 渲染模板
	err = tmpl.ExecuteTemplate(w, name, data)
	logger.LogError(err)
}

func getTemplateFiles(tplFiles ...string) []string {
	// 设置模板相对路径
	viewDir := "resources/views/"

	// 遍历传参文件列表，Slice, 设置正确路径，支持dir.filename 语法糖
	for i, f := range tplFiles {
		tplFiles[i] = viewDir + strings.Replace(f, ".", "/", -1) + ".gohtml"
	}

	// 所有布局模板文件 Slice
	layoutFiles, err := filepath.Glob(viewDir + "layouts/*.gohtml")
	logger.LogError(err)

	// 合并所有文件
	return append(layoutFiles, tplFiles...)
}