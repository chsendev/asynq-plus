package asynqplus

import (
	"reflect"
	"runtime"
)

// GetFunctionName 返回给定函数的名称
func GetFunctionName(i interface{}) string {
	// 获取函数的反射值
	fn := reflect.ValueOf(i)
	// 获取函数的指针
	ptr := fn.Pointer()
	// 获取函数名，包括包名
	fullName := runtime.FuncForPC(ptr).Name()
	//shortName := fullName[strings.LastIndex(fullName, "/")+1:]
	return fullName
}
