package varspool

import (
	"fmt"
	"reflect"

	"github.com/seerx/goql/pkg/log"

	"github.com/seerx/goql/internal/types"

	"github.com/graphql-go/graphql"
	"github.com/seerx/goql/internal/reflects"
)

//func ()

type outputPair struct {
	pkg    string
	output graphql.Output
}

// OutputVarsPool 输出变量管理
type OutputVarsPool struct {
	log    log.Logger
	varMap map[string]*outputPair
}

// NewOutputVarsPool 新建输出对象管理
func NewOutputVarsPool(log log.Logger) *OutputVarsPool {
	return &OutputVarsPool{
		log:    log,
		varMap: map[string]*outputPair{},
	}
}

// ConvertToGraphQL 转换
func (ovp *OutputVarsPool) ConvertToGraphQL(typ reflect.Type) graphql.Output {
	tag := typ.Name() + "@" + typ.PkgPath()
	var out graphql.Output
	tp := reflects.ParseType(typ)
	if tp.IsPrimitive {
		// 原生类型
		gobj := primitiveTypeToGraphQLType(tp.RealType)
		if gobj == nil {
			panic("Unknown primitive Type: " + tag)
		}
		out = gobj
		if tp.IsSlice {
			// 切片
			out = graphql.NewList(gobj)
		}
		//ovp.varMap[typ] = out
		return out
	}

	varName := tp.Name
	if tp.IsSlice {
		varName += "s"
	}

	var ok bool
	// 去变量表中查找
	//os.IsExist()
	existType, ok := ovp.varMap[varName]
	//out, ok = ovp.varMap[varName]
	if ok {
		// 找到，直接返回
		if existType.pkg != tp.Package {
			// 名称相同，但是所在包不同，报出错误
			panic(fmt.Errorf("Output type <%s> exist while defining from package [%s] \n<%s> is defined in package [%s] before", varName, tp.Package, varName, existType.pkg))
		}
		out = existType.output
		ovp.log.Debug("Find in pool:" + tag)
		return out
	}
	// 没找到

	ovp.log.Debug("Create type: " + varName)

	// 非原生类型，一般指结构类型
	//name := tp.Name
	objFields := graphql.Fields{}
	gobj := graphql.NewObject(graphql.ObjectConfig{
		Name:   varName,
		Fields: objFields,
	})

	if tp.IsSlice {
		// 切片
		out = graphql.NewList(gobj)
	} else {
		out = gobj
	}
	ovp.varMap[varName] = &outputPair{
		output: out,
		pkg:    tp.Package,
	}

	for n := 0; n < tp.RealType.NumField(); n++ {
		fd := tp.RealType.Field(n)
		//fdType := reflects.ParseField(&fd)
		tg := reflects.ParseTag(&fd)
		if tg.FieldName == "" {
			continue
		}
		gqlField := &graphql.Field{
			Type:        ovp.ConvertToGraphQL(fd.Type),
			Name:        tg.FieldName,
			Description: tg.Description,
		}
		objFields[tg.FieldName] = gqlField
	}

	return out
}

func primitiveTypeToGraphQLType(typ reflect.Type) graphql.Output {
	kind := typ.Kind()
	if types.IsInt(kind) {
		return graphql.Int
	}
	if types.IsFloat(kind) {
		return graphql.Float
	}
	if types.IsString(kind) {
		return graphql.String
	}
	if types.IsBool(kind) {
		return graphql.Boolean
	}
	if types.IsTime(typ) {
		return graphql.DateTime
	}

	return nil
}

func primitiveTypeToGraphQLInput(typ reflect.Type) graphql.Input {
	kind := typ.Kind()
	if types.IsInt(kind) {
		return graphql.Int
	}
	if types.IsFloat(kind) {
		return graphql.Float
	}
	if types.IsString(kind) {
		return graphql.String
	}
	if types.IsBool(kind) {
		return graphql.Boolean
	}
	if types.IsTime(typ) {
		return graphql.DateTime
	}

	return nil
}
