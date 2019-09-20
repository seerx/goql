package reflects

import (
	"fmt"
	"reflect"
	"strconv"
)

// TypeDef 类型定义
type TypeDef struct {
	StructField *reflect.StructField // 结构字段信息，可能为 nil
	//FieldName   string               // StructField 相关联，字段名称
	//JSONName    string               // StructField 相关联，gql 字段名称，如果设置 json tag 则使用 json tag，否则使用 field 的名称，如果为空则不操作
	Type     reflect.Type // 原始 type
	RealKind reflect.Kind // 真正的类型
	RealType reflect.Type // 真正的 type，如指针、slice 等需要解析真正的类型

	Name    string // 真正的类型名称
	Package string // 真正的包名称

	IsPtr            bool // 是否指针类型
	IsSlice          bool // 是否 Slice 类型
	IsSliceItemIsPtr bool // Slice 元素是否指针 类型
	IsPrimitive      bool // 是否原生类型
	IsStruct         bool // 是否结构体
	IsInterface      bool // 是否接口
	IsFunc           bool // 是否函数
}

func ParseField(field *reflect.StructField) *TypeDef {
	def := ParseType(field.Type)
	def.StructField = field

	//tag := ParseTag(field)
	//def.JSONName = tag.FieldName
	//def.FieldName = field.Name

	return def
}

func ParseType(typ reflect.Type) *TypeDef {
	var def = &TypeDef{
		Type:        typ,
		RealType:    typ,
		RealKind:    typ.Kind(),
		IsPtr:       typ.Kind() == reflect.Ptr,
		IsSlice:     typ.Kind() == reflect.Slice,
		IsStruct:    typ.Kind() == reflect.Struct,
		IsInterface: typ.Kind() == reflect.Interface,
		IsFunc:      typ.Kind() == reflect.Func,
		Package:     typ.PkgPath(),
		Name:        typ.Name(),
	}
	//pkg := typ.PkgPath()

	if def.IsPtr || def.IsSlice {
		// 指针或切片
		// 递归解析
		sub := ParseType(typ.Elem())
		sub.Type = typ
		sub.IsSlice = def.IsSlice
		if def.IsSlice {
			sub.IsSliceItemIsPtr = sub.IsPtr
		}
		sub.IsPtr = def.IsPtr

		def = sub
	}
	def.IsPrimitive = def.Package == ""

	return def
}

// CheckType 检查类型，递归检查
// 数据类型必须是 原生类型、Slice、指针
func (td *TypeDef) CheckType() error {
	//sets := make(map[string]int)
	return td.Recursive(func(parent *TypeDef, def *TypeDef) error {
		if !def.IsStruct && !def.IsPrimitive {
			return fmt.Errorf("不支持的数据类型 %s", td.Key())
		}
		return nil
	})
	//return td.doCheckType(sets)
}

func (td *TypeDef) Key() string {
	return fmt.Sprintf("%s.%s", td.Package, td.Name)
}

func (td *TypeDef) Summary() string {
	return "package: " + td.Package +
		"\nraw name: " + td.Type.Name() +
		"\nname: " + td.Name +
		"\nslice: " + strconv.FormatBool(td.IsSlice) +
		//"\nslice element is ptr: " + strconv.FormatBool(td.IsElementPtr) +
		"\nptr: " + strconv.FormatBool(td.IsPtr) +
		"\nstruct: " + strconv.FormatBool(td.IsStruct) +
		"\ninterface: " + strconv.FormatBool(td.IsInterface)

	//return fmt.Sprintf("%s.%s", td.Package, td.Name)
}

// Recursive 递归遍历数据类型及其下属数据类型
// @param fn 为遍历到的数据类型，返回 false 时中断
func (td *TypeDef) Recursive(fn func(parent *TypeDef, sub *TypeDef) error) error {
	sets := make(map[string]int)
	return td.doRecursive(sets, nil, fn)
}

// doRecursive 递归遍历数据类型及其下属数据类型
// @param fn 为遍历到的数据类型，返回 false 时中断
// @param sets 防止无限递归
func (td *TypeDef) doRecursive(sets map[string]int, parent *TypeDef, fn func(parent *TypeDef, sub *TypeDef) error) error {
	if td.IsPrimitive {
		return fn(parent, td)
		//return nil
	}
	if td.IsStruct {
		if _, ok := sets[td.Key()]; ok {
			// 已经遍历过的，不再遍历
			return nil
		}
		// 记录已经遍历过的类型
		sets[td.Key()] = 0
		if err := fn(parent, td); err != nil {
			// 不再继续
			return err
		}

		c := td.RealType.NumField()
		for n := 0; n < c; n++ {
			field := td.RealType.Field(n)
			sub := ParseField(&field)
			if err := sub.doRecursive(sets, td, fn); err != nil {
				return err
			}
		}
		return nil
	}

	return fn(parent, td)
}
