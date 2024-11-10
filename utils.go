package asynqplus

import (
	"reflect"
	"runtime"
)

// getFunctionName returns function name
func getFunctionName(i interface{}) string {
	fn := reflect.ValueOf(i)
	ptr := fn.Pointer()
	fullName := runtime.FuncForPC(ptr).Name()
	return fullName
}
