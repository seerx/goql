package parser

import (
	"fmt"
	"reflect"

	"github.com/seerx/goql/internal/reflects"
	"github.com/seerx/goql/pkg/require"
)

// prefixField 前缀字段名称
const prefixField = "Prefix"

// StructDef 结构解析
type StructDef struct {
	Instance     interface{} // 注册时传入的结构实例
	Type         reflect.Type
	InjectFields []*InjectPair // 注入字段列表
	RequireField string        // requirement
	Prefix       string        // graphql 方法名称前缀
}

type InjectPair struct {
	Field  string      // 对应结构属性名称
	Inject *InjectInfo // 注入属性
}

// ParseStruct 解析 struct
func ParseStruct(instance interface{},
	injectQuery func(injectType reflect.Type) (*InjectInfo, error)) (*StructDef, error) {
	typ := reflects.ParseType(reflect.TypeOf(instance))
	if !typ.IsStruct {
		return nil, fmt.Errorf("%s 不是 struct 类型", typ.Name)
	}

	def := &StructDef{
		Instance: instance,
		Type:     typ.RealType,
	}

	// 解析注入字段
	fc := typ.RealType.NumField()
	for n := 0; n < fc; n++ {
		fd := typ.RealType.Field(n)
		if fd.Name == prefixField {
			// 解析前缀
			tag := reflects.ParseTag(&fd)
			def.Prefix = tag.Prefix
		} else {
			fdTyp := reflects.ParseField(&fd)
			if require.IsRequirement(fdTyp.Type) {
				// Requirement
				def.RequireField = fd.Name
			} else if fdTyp.IsInterface || fdTyp.IsStruct {
				// 接口或 struct 类型，可以注入
				inject, err := injectQuery(fdTyp.RealType)
				if err == nil && inject != nil {
					// 可以注入
					//inject.Field = &fd
					def.InjectFields = append(def.InjectFields, &InjectPair{
						Inject: inject,
						Field:  fd.Name,
					})
				}
			}
			// 其它类型，不可以注入
		}
	}

	return def, nil
}
