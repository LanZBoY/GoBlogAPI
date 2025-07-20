package v1_test

import (
	"reflect"
	"runtime"
)

func handlerName(f any) string {
	return runtime.FuncForPC(reflect.ValueOf(f).Pointer()).Name()
}
