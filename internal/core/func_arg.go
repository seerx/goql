package core

import (
	"reflect"

	"github.com/seerx/goql/pkg/param"

	"github.com/seerx/goql/internal/inject"

	"github.com/graphql-go/graphql"
	"github.com/seerx/goql/internal/reflects"
)

type ArgContext struct {
	Param          *graphql.ResolveParams
	InjectValueMap map[reflect.Type]reflect.Value
	Input          reflect.Value
	Require        *param.Requirement
	//Validator      *param.InputValidator
}

type FuncArg interface {
	CreateValue(ctx *ArgContext) reflect.Value
	IsInjectInterface() bool
}

// StructArg 结构体参数，即函数所属结构体
type StructArg struct {
	Def   *StructDef
	IsPtr bool
	//ValueMap map[reflect.Type]reflect.Value
}

func (arg *StructArg) CreateValue(ctx *ArgContext) reflect.Value {
	// 结构体参数
	var argIn reflect.Value
	var structArg reflect.Value
	if arg.Def != nil {
		structArg = reflect.New(arg.Def.Type)
		if arg.IsPtr {
			argIn = structArg
			structArg = structArg.Elem()
		} else {
			argIn = structArg.Elem()
		}

		if arg.Def.RequireField != "" {
			fd := structArg.FieldByName(arg.Def.RequireField)
			typ := fd.Type()
			if typ.Kind() == reflect.Ptr {
				// 指针
				fd.Set(reflect.ValueOf(ctx.Require))
			} else {
				// 非指针
				fd.Set(reflect.ValueOf(ctx.Require).Elem())
			}
		}

		//elem := structArg.Elem()
		//arg.ValueMap = map[reflect.Type]reflect.Value{}
		// 注入结构中的数据
		for _, fd := range arg.Def.InjectFields {
			fdVal := fd.Inject.InjectValue(ctx.Param)
			sfd := structArg.FieldByName(fd.Field)

			if sfd.Type().Kind() != reflect.Ptr {
				sfd.Set(fdVal.Elem())
			} else {
				sfd.Set(fdVal)
			}

			// 存储注入值
			ctx.InjectValueMap[fd.Inject.Type] = fdVal
		}
	}
	return argIn
}

func (StructArg) IsInjectInterface() bool {
	return false
}

// ResolveParamsArg graphql.ResolveParams 参数
type ResolveParamsArg struct {
	IsPtr bool
}

func (arg *ResolveParamsArg) CreateValue(ctx *ArgContext) reflect.Value {
	return reflect.ValueOf(ctx.Param)
}

func (arg *ResolveParamsArg) IsInjectInterface() bool {
	return false
}

// RequireArg Reqiurement 参数
type RequireArg struct {
	IsPtr bool
}

func (arg *RequireArg) CreateValue(ctx *ArgContext) reflect.Value {
	if arg.IsPtr {
		return reflect.ValueOf(ctx.Require)
	}
	return reflect.ValueOf(ctx.Require).Elem()
}

func (arg *RequireArg) IsInjectInterface() bool {
	return false
}

// InjectArg 注入参数
type InjectArg struct {
	IsInterface bool
	Inject      *inject.InjectInfo
}

func (arg *InjectArg) CreateValue(ctx *ArgContext) reflect.Value {
	val, ok := ctx.InjectValueMap[arg.Inject.Type]
	if ok {
		return val
	}
	val = arg.Inject.InjectValue(ctx.Param)
	ctx.InjectValueMap[arg.Inject.Type] = val
	return val
}

func (arg *InjectArg) IsInjectInterface() bool {
	return arg.IsInterface
}

// RequestArg 请求参数
type RequestArg struct {
	ArgType *reflects.TypeDef
}

func (arg *RequestArg) CreateValue(ctx *ArgContext) reflect.Value {
	// 待实现，从 gqlParam 中解析数据到 ArgType 中，并返回
	return ctx.Input
}

func (arg *RequestArg) IsInjectInterface() bool {
	return false
}
