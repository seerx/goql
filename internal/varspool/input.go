package varspool

import (
	"fmt"
	"reflect"

	"github.com/seerx/goql/internal/reflects/validators"

	"github.com/seerx/goql/internal/reflects"

	"github.com/seerx/goql/pkg/log"

	"github.com/graphql-go/graphql"
)

type InputVar struct {
	GqlInput   graphql.Input
	JSONName   string
	FieldName  string
	Descrpiton string
	Type       *reflects.TypeDef
	ItemVar    *InputVar              // Slice 元素类型
	Children   []*InputVar            // 结构属性类型
	Validators []validators.Validator // 数据有效性检查
	pkg        string                 // 结构定义所在包
}

type gtpair struct {
	Input graphql.Input
	Var   *InputVar
}

// InputVarsPool 输入变量管理
type InputVarsPool struct {
	log          log.Logger
	gqlObjectMap map[string]gtpair // graphql.Input
	varMap       map[reflect.Type]*InputVar
}

// NewInputVarsPool 新建输入对象管理
func NewInputVarsPool(log log.Logger) *InputVarsPool {
	return &InputVarsPool{
		log:          log,
		gqlObjectMap: map[string]gtpair{},
		varMap:       map[reflect.Type]*InputVar{},
	}
}

func (ivp *InputVarsPool) FindInputVar(typ reflect.Type) *InputVar {
	if iv, ok := ivp.varMap[typ]; ok {
		return iv
	}
	return nil
}

func (ivp *InputVarsPool) GenerateArgs(typ reflect.Type) graphql.FieldConfigArgument {
	_, vars := ivp.convertToGraphQL(typ, true)
	//vars.Children
	//args := &graphql.FieldConfigArgument{}
	mp := map[string]*graphql.ArgumentConfig{}
	//args := []*graphql.FieldConfigArgument{}
	for _, v := range vars.Children {

		mp[v.JSONName] = &graphql.ArgumentConfig{
			Type:        v.GqlInput,
			Description: v.Descrpiton,
		}
	}
	return graphql.FieldConfigArgument(mp)
}

func (ivp *InputVarsPool) ConvertToGraphQL(typ reflect.Type) graphql.Input {
	input, _ := ivp.convertToGraphQL(typ, false)
	return input
}

func (ivp *InputVarsPool) convertToGraphQL(typ reflect.Type, rootParam bool) (graphql.Input, *InputVar) {
	tag := typ.Name() + "@" + typ.PkgPath()
	tp := reflects.ParseType(typ)
	if tp.IsPrimitive {
		// 原生类型
		gobj := primitiveTypeToGraphQLInput(tp.RealType)
		if gobj == nil {
			panic("Unknown primitive Type: " + tag)
		}
		//in = &InputVar{
		//	GqlInput: gobj,
		//	Type:     tp,
		//}
		if tp.IsSlice {
			// 切片
			//in.GqlInput = graphql.NewList(gobj)
			gobj = graphql.NewList(gobj)
		}
		//ivp.varMap[typ] = in
		return gobj, nil
	}
	// 非原生类型

	varName := "in_" + tp.Name
	if tp.IsSlice {
		varName += "s"
	}

	var in *InputVar
	var sliceItamVar *InputVar
	//var ok bool
	var gobj graphql.Input
	objFields := graphql.InputObjectConfigFieldMap{}
	// 去变量表中查找
	pair, ok := ivp.gqlObjectMap[varName]
	if ok {
		if pair.Var.pkg != tp.Package {
			// 名称相同，但是所在包不同，报出错误
			panic(fmt.Errorf("Input type <%s> exist while defining from package [%s] \n<%s> is defined in package [%s] before", varName, tp.Package, varName, pair.Var.pkg))
		}

		// 找到，直接返回
		ivp.log.Debug("Find in pool:" + varName)
		//ivar, ok := ivp.varMap[typ]
		//if !ok {
		//	panic(fmt.Errorf("Cann't find InputVar with %s", varName))
		//}
		return pair.Input, pair.Var
	} else {
		// 没有找到
		ivp.log.Debug("Create input type: " + varName)
		if tp.IsSlice {
			// 列表
			var itemType graphql.Input
			itemType, sliceItamVar = ivp.convertToGraphQL(tp.RealType, false)
			gobj = graphql.NewList(itemType)
		} else {
			// 结构
			// 注册单个查询对象
			gobj = graphql.NewInputObject(
				graphql.InputObjectConfig{
					Name:   varName,
					Fields: objFields,
				})
		}
	}

	in = &InputVar{
		Type:     tp,
		ItemVar:  sliceItamVar,
		GqlInput: gobj,
	}

	//if tp.IsSlice {
	//	fmt.Println(".... set ", varName, " item type ...")
	//	in.ItemVar = sliceItamVar
	//}
	//if tp.IsSlice {
	//	// 切片
	//	in.GqlInput = graphql.NewList(gobj)
	//} else {
	//	in.GqlInput = gobj
	//}
	if !rootParam {
		// 注册 graphql 类型
		ivp.gqlObjectMap[varName] = gtpair{
			Input: gobj,
			Var:   in,
		}
	}
	// 注册参数类型
	ivp.varMap[typ] = in

	//fmt.Println("Struct", tp.Name)
	// 确定是结构
	for n := 0; n < tp.RealType.NumField(); n++ {
		fd := tp.RealType.Field(n)
		//fdType := reflects.ParseField(&fd)
		tg := reflects.ParseTag(&fd)
		if tg.FieldName == "" {
			continue
		}
		gvar, ivar := ivp.convertToGraphQL(fd.Type, false)
		child := &InputVar{
			GqlInput:   gvar,
			JSONName:   tg.FieldName,
			FieldName:  fd.Name,
			Type:       reflects.ParseField(&fd),
			Descrpiton: tg.Description,
			//Children:  ivar.Children,
		}
		if ivar == nil {
			// 原生类型，解析数据检查列表
			child.Validators = tg.GenerateValidators(fd.Type)

		} else {
			// 非原生类型
			child.ItemVar = ivar.ItemVar
			child.Children = ivar.Children
		}
		//child.JSONName = tg.FieldName
		//child.FieldName = fd.Name
		in.Children = append(in.Children, child)
		gqlField := &graphql.InputObjectFieldConfig{
			Type:        gvar,
			Description: tg.Description,
			//Name:        tg.FieldName,
		}
		//fmt.Println("Struct", tp.Name, "Field", tg.FieldName)
		objFields[tg.FieldName] = gqlField
	}

	return gobj, in
}
