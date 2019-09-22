package parser

import (
	"fmt"
	"reflect"

	"github.com/seerx/goql/pkg/require"

	"github.com/graphql-go/graphql"
	"github.com/seerx/goql/internal/reflects"
)

// 解析请求参数
type RequestParam struct {
	InputVars *InputVarsPool
	Require   *require.Requirement
}

func (rp *RequestParam) Parse(p *graphql.ResolveParams, typ *reflects.TypeDef, explodeParams bool) (reflect.Value, error) {
	ivar := rp.InputVars.FindInputVar(typ.Type)
	rp.Require = require.New()
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

func (rp *RequestParam) parseParam(parentParam string, ivar *InputVar, input interface{}, forcePtr bool) (reflect.Value, error) {
	if input == nil {
		//fmt.Println("解析:", parentParam+"."+ivar.JSONName, " 无数据")
		//if parentParam != "" {
		//	rp.Require.Add(parentParam + "." + ivar.JSONName)
		//} else {
		//	rp.Require.Add(ivar.JSONName)
		//}
		return reflect.ValueOf(nil), nil
	}
	typ := ivar.Type
	if typ.IsSlice {
		// 列表
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
		data, ok := input.(map[string]interface{})
		if !ok {
			panic(fmt.Errorf("Cann't parse %s as struct", ivar.FieldName))
		}

		// 结构
		val := reflect.New(typ.RealType)
		elem := val.Elem()

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
			} else {
				rp.Require.Add(pParam)
			}
		}
		if typ.IsPtr || forcePtr {
			return val, nil
		}
		return elem, nil
	} else if typ.IsPrimitive {
		// 原生类型
		// 检查数据合法性
		for _, cker := range ivar.Validators {
			if err := cker.Check(input); err != nil {
				panic(err)
			}
		}
		return reflect.ValueOf(input), nil
	}
	return reflect.ValueOf(nil), fmt.Errorf("Unsupport data type : " + typ.Key())
}
