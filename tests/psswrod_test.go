package tests

import (
	"fmt"
	"goblog/pkg/password"
	"goblog/pkg/sedemail"
	"math/rand"
	"testing"
	"time"
)

func TestPasswod(t *testing.T)  {
	//$2a$14$chn0jLyrdgAYtHHnMfpb3u8Ck3eDGlySpWH8g8fyrIS0CEH8tTPpe
	str := password.Hash("123456789")
	fmt.Println(str)
}

// 邮件发送
func TestSendEmail(t *testing.T)  {
	rand.Seed(time.Now().Unix())
	num := rand.Intn(10000)
	text := fmt.Sprintf("您的验证码是：%d", num)
	sedemail.SendEmail("ht19910000@163.com","1101955127@qq.com", "goblog 博客密码找回", "TFXQXEGWNOJVWOVE", text)
}


