package sedemail

import (
	"fmt"
	"github.com/jordan-wright/email"
	"goblog/pkg/logger"
	"net/smtp"
)

func SendEmail(from, _email, subject, password, text string)  {

	em := email.NewEmail()
	// 设置 sender 发送方邮箱
	em.From = from

	// 设置接收方邮箱
	em.To = []string{_email}

	// 设置主题
	em.Subject = subject

	// 邮件发送内容
	em.Text = []byte(text)

	// 设置服务器设置相关
	err := em.Send("smtp.163.com:25", smtp.PlainAuth("", em.From, password, "smtp.163.com"))

	if err != nil {
		logger.LogError(err)
		fmt.Println(err)
	}
	logger.LogInfo(text)
}