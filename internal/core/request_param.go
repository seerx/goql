package core

import (
	"fmt"
	"reflect"

	"github.com/seerx/goql/pkg/param"

	"github.com/seerx/goql/internal/varspool"

	"github.com/graphql-go/graphql"
	"github.com/seerx/goql/internal/reflects"
)

// 解析请求参数
type RequestParam struct {
	InputVars *varspool.InputVarsPool
	Require   *param.Requirement
}

func (rp *RequestParam) Parse(p *graphql.ResolveParams, typ *reflects.TypeDef, explodeParams bool) (reflect.Value, error) {
	ivar := rp.InputVars.FindInputVar(typ.Type)
	rp.Require = param.New()
	if explodeParams {
		return rp.parseParam("", ivar, p.Args, false)
	}
	in := p.Args["in"]

	//if !ok {
	//	return reflect.ValueOf(nil), fmt.Errorf("Required arguments: in")
	//}
	// 如果打散顶层结构，则使用 p.Args 做参数
	//ivar := rp.InputVars.FindInputVar(typ.Type)
	// 开始递归解析
	return rp.parseParam("", ivar, in, false)
}

func (rp *RequestParam) parseParam(parentParam string, ivar *varspool.InputVar, input interface{}, forcePtr bool) (reflect.Value, error) {
	//if input == nil {
	//	//fmt.Println("解析:", parentParam+"."+ivar.JSONName, " 无数据")
	//	//if parentParam != "" {
	//	//	rp.Require.Add(parentParam + "." + ivar.JSONName)
	//	//} else {
	//	//	rp.Require.Add(ivar.JSONName)
	//	//}
	//	return reflect.ValueOf(nil), nil
	//}
	typ := ivar.Type
	if typ.IsSlice {
		// 列表
		if input == nil {
			return reflect.MakeSlice(typ.Type, 0, 0), nil
		}
		//ivp := rp.InputVars
		//vr := ivp.FindInputVar(typ.RealType)
		ary, ok := input.([]interface{})
		if !ok {
			panic(fmt.Errorf("Cann't parse %s as slice", ivar.FieldName))
		}

		slice := reflect.MakeSlice(typ.Type, 0, len(ary))
		for _, v := range ary {
			if v == nil {
				continue
			}
			val, err := rp.parseParam(parentParam, ivar.ItemVar, v, true)
			if err != nil {
				return reflect.ValueOf(nil), err
			}

			if typ.IsSliceItemIsPtr {
				slice = reflect.Append(slice, val)
			} else {
				slice = reflect.Append(slice, val.Elem())
			}
		}
		return slice, nil
	} else if typ.IsStruct {
		//ivp := rp.InputVars
		//u := ivp.ConvertToGraphQL(typ.Type)
		// 结构
		val := reflect.New(typ.RealType)
		elem := val.Elem()

		if input != nil {
			data, ok := input.(map[string]interface{})
			if !ok {
				panic(fmt.Errorf("Cann't parse %s as struct", ivar.FieldName))
			}

			for _, child := range ivar.Children {
				v, ok := data[child.JSONName]
				pParam := child.JSONName
				if parentParam != "" {
					pParam = parentParam + "." + pParam
				}

				if ok {
					fd := elem.FieldByName(child.FieldName)
					//fd.Set()
					wrapedData, err := rp.parseParam(pParam, child, v, false)
					if err != nil {
						return reflect.ValueOf(nil), err
					}
					fd.Set(wrapedData)
					rp.Require.Add(pParam)
				}
			}
		}
		if typ.IsPtr || forcePtr {
			return val, nil
		}
		return elem, nil
	} else if typ.IsPrimitive {
		// 原生类型
		if input == nil {
			return reflect.ValueOf(nil), nil
		}
		// 检查数据合法性
		for _, cker := range ivar.Validators {
			if err := cker.Check(input); err != nil {
				panic(err)
			}
		}
		//var a int = 10
		////b := float64(a)
		////var aa int64 = 10
		//b := float64(a)
		//fmt.Printf("%f\n", b)
		p := convert(typ, input)
		if p == nil {
			panic(fmt.Errorf("参数 %s 无法转换为 %s", parentParam, typ.Name))
		}

		return reflect.ValueOf(p), nil
	}
	return reflect.ValueOf(nil), fmt.Errorf("Unsupport data type : " + typ.Key())
}

func convert(typ *reflects.TypeDef, val interface{}) interface{} {
	valType := reflect.TypeOf(val)
	if valType.Kind() == reflect.Ptr {
		valType = valType.Elem()
	}
	if valType.ConvertibleTo(typ.RealType) {
		var Int int
		var Float64 float64
		switch valType.Kind() {
		case reflect.Int:
			Int = val.(int)
			return intToType(typ.RealType.Kind(), Int)
		case reflect.Float64:
			Float64 = val.(float64)
			return float64ToType(typ.RealType.Kind(), Float64)
		case reflect.String:
			return val
		}

	}
	return val
}

func intToType(kind reflect.Kind, val int) interface{} {
	// 可以转换
	switch kind {
	case reflect.Int:
		return int(val)
	case reflect.Int8:
		return int8(val)
	case reflect.Int32:
		return int32(val)
	case reflect.Int64:
		return int64(val)

	case reflect.Uint:
		return uint(val)
	case reflect.Uint8:
		return uint8(val)
	case reflect.Uint16:
		return uint16(val)
	case reflect.Uint32:
		return uint32(val)
	case reflect.Uint64:
		return uint64(val)

	case reflect.Float32:
		return float32(val)
	case reflect.Float64:
		return val
	}
	return val
}

func float64ToType(kind reflect.Kind, val float64) interface{} {
	// 可以转换
	switch kind {
	case reflect.Int:
		return int(val)
	case reflect.Int8:
		return int8(val)
	case reflect.Int32:
		return int32(val)
	case reflect.Int64:
		return int64(val)

	case reflect.Uint:
		return uint(val)
	case reflect.Uint8:
		return uint8(val)
	case reflect.Uint16:
		return uint16(val)
	case reflect.Uint32:
		return uint32(val)
	case reflect.Uint64:
		return uint64(val)

	case reflect.Float32:
		return float32(val)
	case reflect.Float64:
		return val
	}
	return val
}
