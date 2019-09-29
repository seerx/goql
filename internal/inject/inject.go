package inject

import (
	"errors"
	"reflect"

	"github.com/graphql-go/graphql"

	"github.com/seerx/goql/internal/reflects"
)

const (
	KeyOfConext  = "context"
	KeyOfRequest = "submit"
)

// InjectInfo 注入信息
type InjectInfo struct {
	Type reflect.Type  // 注入类型
	Func reflect.Value // 生成注入对象的函数
	Info *reflects.FuncInfo

	//FuncName string        // 函数名称
	//Package  string        // 包名称
}

func (i *InjectInfo) InjectValue(p *graphql.ResolveParams) reflect.Value {
	root, ok := p.Info.RootValue.(map[string]interface{})
	if !ok {
		panic(errors.New("Param graphql.ResolveParams.Info.RootValue is not type of map[string]interface{}"))
	}
	ctx, ok := root[KeyOfConext]
	if !ok {
		panic(errors.New("No context.Context in graphql.ResolveParams.Info.RootValue"))
	}
	req, ok := root[KeyOfRequest]
	if !ok {
		panic(errors.New("No http.Request in graphql.ResolveParams.Info.RootValue"))
	}

	args := []reflect.Value{
		reflect.ValueOf(ctx),
		reflect.ValueOf(req),
		reflect.ValueOf(p),
	}
	res := i.Func.Call(args)
	return res[0]
}
