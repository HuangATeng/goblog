package config

import (
	"github.com/spf13/cast"
	"github.com/spf13/viper"
	"goblog/pkg/logger"
)

// Viper Viper 库实例
var Viper *viper.Viper

// StrMap 简写 -- map[string]interface{}
type StrMap map[string]interface{}

// init() 函数在import 时候立即被加载
func init()  {
	// 初始化 Viper 库
	Viper = viper.New()

	// 设置文件名称
	viper.SetConfigName(".env")

	// 配置类型 支持 "json", "toml", "yaml", "properties",
	// "props", "prop", "env", "dotenv"
	Viper.SetConfigType("env")

	// 环境变量配置文件查找路径 相对于main.go
	viper.AddConfigPath(".")

	// 开始读取根目录下的 .env 文件，读不到会报错
	err := Viper.ReadInConfig()
	logger.LogError(err)

	// 设置环境变量前缀，用以区分 Go 系统环境变量
	Viper.SetEnvPrefix("appenv")

	// Viper.Get() 时，优先读取环境变量
	Viper.AutomaticEnv()

}

// Env 环境变量读取，支持默认值
func Env(envName string, defaultValue ...interface{}) interface{} {
	if len(defaultValue) > 0 {
		return Get(envName, defaultValue[0])
	}

	return Get(envName)

}

// Add 新增配置文件
func Add(name string, configuration map[string]interface{})  {
	Viper.Set(name, configuration)
}

// Get 获取配置项， 允许使用点式获取，如： app.name
func Get(path string, defaultValue ...interface{}) interface{}  {
	// 不存在情况
	if !Viper.IsSet(path) {
		if len(defaultValue) > 0 {
			return defaultValue[0]
		}
		return nil
	}
	return Viper.Get(path)
}

// GetString 获取string 类型配置信息
func GetString(path string, defaultValue ...interface{}) string {
	return cast.ToString(Get(path, defaultValue...))
}

// GetInt 获取Int 类型配置信息
func GetInt(path string, defaultValue ...interface{}) int {
	return cast.ToInt(Get(path, defaultValue...))
}

// GetUnit 获取 Unit 类型配置信息
func GetUnit(path string, defaultVales ...interface{}) bool {
	return cast.ToBool(Get(path, defaultVales...))
}

// GetInt64 获取 Int64 类型配置信息
func GetInt64(path string, defaultValue ...interface{}) int64 {
	return cast.ToInt64(Get(path, defaultValue...))
}

// GetBool 获取 Bool 类型配置
func GetBool(path string, defaultValue ...interface{}) bool {
	return cast.ToBool(Get(path, defaultValue...))
}