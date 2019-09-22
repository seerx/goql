package gqlh

import (
	"fmt"
	"reflect"

	"github.com/seerx/goql/pkg/log"

	"github.com/graphql-go/graphql"

	"github.com/seerx/goql/internal/parser"
	"github.com/seerx/goql/internal/reflects"
)

// ResolverManager 功能接口管理
type ResolverManager struct {
	log           log.Logger
	resolverFuncs []*parser.FuncDef
	resolverMap   map[string]*parser.FuncDef
}

// NewResolverManager 创建
func NewResolverManager(log log.Logger) *ResolverManager {
	return &ResolverManager{
		log:         log,
		resolverMap: map[string]*parser.FuncDef{},
	}
}

func (rm *ResolverManager) GenerateResolvers(ivp *parser.InputVarsPool,
	ovp *parser.OutputVarsPool,
	callback func(name string, resolver *graphql.Field)) {
	//vm := NewOutputVarsPool(rm.log)
	for _, fn := range rm.resolverFuncs {
		out := ovp.ConvertToGraphQL(fn.Result.Type)

		item := &graphql.Field{
			Description: fn.Prop.Desc,
			Type:        out,
		}
		if fn.RequestArg != nil {
			// 有输入参数
			if fn.Struct != nil && fn.Struct.ExplodeParams {
				// 拆开输入参数
				item.Args = ivp.GenerateArgs(fn.RequestArg.Type)
			} else {
				in := ivp.ConvertToGraphQL(fn.RequestArg.Type)
				item.Args = graphql.FieldConfigArgument{
					"in": &graphql.ArgumentConfig{
						Type: in,
					},
				}
			}
		}

		// 生成 Resolve
		item.Resolve = fn.CreateResolver(ivp, ovp)

		callback(fn.GetResolveName(), item)
	}
}

// ParseStruct 解析结构中的方法
func (rm *ResolverManager) ParseStruct(structInstance interface{}, injectQuery parser.InjectQuery) {
	structDef, err := parser.ParseStruct(structInstance, injectQuery)
	if err != nil {
		rm.log.Error("ParseStruct: " + err.Error())
		return
	}

	typ := reflects.ParseType(reflect.TypeOf(structInstance))

	// 解析函数列表
	mc := typ.Type.NumMethod()
	for n := 0; n < mc; n++ {
		method := typ.Type.Method(n)

		if parser.IsDescFunc(method) {
			// 是描述函数，忽略
			continue
		}

		info := &reflects.FuncInfo{
			Name:    method.Name,
			Struct:  typ.Name,
			Package: typ.Package,
		}
		funcDef, err := parser.ParseFunc(method.Func, method.Type, info, structDef, injectQuery)
		if err != nil {
			rm.log.Error(fmt.Sprintf("Invalid function Signature <%s> : %s", info.String(), err.Error()))
			//rm.log.Error(fmt.Sprintf("&s.%s 无法解析:%s", typ.Name, method.Name, err.Error()))
		} else {
			name := funcDef.GetResolveName()
			if _, ok := rm.resolverMap[name]; ok {
				rm.log.Warn("API " + name + " is exists, method " + info.String() + " is ignored")
				continue
			}

			// 解析函数说明
			funcDef.ParseDescription(structInstance)
			rm.resolverFuncs = append(rm.resolverFuncs, funcDef)
			rm.resolverMap[name] = funcDef
			rm.log.Debug(fmt.Sprintf("Passed %s.%s", typ.Name, info.Name))

			//fns = append(fns, funcDef)
		}
	}
}

// ParserFunction 解析独立函数，不属于任何结构
func (rm *ResolverManager) ParserFunction(funcObj interface{}, injectQuery parser.InjectQuery) {
	// 函数
	info := reflects.ParseFuncInfo(funcObj)
	fnTyp := reflect.TypeOf(funcObj)
	fnObj := reflect.ValueOf(funcObj)
	funcDef, err := parser.ParseFunc(fnObj, fnTyp, info, nil, injectQuery)
	if err != nil {
		rm.log.Error(fmt.Sprintf("Invalid function Signature <%s> : %s", info.String(), err.Error()))
		return
	}

	name := funcDef.GetResolveName()
	if _, ok := rm.resolverMap[name]; ok {
		rm.log.Warn("API " + name + " is exists, method " + info.String() + " is ignored")
		return
	}

	rm.resolverFuncs = append(rm.resolverFuncs, funcDef)
	rm.resolverMap[name] = funcDef
	rm.log.Debug(fmt.Sprintf("Passed %s", info.Name))
}
