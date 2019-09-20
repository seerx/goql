package reflects

import (
	"fmt"
	"reflect"
	"runtime"
	"strings"
)

// FuncInfo 函数信息
type FuncInfo struct {
	Name    string // 名称
	Struct  string // 所属结构
	Package string // 包
}

func (fi *FuncInfo) String() string {
	if fi.Struct == "" {
		return fmt.Sprintf("%s@%s", fi.Name, fi.Package)
	}
	return fmt.Sprintf("%s.%s@%s", fi.Struct, fi.Name, fi.Package)
}

// ParseFuncInfo 解析函数信息
func ParseFuncInfo(aFunc interface{}) *FuncInfo {
	// 获取函数名称
	fn := runtime.FuncForPC(reflect.ValueOf(aFunc).Pointer()).Name()

	p := strings.LastIndex(fn, ".")

	if p > 0 {
		return &FuncInfo{
			Name:    fn[p+1:],
			Package: fn[0:p],
		}
	}

	return nil
}
