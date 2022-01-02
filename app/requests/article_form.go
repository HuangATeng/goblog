package requests

import (
	"github.com/thedevsaddam/govalidator"
	"goblog/app/models/article"
)

// ValidateArticleForm 验证表单 返回 errs 长度等于零即通过
func ValidateArticleForm(data article.Article) map[string][]string {
	// 定制认证规则
	rules := govalidator.MapData{
		"title": []string{"required", "min:3", "max:40"},
		"body": []string{"required", "min:10"},
	}
	// 错误消息
	messages := govalidator.MapData{
		"title": []string{
			"required:标题必填",
			"min:标题长度须大于 3",
			"max:标题长度需小于 40",
		},
		"body": []string{
			"required: 文章内容为必填",
			"min:长度须大于 10",
		},
	}

	// 配置初始化
	opts := govalidator.Options{
		Data:			&data,
		Rules: 			rules,
		TagIdentifier: 	"valid", // 模型中的 Struct 标签标示符
		Messages:		messages,
	}

	// 开始验证
	return govalidator.New(opts).ValidateStruct()

}
