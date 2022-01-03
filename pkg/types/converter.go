package types

import (
	"goblog/pkg/logger"
	"strconv"
)

// Int64ToString 将int64 转为 string
func Int64ToString(num int64) string  {
	return strconv.FormatInt(num, 10)
}

// StringToUnit64 字符串转为 uint64
func StringToUnit64(str string) uint64  {
	i, err := strconv.ParseUint(str, 10 ,64)
	if err != nil {
		logger.LogError(err)
	}
	return i
}

// Uint64ToString 将 uint64 转为 string
func Uint64ToString(num uint64) string {
	return strconv.FormatUint(num, 10)
}

// StringToInt 将字符串转为 int
func StringToInt(str string) int {
	i, err := strconv.Atoi(str)

	if err != nil {
		logger.LogError(err)
	}
	return i
}