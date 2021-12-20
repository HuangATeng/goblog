package view

import (
	"goblog/pkg/logger"
	"goblog/pkg/route"
	"html/template"
	"io"
	"path/filepath"
	"strings"
)

// render 渲染视图

func Render(w io.Writer, data interface{}, tplFiles...string)  {
	// 1 设置模板相对路径
	viewDir := "resources/views/"

	// 语法糖 将article.show 更正为articles/show  n 允许替换次数 -1 替换所有
	//name = strings.Replace(name, ".", "/", -1)
	for i, f := range tplFiles {
		tplFiles[i] = viewDir + strings.Replace(f, ".", "/", -1) + ".gohtml"
	}
	// 布局所有模板文件
	layoutFiles, err := filepath.Glob(viewDir + "layouts/*.gohtml")
	//logger.LogError(err)

	// 在 Slice 里新增目标文件
	//newFiles := append(files, viewDir + name + ".gohtml")
	allFiles := append(layoutFiles, tplFiles...)
	// 解析所有模板文件
	tmpl, err := template.New("").
		Funcs(template.FuncMap{
			"RouteName2URL": route.Name2Url,
	}).ParseFiles(allFiles...)
	logger.LogError(err)

	// 渲染模板
	err = tmpl.ExecuteTemplate(w, "app", data)
	logger.LogError(err)
}
